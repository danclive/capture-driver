package snap7

import (
	"net/url"
	"strconv"
	"sync"
	"time"

	nson "github.com/danclive/nson-go"
	queen "github.com/danclive/queen-go"
	snap7 "github.com/danclive/snap7-go"
)

const (
	CONNECT       = "driver/snap7/connect"
	CONNECT_ACK   = "driver/snap7/connect/ack"
	RECONNECT     = "driver/snap7/reconnect"
	RECONNECT_ACK = "driver/snap7/reconnect/ack"
	READ          = "driver/snap7/read"
	READ_ACK      = "driver/snap7/read/ack"
	WRITE         = "driver/snap7/write"
	WRITE_ACK     = "driver/snap7/write/ack"
	CLOSE         = "driver/snap7/close"
	CLOSE_ACK     = "driver/snap7/close/ack"
)

func InitSnap7(event_emiter *queen.EventEmiter) {
	event_emiter.On(CONNECT, connect)
	event_emiter.On(READ, read)
	event_emiter.On(WRITE, write)
	event_emiter.On(CLOSE, close)
	event_emiter.On("driver/snap7/status", status)

	event_emiter.On(RECONNECT, reconnect)

	doc := nson.Message{"driver": nson.String("snap7"), "ok": nson.Bool(true)}
	event_emiter.Emit("init/driver/ack", doc)
}

var conns = make(map[int32]*conn_t)
var lock sync.RWMutex

type conn_t struct {
	id     int32
	config string
	host   string
	rank   int
	slot   int
	islink bool  // 是否已经连接
	retry  int32 // 是否重试&重试时间
	tick   int32 // 读
	tags   nson.Array
	client snap7.Snap7Client
	lock   sync.Mutex
	// isBIGEndian bool
}

func parse_config(msg *nson.Message) (config string, host string, rank int, slot int, ok bool) {
	var err error
	config, err = msg.GetString("config")
	if err != nil {
		msg.Insert("ok", nson.Bool(false))
		msg.Insert("error", nson.String("Message format error: can't not get config!"))
		return
	}

	params, err := url.Parse(config)
	if err != nil {
		msg.Insert("ok", nson.Bool(false))
		msg.Insert("error", nson.String("Config parse error: "+err.Error()))
		return
	}

	// S7-TCP://192.168.0.2?rank=0&slot=2&isBIGEndian=true
	query, err := url.ParseQuery(params.RawQuery)
	if err != nil {
		msg.Insert("ok", nson.Bool(false))
		msg.Insert("error", nson.String("Config parse error: "+err.Error()))
		return
	}

	foo := func(query url.Values, key string) (int, bool) {
		if rank, ok := query[key]; ok && len(rank) > 0 {
			foo_int, err := strconv.Atoi(rank[0])
			if err != nil {
				msg.Insert("ok", nson.Bool(false))
				msg.Insert("error", nson.String("Config parse error, "+key+" : "+err.Error()))
				return 0, false
			}

			return foo_int, true
		} else {
			msg.Insert("ok", nson.Bool(false))
			msg.Insert("error", nson.String("Config parse error: can't get "+key))
			return 0, false
		}
	}

	rank, ok = foo(query, "rank")
	if !ok {
		return
	}

	slot, ok = foo(query, "slot")
	if !ok {
		return
	}

	host = params.Host
	ok = true
	return
}

func connect(context queen.Context) {
	msg, ok := context.Message.(nson.Message)
	if !ok {
		msg.Insert("ok", nson.Bool(false))
		msg.Insert("error", nson.String("Message format error: message isn't nson!"))
		context.Event_emiter.Emit(CONNECT_ACK, msg)
		return
	}

	if id, err := msg.GetI32("id"); err == nil {
		lock.Lock()
		conn, ok := conns[id]
		// lock.Unlock()

		if ok {
			conn.lock.Lock()
			conn.client.Close()
			conn.lock.Unlock()

			msg.Insert("cover", nson.String(conn.config))
		}

		config, host, rank, slot, ok := parse_config(&msg)
		if !ok {
			context.Event_emiter.Emit(CONNECT_ACK, msg)
			return
		}

		retry, _ := msg.GetI32("retry")

		client, err := snap7.ConnentTo(host, rank, slot)

		conn2 := conn_t{
			id:     id,
			config: config,
			host:   host,
			rank:   rank,
			slot:   slot,
			retry:  retry,
			client: client,
		}

		if err != nil {
			msg.Insert("ok", nson.Bool(false))
			msg.Insert("error", nson.String("Snap7 connect error: "+err.Error()))

			if retry > 0 {
				context.Event_emiter.Emit(
					RECONNECT,
					nson.Message{"id": nson.I32(id), "retry": nson.I32(retry)})
			}
		} else {
			conn2.islink = true
			msg.Insert("ok", nson.Bool(true))
		}

		// lock.Lock()
		conns[id] = &conn2
		lock.Unlock()

		// msg.Insert("ok", true)
		// context.Event_emiter.Emit(CONNECT_ACK, msg)
	} else {
		msg.Insert("ok", nson.Bool(false))
		msg.Insert("error", nson.String("Message format error: can't get id!"))
		// context.Event_emiter.Emit(CONNECT_ACK, msg)
	}

	context.Event_emiter.Emit(CONNECT_ACK, msg)
}

