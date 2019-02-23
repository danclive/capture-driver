package modbus

import (
	"strings"
	"sync"
	"time"

	nson "github.com/danclive/nson-go"
	queen "github.com/danclive/queen-go"
	"github.com/goburrow/modbus"
)

var (
	PREFIX        = ""
	CONNECT       = PREFIX + "driver/modbus/connect"
	CONNECT_ACK   = PREFIX + "driver/modbus/connect/ack"
	RECONNECT     = PREFIX + "driver/modbus/reconnect"
	RECONNECT_ACK = PREFIX + "driver/modbus/reconnect/ack"
	READ          = PREFIX + "driver/modbus/read"
	READ_ACK      = PREFIX + "driver/modbus/read/ack"
	WRITE         = PREFIX + "driver/modbus/write"
	WRITE_ACK     = PREFIX + "driver/modbus/write/ack"
	CLOSE         = PREFIX + "driver/modbus/close"
	CLOSE_ACK     = PREFIX + "driver/modbus/close/ack"
	STATUS        = PREFIX + "driver/modbus/status"
	STATUS_ACK    = PREFIX + "driver/modbus/status/ack"
)

func InitModbus(queen *queen.Queen) {
	queen.On(CONNECT, connect)
	queen.On(RECONNECT, reconnect)
	queen.On(READ, read)
	queen.On(WRITE, write)
	queen.On(CLOSE, close)
	queen.On(STATUS, status)

	msg := nson.Message{"driver": nson.String("modbus"), "ok": nson.Bool(true)}
	queen.Emit("init/driver/ack", msg)
}

var conns = make(map[int64]*conn_t)
var lock sync.RWMutex

type conn_t struct {
	id         int64
	config     string
	istcp      bool
	address    string
	isbig      bool
	islink     bool  // 是否已经连接
	retry      int32 // 是否重试&重试时间
	tick       int32 // 读间隔
	tags       nson.Array
	client     modbus.Client
	tcpHandler *modbus.TCPClientHandler
	rtuHandler *modbus.RTUClientHandler
	lock       sync.Mutex
	// isBIGEndian bool
}

func parse_config(msg *nson.Message) (config string, istcp bool, address string, isbig bool, ok bool) {
	var err error
	config, err = msg.GetString("config")
	if err != nil {
		msg.Insert("ok", nson.Bool(false))
		msg.Insert("error", nson.String("Message format error: can't not get config!"))
		return
	}

	if strings.HasPrefix(config, "TCP") {
		istcp = true
	}

	i := strings.Index(config, "?")

	if i < 0 {
		if istcp {
			address = config[4:]
		} else {
			address = config[7:]
		}
	} else {
		if istcp {
			address = config[4:i]
		} else {
			address = config[7:i]
		}

		options := config[i+1:]

		if options == "isBIGEndian=false" {
			isbig = false
		} else {
			isbig = true
		}
	}

	ok = true
	return
}

