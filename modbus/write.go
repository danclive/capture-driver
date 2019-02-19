package modbus

import (
	nson "github.com/danclive/nson-go"
	queen "github.com/danclive/queen-go"
)

func write_ext(context queen.Context, msg *nson.Message, conn *conn_t, tags nson.Array) {
	/*
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

			mbaddr := ParseTagAddress(address, format)

			var data []byte
			switch format {
			case "BYTE", "CHAR":
				if value, err := tag.GetI32("value"); err == nil {
					data = []byte{uint8(value)}
				}
			case "SHORT", "WORD":
				if value, err := tag.GetI32("value"); err == nil {
					data = util.IntTo2Bytes(int(value), conn.isbig)
				}
			case "LONG", "DWORD":
				if value, err := tag.GetI32("value"); err == nil {
					data = util.IntTo4Bytes(int(value), conn.isbig)
				}
			case "FLOAT", "REAL":
				if value, err := tag.GetF32("value"); err == nil {
					data = util.Float32To4Bytes(float32(value), conn.isbig)
				}
			case "BOOL":
				if value, err := tag.GetBool("value"); err == nil {
					var bits []byte
					var err error

					if mbaddr.Area == "0" {
						bits, err = conn.client.ReadCoils(uint16(mbaddr.Address), uint16(mbaddr.Size))
					}

					// todo

					if err != nil {
						conn.islink = false
						// 连接断开时，尝试重新连接
						context.Queen.Emit(
							RECONNECT,
							nson.Message{"id": nson.I32(conn.id), "retry": nson.I32(conn.retry)})
					}
				}
			}

			_ = data
		}

		msg.Insert("ok", nson.Bool(true))
	*/
}
