package nson

import (
	"bytes"
	"fmt"
)

func decode_value(buf *bytes.Buffer, tag uint8) (string, Value, error) {
	switch tag {
	case TAG_F32:
		key, err := readCstring(buf)

		if err != nil {
			return "", nil, err
		}

		value, err := F32(0).Decode(buf)
		if err != nil {
			return "", nil, err
		}

		return key, value, nil
	case TAG_F64:
		key, err := readCstring(buf)

		if err != nil {
			return "", nil, err
		}

		value, err := F64(0).Decode(buf)
		if err != nil {
			return "", nil, err
		}

		return key, value, nil
	case TAG_I32:
		key, err := readCstring(buf)

		if err != nil {
			return "", nil, err
		}

		value, err := I32(0).Decode(buf)
		if err != nil {
			return "", nil, err
		}

		return key, value, nil
	case TAG_I64:
		key, err := readCstring(buf)

		if err != nil {
			return "", nil, err
		}

		value, err := I64(0).Decode(buf)
		if err != nil {
			return "", nil, err
		}

		return key, value, nil
	case TAG_U32:
		key, err := readCstring(buf)

		if err != nil {
			return "", nil, err
		}

		value, err := U32(0).Decode(buf)
		if err != nil {
			return "", nil, err
		}

		return key, value, nil
	case TAG_U64:
		key, err := readCstring(buf)

		if err != nil {
			return "", nil, err
		}

		value, err := U64(0).Decode(buf)
		if err != nil {
			return "", nil, err
		}

		return key, value, nil
	case TAG_STRING:
		key, err := readCstring(buf)

		if err != nil {
			return "", nil, err
		}

		value, err := String("").Decode(buf)
		if err != nil {
			return "", nil, err
		}

		return key, value, nil
	case TAG_ARRAY:
		key, err := readCstring(buf)

		if err != nil {
			return "", nil, err
		}

		value, err := Array{}.Decode(buf)
		if err != nil {
			return "", nil, err
		}

		return key, value, nil
	case TAG_BOOLEAN:
		key, err := readCstring(buf)

		if err != nil {
			return "", nil, err
		}

		value, err := Bool(false).Decode(buf)
		if err != nil {
			return "", nil, err
		}

		return key, value, nil
	case TAG_NULL:
		key, err := readCstring(buf)

		if err != nil {
			return "", nil, err
		}

		value, err := Null{}.Decode(buf)
		if err != nil {
			return "", nil, err
		}

		return key, value, nil
	case TAG_BINARY:
		key, err := readCstring(buf)

		if err != nil {
			return "", nil, err
		}

		value, err := Binary{}.Decode(buf)
		if err != nil {
			return "", nil, err
		}

		return key, value, nil
	case TAG_TIMESTAMP:
		key, err := readCstring(buf)

		if err != nil {
			return "", nil, err
		}

		value, err := Timestamp(0).Decode(buf)
		if err != nil {
			return "", nil, err
		}

		return key, value, nil
	case TAG_UTC_DATETIME:
		key, err := readCstring(buf)

		if err != nil {
			return "", nil, err
		}

		value, err := UTCDateTime{}.Decode(buf)
		if err != nil {
			return "", nil, err
		}

		return key, value, nil
	case TAG_MESSAGE_ID:
		key, err := readCstring(buf)

		if err != nil {
			return "", nil, err
		}

		value, err := MessageId{}.Decode(buf)
		if err != nil {
			return "", nil, err
		}

		return key, value, nil
	case TAG_MESSAGE:
		key, err := readCstring(buf)

		if err != nil {
			return "", nil, err
		}

		value, err := Message{}.Decode(buf)
		if err != nil {
			return "", nil, err
		}

		return key, value, nil
	default:
		return "", nil, fmt.Errorf("Unsupported type '%X'.", tag)
	}
}