// TCP:localhost:502?isBIGEndian=false
// Serial:/dev/ttyS0?isBIGEndian=false
func connect(context queen.Context) {
	msg, ok := context.Message.(nson.Message)
	if !ok {
		msg.Insert("ok", nson.Bool(false))
		msg.Insert("error", nson.String("Message format error: message isn't nson!"))
		context.Queen.Emit(CONNECT_ACK, msg)
		return
	}

	if id, err := msg.GetI64("id"); err == nil {
		lock.Lock()
		conn, ok := conns[id]
		// lock.Unlock()

		if ok {
			conn.lock.Lock()
			if conn.istcp {
				conn.tcpHandler.Close()
			} else {
				conn.rtuHandler.Close()
			}
			conn.lock.Unlock()

			msg.Insert("cover", nson.String(conn.config))
		}

		config, istcp, address, isbig, ok := parse_config(&msg)
		if !ok {
			context.Queen.Emit(CONNECT_ACK, msg)
			return
		}

		retry, _ := msg.GetI32("retry")

		if istcp {
			handler := modbus.NewTCPClientHandler(address)

			conn2 := conn_t{
				id:         id,
				config:     config,
				istcp:      istcp,
				isbig:      isbig,
				address:    address,
				retry:      retry,
				tcpHandler: handler,
			}

			err := handler.Connect()
			if err != nil {
				msg.Insert("ok", nson.Bool(false))
				msg.Insert("error", nson.String("modbus connect error: "+err.Error()))

				if retry > 0 {
					context.Queen.Emit(
						RECONNECT,
						nson.Message{"id": nson.I64(id), "retry": nson.I32(retry)})
				}
			} else {
				conn2.islink = true

				conn2.client = modbus.NewClient(handler)

				msg.Insert("ok", nson.Bool(true))
			}

			conns[id] = &conn2
			lock.Unlock()
		} else {
			handler := modbus.NewRTUClientHandler(address)

			conn2 := conn_t{
				id:         id,
				config:     config,
				istcp:      istcp,
				isbig:      isbig,
				address:    address,
				retry:      retry,
				rtuHandler: handler,
			}

			err := handler.Connect()
			if err != nil {
				msg.Insert("ok", nson.Bool(false))
				msg.Insert("error", nson.String("modbus connect error: "+err.Error()))

				if retry > 0 {
					context.Queen.Emit(
						RECONNECT,
						nson.Message{"id": nson.I32(id), "retry": nson.I32(retry)})
				}
			} else {
				conn2.islink = true

				conn2.client = modbus.NewClient(handler)

				msg.Insert("ok", nson.Bool(true))
			}

			conns[id] = &conn2
			lock.Unlock()
		}
	} else {
		msg.Insert("ok", nson.Bool(false))
		msg.Insert("error", nson.String("Message format error: can't get id!"))
	}

	context.Queen.Emit(CONNECT_ACK, msg)
}

func reconnect(context queen.Context) {
	msg, ok := context.Message.(nson.Message)
	if !ok {
		msg.Insert("ok", nson.Bool(false))
		msg.Insert("error", nson.String("Message format error: message isn't nson!"))
		context.Queen.Emit(RECONNECT_ACK, msg)
		return
	}

	retry, err := msg.GetI32("retry")
	if err != nil {
		msg.Insert("ok", nson.Bool(false))
		msg.Insert("error", nson.String("Message format error: can't get retry!"))
		context.Queen.Emit(RECONNECT_ACK, msg)
		return
	}

	time.Sleep(time.Duration(retry) * time.Second)

	if id, err := msg.GetI64("id"); err == nil {
		lock.Lock()
		conn, ok := conns[id]
		if ok {
			conn.lock.Lock()
			if !conn.islink {
				if conn.istcp {
					err := conn.tcpHandler.Connect()
					if err != nil {
						context.Queen.Emit(RECONNECT,
							nson.Message{"id": nson.I32(id), "retry": nson.I32(retry)})

						msg.Insert("ok", nson.Bool(false))
						msg.Insert("error", nson.String("modbus connect error: "+err.Error()))
					} else {
						conn.client = modbus.NewClient(conn.tcpHandler)
						conn.islink = true
					}
				} else {
					err := conn.rtuHandler.Connect()
					if err != nil {
						context.Queen.Emit(RECONNECT,
							nson.Message{"id": nson.I32(id), "retry": nson.I32(retry)})

						msg.Insert("ok", nson.Bool(false))
						msg.Insert("error", nson.String("modbus connect error: "+err.Error()))
					} else {
						conn.client = modbus.NewClient(conn.rtuHandler)
						conn.islink = true
					}
				}
			}

			conn.lock.Unlock()
		}
		lock.Unlock()
	} else {
		msg.Insert("ok", nson.Bool(false))
		msg.Insert("error", nson.String("Message format error: can't get id!"))
	}

	context.Queen.Emit(RECONNECT_ACK, msg)
}