func reconnect(context queen.Context) {
	msg, ok := context.Message.(nson.Message)
	if !ok {
		msg.Insert("ok", nson.Bool(false))
		msg.Insert("error", nson.String("Message format error: message isn't nson!"))
		context.Event_emiter.Emit(RECONNECT_ACK, msg)
		return
	}

	retry, err := msg.GetI32("retry")
	if err != nil {
		msg.Insert("ok", nson.Bool(false))
		msg.Insert("error", nson.String("Message format error: can't get retry!"))
		context.Event_emiter.Emit(RECONNECT_ACK, msg)
		return
	}

	time.Sleep(time.Duration(retry) * time.Second)

	if id, err := msg.GetI32("id"); err == nil {
		lock.Lock()
		conn, ok := conns[id]
		if ok {
			conn.lock.Lock()
			if !conn.islink {
				client, err := snap7.ConnentTo(conn.host, conn.rank, conn.slot)
				if err != nil {

					context.Event_emiter.Emit(RECONNECT,
						nson.Message{"id": nson.I32(id), "retry": nson.I32(retry)})

					msg.Insert("ok", nson.Bool(false))
					msg.Insert("error", nson.String("Snap7 connect error: "+err.Error()))
				} else {
					conn.client = client
					conn.islink = true

					msg.Insert("ok", nson.Bool(true))
				}
			}

			conn.lock.Unlock()
		}
		lock.Unlock()
	} else {
		msg.Insert("ok", nson.Bool(false))
		msg.Insert("error", nson.String("Message format error: can't get id!"))
	}

	context.Event_emiter.Emit(RECONNECT_ACK, msg)
}

func read(context queen.Context) {
	msg, ok := context.Message.(nson.Message)
	if !ok {
		msg.Insert("ok", nson.Bool(false))
		msg.Insert("error", nson.String("Message format error: message isn't nson!"))
		context.Event_emiter.Emit(READ_ACK, msg)
		return
	}

	if id, err := msg.GetI32("id"); err == nil {
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

			// 如果没有连接的话,跳过本次读取
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
					go func(id int32, tick int32, event_emit *queen.EventEmiter) {
						// 延迟
						time.Sleep(time.Duration(tick) * time.Millisecond)

						msg2 := nson.Message{"id": nson.I32(id)}
						event_emit.Emit(READ, msg2)

					}(id, tick, context.Event_emiter)
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

	context.Event_emiter.Emit(READ_ACK, msg)
}

func write(context queen.Context) {
	msg, ok := context.Message.(nson.Message)
	if !ok {
		msg.Insert("ok", nson.Bool(false))
		msg.Insert("error", nson.String("Message format error: message isn't nson!"))
		context.Event_emiter.Emit(WRITE_ACK, msg)
		return
	}

	if id, err := msg.GetI32("id"); err == nil {
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

	context.Event_emiter.Emit(WRITE_ACK, msg)
}

func close(context queen.Context) {
	msg, ok := context.Message.(nson.Message)
	if !ok {
		msg.Insert("ok", nson.Bool(false))
		msg.Insert("error", nson.String("Message format error: message isn't nson!"))
		context.Event_emiter.Emit(CLOSE_ACK, msg)
		return
	}

	if id, err := msg.GetI32("id"); err == nil {
		lock.RLock()
		if conn, ok := conns[id]; ok {
			conn.lock.Lock()
			conn.client.Close()
			conn.lock.Unlock()

			delete(conns, id)
		} else {
			msg.Insert("ok", nson.Bool(false))
			msg.Insert("error", nson.String("Message format error: can't get conn!"))
		}
		lock.RUnlock()
	} else {
		msg.Insert("ok", nson.Bool(false))
		msg.Insert("error", nson.String("Message format error: can't get id!"))
	}

	context.Event_emiter.Emit(CLOSE_ACK, msg)
}

func status(context queen.Context) {
	context.Event_emiter.Emit("driver/snap7/status/ack", nson.Message{"ok": nson.Bool(true)})
}
