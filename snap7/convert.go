package snap7

import (
	"github.com/danclive/capture-driver/util"
	"github.com/danclive/capture/log"
	nson "github.com/danclive/nson-go"
)

func convert_raw_to_nson(tag *tag_t2) (nson.Value, bool) {

	switch tag.tag.format {
	case "BYTE":
		if len(tag.tag.data) != 1 {
			log.Error("Convert data BYTE: len(tag.Data) != 1 ", tag)
		} else {
			u8 := uint8(tag.tag.data[0])

			return nson.I32(int32(u8)), true
		}
	case "CHAR":
		if len(tag.tag.data) != 1 {
			log.Error("Convert data CHAR: len(tag.tag.data) != 1 ", tag)
		} else {
			i8 := int8(tag.tag.data[0])

			return nson.I32(int32(i8)), true
		}
	case "SHORT":
		if len(tag.tag.data) != 2 {
			log.Error("Convert data SHORT: len(tag.tag.data) != 2 ", tag)
		} else {
			i16 := util.BytesToInt16(tag.tag.data, true)

			return nson.I32(int32(i16)), true
		}
	case "WORD":
		if len(tag.tag.data) != 2 {
			log.Error("Convert data WORD: len(tag.tag.data) != 2 ", tag)
		} else {
			u16 := util.BytesToUInt16(tag.tag.data, true)

			return nson.I32(int32(u16)), true
		}
	case "LONG":
		if len(tag.tag.data) != 4 {
			log.Error("Convert data LONG: len(tag.tag.data) != 4 ", tag)
		} else {
			i32 := util.BytesToInt32(tag.tag.data, true)

			return nson.I32(int32(i32)), true
		}
	case "DWORD":
		if len(tag.tag.data) != 4 {
			log.Error("Convert data DWORD: len(tag.tag.data) != 4 ", tag)
		} else {
			u32 := util.BytesToUInt32(tag.tag.data, true)

			return nson.I32(int32(u32)), true
		}
	case "FLOAT", "REAL":
		if len(tag.tag.data) != 4 {
			log.Error("Convert data FLOAT, REAL: len(tag.tag.data) != 4 ", tag)
		} else {
			f32 := util.BytesToFloat32(tag.tag.data, true)

			return nson.F32(float32(f32)), true
		}
	case "BOOL":
		switch tag.s7addr.Format {
		case "B", "C", "X", "I", "DBB", "DBC", "DBX":
			if len(tag.tag.data) != 1 {
				log.Error("Convert data BOOL -> B: len(tag.tag.data) != 1 ", tag)
			} else {
				b := util.GetBitFromBites(tag.tag.data[0], tag.s7addr.Bit)

				return nson.Bool(b), true
			}
		case "W", "DBW":
			if len(tag.tag.data) != 2 {
				log.Error("Convert data BOOL -> W, X: len(tag.tag.data) != 2 ", tag)
			} else {
				i := tag.s7addr.Bit / 8
				bit := tag.s7addr.Bit % 8
				if i <= 1 {
					b := util.GetBitFromBites(tag.tag.data[i], bit)

					return nson.Bool(b), true
				} else {
					log.Error("Convert data BOOL -> W, X: i > 1 ", tag.s7addr)
				}
			}
		case "D", "DBD", "DI", "DBDI", "REAL", "DBREAL":
			if len(tag.tag.data) != 4 {
				log.Error("Convert data BOOL -> D: len(tag.tag.data) != 4 ", tag)
			} else {
				i := tag.s7addr.Bit / 8
				bit := tag.s7addr.Bit % 8
				if i <= 3 {
					b := util.GetBitFromBites(tag.tag.data[i], bit)

					return nson.Bool(b), true
				} else {
					log.Error("Convert data BOOL -> D: i > 3 ", tag.s7addr)
				}
			}
		}
	}

	return nil, false
}