func read(context queen.Context) {
	msg, ok := context.Message.(nson.Message)
	if !ok {
		msg.Insert("ok", nson.Bool(false))
		msg.Insert("error", nson.String("Message format error: message isn't nson!"))
		context.Queen.Emit(READ_ACK, msg)
		return
	}

	if id, err := msg.GetI64("id"); err == nil {
		// 是否是紧急数据
		urgent, _ := msg.GetBool("urgent")

		lock.RLock()
		if conn, ok := conns[id]; ok {
			conn.lock.Lock()

			tick, err := msg.GetI32("tick")
			if err != nil {
				tick = conn.tick
			} else {
				conn.tick = tick
			}

			if conn.islink {
				// 尝试获取tags,如果是首次读的话,必须提供tags
				if tags, err := msg.GetArray("tags"); err == nil {

					// 如果不是紧急数据,就将tags写入到conn中,以便定时读的时候能够获取
					if !urgent {
						conn.tags = tags
					}

					read_ext(context, &msg, conn, tags)
				} else {
					// 如果没有从传入的msg中获取到tags的话,可能是定时读,先判断是否有tags
					// 要注意,这里如果tags为空的话,应该传入参数错误,应当终止读取
					if len(conn.tags) == 0 {
						// 清除定时读
						conn.tick = 0
						// msg.Insert("ok", false)
						// msg.Insert("error", "Message format error: can't get tags!")
					} else {
						read_ext(context, &msg, conn, conn.tags)
					}
				}

				if conn.tick > 0 {
					go func(id int64, tick int32, event_emit *queen.Queen) {
						// 延迟
						time.Sleep(time.Duration(tick) * time.Millisecond)

						msg2 := nson.Message{"id": nson.I64(id)}
						event_emit.Emit(READ, msg2)

					}(id, tick, context.Queen)
				}
			}

			conn.lock.Unlock()
		} else {
			msg.Insert("ok", nson.Bool(false))
			msg.Insert("error", nson.String("Message format error: can't get conn!"))
		}
		lock.RUnlock()
	} else {
		msg.Insert("ok", nson.Bool(false))
		msg.Insert("error", nson.String("Message format error: can't get id!"))
	}

	context.Queen.Emit(READ_ACK, msg)
}

func write(context queen.Context) {
	msg, ok := context.Message.(nson.Message)
	if !ok {
		msg.Insert("ok", nson.Bool(false))
		msg.Insert("error", nson.String("Message format error: message isn't nson!"))
		context.Queen.Emit(WRITE_ACK, msg)
		return
	}

	if id, err := msg.GetI64("id"); err == nil {
		lock.RLock()
		if conn, ok := conns[id]; ok {
			conn.lock.Lock()

			if conn.islink {
				if tags, err := msg.GetArray("tags"); err == nil {
					write_ext(context, &msg, conn, tags)
				}
			}

			conn.lock.Unlock()
		} else {
			msg.Insert("ok", nson.Bool(false))
			msg.Insert("error", nson.String("Message format error: can't get conn!"))
		}
		lock.RUnlock()
	} else {
		msg.Insert("ok", nson.Bool(false))
		msg.Insert("error", nson.String("Message format error: can't get id!"))
	}

	context.Queen.Emit(WRITE_ACK, msg)
}

func close(context queen.Context) {
	msg, ok := context.Message.(nson.Message)
	if !ok {
		msg.Insert("ok", nson.Bool(false))
		msg.Insert("error", nson.String("Message format error: message isn't nson!"))
		context.Queen.Emit(CLOSE_ACK, msg)
		return
	}

	if id, err := msg.GetI64("id"); err == nil {
		lock.Lock()

		if conn, ok := conns[id]; ok {
			conn.lock.Lock()
			if conn.istcp {
				conn.tcpHandler.Close()
			} else {
				conn.rtuHandler.Close()
			}

			delete(conns, id)
			conn.lock.Unlock()

			msg.Insert("ok", nson.Bool(true))
		} else {
			msg.Insert("ok", nson.Bool(false))
			msg.Insert("error", nson.String("Message format error: can't get conn!"))
		}

		lock.Unlock()
	} else {
		msg.Insert("ok", nson.Bool(false))
		msg.Insert("error", nson.String("Message format error: can't get id!"))
	}
}

func status(context queen.Context) {
	context.Queen.Emit(STATUS_ACK, nson.Message{"ok": nson.Bool(true)})
}
