package modbus

import (
	"github.com/danclive/capture-driver/util"
	nson "github.com/danclive/nson-go"
)

func convert_raw_to_nson(conn *conn_t, tag string, mdaddr MBAddr, data []byte) (nson.Value, bool) {
	switch mdaddr.Format {
	case "BYTE":
		if len(data) != 1 {
			util.Error("Convert data BYTE: len(tag.Data) != 1 ", tag)
		} else {
			u8 := uint8(data[0])

			return nson.I32(int32(u8)), true
		}
	case "CHAR":
		if len(data) != 1 {
			util.Error("Convert data CHAR: len(data) != 1 ", tag)
		} else {
			i8 := int8(data[0])

			return nson.I32(int32(i8)), true
		}
	case "SHORT":
		if len(data) != 2 {
			util.Error("Convert data SHORT: len(data) != 2 ", tag)
		} else {
			i16 := util.BytesToInt16(data, conn.isbig)

			return nson.I32(int32(i16)), true
		}
	case "WORD":
		if len(data) != 2 {
			util.Error("Convert data WORD: len(data) != 2 ", tag)
		} else {
			u16 := util.BytesToUInt16(data, conn.isbig)

			return nson.I32(int32(u16)), true
		}
	case "LONG":
		if len(data) != 4 {
			util.Error("Convert data LONG: len(data) != 4 ", tag)
		} else {
			i32 := util.BytesToInt32(data, conn.isbig)

			return nson.I32(int32(i32)), true
		}
	case "DWORD":
		if len(data) != 4 {
			util.Error("Convert data DWORD: len(data) != 4 ", tag)
		} else {
			u32 := util.BytesToUInt32(data, conn.isbig)

			return nson.I32(int32(u32)), true
		}
	case "FLOAT", "REAL":
		if len(data) != 4 {
			util.Error("Convert data FLOAT, REAL: len(data) != 4 ", tag)
		} else {
			f32 := util.BytesToFloat32(data, conn.isbig)

			return nson.F32(float32(f32)), true
		}
	case "BOOL":
		if len(data) != 2 {

		} else {
			u16 := util.BytesToUInt16(data, conn.isbig)

			value := false

			if (u16 & (1 << uint32(mdaddr.Bit))) > 0 {
				value = true
			}

			return nson.Bool(value), true
		}
	}

	return nil, false
}
