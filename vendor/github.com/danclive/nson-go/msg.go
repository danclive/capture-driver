package nson

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"
)

// Message
func (self Message) Tag() uint8 {
	return TAG_MESSAGE
}

func (self Message) String() string {
	buf := new(bytes.Buffer)

	fmt.Fprintf(buf, "Message{")

	l := len(self)
	i := 0

	for k, v := range self {
		i++

		if i == l {
			fmt.Fprintf(buf, "%v: %v", k, v.String())
		} else {
			fmt.Fprintf(buf, "%v: %v, ", k, v.String())
		}
	}

	fmt.Fprintf(buf, "}")

	return buf.String()
}

func (self Message) Encode(buff *bytes.Buffer) error {
	buf := new(bytes.Buffer)

	if err := writeInt32(buf, 0); err != nil {
		return err
	}

	for k, v := range self {
		if err := buf.WriteByte(v.Tag()); err != nil {
			return err
		}

		if err := writeCstring(buf, k); err != nil {
			return err
		}

		if err := v.Encode(buf); err != nil {
			return err
		}
	}

	if err := buf.WriteByte(0x00); err != nil {
		return err
	}

	binary.LittleEndian.PutUint32(buf.Bytes(), uint32(buf.Len()))

	if _, err := buf.WriteTo(buff); err != nil {
		return err
	}

	return nil
}

func (self Message) Decode(buf *bytes.Buffer) (Value, error) {
	_, err := readInt32(buf)
	if err != nil {
		return nil, err
	}

	msg := Message{}

	for {
		tag, err := buf.ReadByte()
		if err != nil {
			return nil, err
		}

		if tag == 0 {
			break
		}

		key, value, err := decode_value(buf, tag)
		if err != nil {
			return nil, err
		}

		msg[key] = value
	}

	return msg, nil
}

func (self *Message) Get(key string) (Value, bool) {
	value, has := (*self)[key]
	return value, has
}

func (self *Message) Contains(key string) bool {
	_, has := (*self)[key]
	return has
}

func (self *Message) Len() int {
	return len(*self)
}

func (self *Message) Insert(key string, value Value) {
	(*self)[key] = value
}

func (self *Message) Remove(key string) bool {
	has := self.Contains(key)

	if has {
		delete(*self, key)
	}

	return has
}

func (self *Message) GetF32(key string) (float32, error) {
	value, has := self.Get(key)
	if !has {
		return 0, fmt.Errorf("Not Present, key: %v", key)
	}

	if value.Tag() != TAG_F32 {
		return 0, fmt.Errorf("Unexpected Type, key: %v, value: %v", key, value)
	}

	return float32(value.(F32)), nil
}

func (self *Message) GetF64(key string) (float64, error) {
	value, has := self.Get(key)
	if !has {
		return 0, fmt.Errorf("Not Present, key: %v", key)
	}

	if value.Tag() != TAG_F64 {
		return 0, fmt.Errorf("Unexpected Type, key: %v, value: %v", key, value)
	}

	return float64(value.(F64)), nil
}

func (self *Message) GetI32(key string) (int32, error) {
	value, has := self.Get(key)
	if !has {
		return 0, fmt.Errorf("Not Present, key: %v", key)
	}

	if value.Tag() != TAG_I32 {
		return 0, fmt.Errorf("Unexpected Type, key: %v, value: %v", key, value)
	}

	return int32(value.(I32)), nil
}

func (self *Message) GetU32(key string) (uint32, error) {
	value, has := self.Get(key)
	if !has {
		return 0, fmt.Errorf("Not Present, key: %v", key)
	}

	if value.Tag() != TAG_U32 {
		return 0, fmt.Errorf("Unexpected Type, key: %v, value: %v", key, value)
	}

	return uint32(value.(U32)), nil
}

