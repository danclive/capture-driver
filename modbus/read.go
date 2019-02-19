package modbus

import (
	"github.com/danclive/capture-driver/util"
	nson "github.com/danclive/nson-go"
	queen "github.com/danclive/queen-go"
)

func read_ext(context queen.Context, msg *nson.Message, conn *conn_t, tags nson.Array) {
	tags_array := make(nson.Array, 0)

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

		var bits []byte

		switch mbaddr.Area {
		case "0":
			bits, err = conn.client.ReadCoils(uint16(mbaddr.Address), uint16(mbaddr.Size))
		case "1":
			bits, err = conn.client.ReadDiscreteInputs(uint16(mbaddr.Address), uint16(mbaddr.Size))
		case "3":
			bits, err = conn.client.ReadInputRegisters(uint16(mbaddr.Address), uint16((mbaddr.Size+1)/2))
		case "4":
			bits, err = conn.client.ReadHoldingRegisters(uint16(mbaddr.Address), uint16((mbaddr.Size+1)/2))
		}

		if err != nil {

			util.ErrorWithFields("read modbus err: ",
				util.Fields{
					"err":     err,
					"tag":     name,
					"address": address,
					"type":    ntype,
					"format":  format,
					"size":    mbaddr.Size,
				})

			conn.islink = false
			// 连接断开时，尝试重新连接
			context.Queen.Emit(
				RECONNECT,
				nson.Message{"id": nson.I32(conn.id), "retry": nson.I32(conn.retry)})
			return
		}

		var value nson.Value

		switch mbaddr.Area {
		case "0", "1":
			if len(bits) > 0 {
				v := false

				if bits[0] > 0 {
					v = true
				}

				value = nson.Bool(v)
			}
		case "3", "4":
			v, o := convert_raw_to_nson(conn, name, mbaddr, bits)
			if o {
				value = v
			}
		}

		if value != nil {
			rtag := nson.Message{
				"name":    nson.String(name),
				"type":    nson.String(ntype),
				"format":  nson.String(format),
				"address": nson.String(address),
				"value":   value,
			}

			tags_array = append(tags_array, rtag)
		}
	}

	msg.Insert("tags", nson.Array(tags_array))
	msg.Insert("ok", nson.Bool(true))
}
