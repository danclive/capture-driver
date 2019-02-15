package snap7

import (
	"github.com/danclive/capture-driver/util"
	nson "github.com/danclive/nson-go"
	"github.com/danclive/queen-go"
)

func write_ext(context queen.Context, msg *nson.Message, conn *conn_t, tags nson.Array) {
	for _, tag_i := range tags {
		tag, ok := tag_i.(nson.Message)
		if !ok {
			continue
		}

		name, err := tag.GetString("name")
		if err != nil {
			continue
		}

		ntype, err := tag.GetString("type")
		if err != nil {
			continue
		}

		format, err := tag.GetString("format")
		if err != nil {
			continue
		}

		address, err := tag.GetString("address")
		if err != nil {
			continue
		}

		_ = name
		_ = ntype

		s7addr := ParseTagAddress(address)
		//data := make([]byte, 0)
		var data []byte

		switch format {
		case "BYTE", "CHAR":
			if value, err := tag.GetI32("value"); err == nil {
				data = []byte{uint8(value)}
			}
		case "SHORT", "WORD":
			if value, err := tag.GetI32("value"); err == nil {
				data = util.IntTo2Bytes(int(value))
			}
		case "LONG", "DWORD":
			if value, err := tag.GetI32("value"); err == nil {
				data = util.IntTo4Bytes(int(value))
			}
		case "FLOAT", "REAL":
			if value, err := tag.GetF32("value"); err == nil {
				data = util.Float32To4Bytes(float32(value))
			}
		case "BOOL":
			if value, err := tag.GetBool("value"); err == nil {
				bits, err := conn.client.ReadArea(s7addr.Cmd[0], s7addr.Cmd[1], s7addr.Cmd[2], s7addr.Cmd[3])
				if err != nil {
					conn.islink = false
					// 连接断开时，尝试重新连接
					context.Event_emiter.Emit(
						RECONNECT,
						nson.Message{"id": nson.I32(conn.id), "retry": nson.I32(conn.retry)})
				}

				switch s7addr.Format {
				case "B", "C", "X", "I", "DBB", "DBC", "DBX":
					if len(bits) != 1 {

					} else {
						bits[0] = util.SetBitFromBites(bits[0], s7addr.Bit, value)
						data = bits
					}
				case "W", "DBW":
					if len(bits) != 2 {

					} else {
						i := s7addr.Bit / 8
						bit := s7addr.Bit % 8

						if i <= 1 {
							bits[i] = util.SetBitFromBites(bits[i], bit, value)
						}
					}
				case "D", "DBD", "DI", "DBDI", "REAL", "DBREAL":
					if len(bits) != 4 {
						i := s7addr.Bit / 8
						bit := s7addr.Bit % 8

						if 1 <= 3 {
							bits[i] = util.SetBitFromBites(bits[i], bit, value)
						}
					}
				}
			}
		}

		if len(data) > 0 {
			err := conn.client.WriteArea(s7addr.Cmd[0], s7addr.Cmd[1], s7addr.Cmd[2], data)
			if err != nil {
				conn.islink = false
				// 连接断开时，尝试重新连接
				context.Event_emiter.Emit(
					RECONNECT,
					nson.Message{"id": nson.I32(conn.id), "retry": nson.I32(conn.retry)})
			}

			msg.Insert("ok", nson.Bool(true))
		}
	}
}