func (self *Message) GetI64(key string) (int64, error) {
	value, has := self.Get(key)
	if !has {
		return 0, fmt.Errorf("Not Present, key: %v", key)
	}

	if value.Tag() != TAG_I64 {
		return 0, fmt.Errorf("Unexpected Type, key: %v, value: %v", key, value)
	}

	return int64(value.(I64)), nil
}

func (self *Message) GetU64(key string) (uint64, error) {
	value, has := self.Get(key)
	if !has {
		return 0, fmt.Errorf("Not Present, key: %v", key)
	}

	if value.Tag() != TAG_U64 {
		return 0, fmt.Errorf("Unexpected Type, key: %v, value: %v", key, value)
	}

	return uint64(value.(U64)), nil
}

func (self *Message) GetString(key string) (string, error) {
	value, has := self.Get(key)
	if !has {
		return "", fmt.Errorf("Not Present, key: %v", key)
	}

	if value.Tag() != TAG_STRING {
		return "", fmt.Errorf("Unexpected Type, key: %v, value: %v", key, value)
	}

	return string(value.(String)), nil
}

func (self *Message) GetArray(key string) (Array, error) {
	value, has := self.Get(key)
	if !has {
		return nil, fmt.Errorf("Not Present, key: %v", key)
	}

	if value.Tag() != TAG_ARRAY {
		return nil, fmt.Errorf("Unexpected Type, key: %v, value: %v", key, value)
	}

	return value.(Array), nil
}

func (self *Message) GetMessage(key string) (Message, error) {
	value, has := self.Get(key)
	if !has {
		return nil, fmt.Errorf("Not Present, key: %v", key)
	}

	if value.Tag() != TAG_MESSAGE {
		return nil, fmt.Errorf("Unexpected Type, key: %v, value: %v", key, value)
	}

	return value.(Message), nil
}

func (self *Message) GetBool(key string) (bool, error) {
	value, has := self.Get(key)
	if !has {
		return false, fmt.Errorf("Not Present, key: %v", key)
	}

	if value.Tag() != TAG_BOOLEAN {
		return false, fmt.Errorf("Unexpected Type, key: %v, value: %v", key, value)
	}

	return bool(value.(Bool)), nil
}

func (self *Message) IsNull(key string) bool {
	value, has := self.Get(key)
	if !has {
		return false
	}

	if value.Tag() != TAG_NULL {
		return false
	}

	return true
}

func (self *Message) GetBinary(key string) ([]byte, error) {
	value, has := self.Get(key)
	if !has {
		return nil, fmt.Errorf("Not Present, key: %v", key)
	}

	if value.Tag() != TAG_BINARY {
		return nil, fmt.Errorf("Unexpected Type, key: %v, value: %v", key, value)
	}

	return []byte(value.(Binary)), nil
}

func (self *Message) GetTimestamp(key string) (int64, error) {
	value, has := self.Get(key)
	if !has {
		return 0, fmt.Errorf("Not Present, key: %v", key)
	}

	if value.Tag() != TAG_TIMESTAMP {
		return 0, fmt.Errorf("Unexpected Type, key: %v, value: %v", key, value)
	}

	return int64(value.(Timestamp)), nil
}

func (self *Message) GetUTCDateTime(key string) (time.Time, error) {
	value, has := self.Get(key)
	if !has {
		return time.Time{}, fmt.Errorf("Not Present, key: %v", key)
	}

	if value.Tag() != TAG_UTC_DATETIME {
		return time.Time{}, fmt.Errorf("Unexpected Type, key: %v, value: %v", key, value)
	}

	return time.Time(value.(UTCDateTime)), nil
}

func (self *Message) GetMessageId(key string) (MessageId, error) {
	value, has := self.Get(key)
	if !has {
		return nil, fmt.Errorf("Not Present, key: %v", key)
	}

	if value.Tag() != TAG_MESSAGE_ID {
		return nil, fmt.Errorf("Unexpected Type, key: %v, value: %v", key, value)
	}

	return value.(MessageId), nil
}
