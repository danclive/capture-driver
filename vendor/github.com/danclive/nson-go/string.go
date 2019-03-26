package nson

import (
	"bytes"
	"fmt"
)

// String
func (self String) Tag() uint8 {
	return TAG_STRING
}

func (self String) String() string {
	return fmt.Sprintf("String(%v)", string(self))
}

func (self String) Encode(buf *bytes.Buffer) error {
	return writeString(buf, string(self))
}

func (self String) Decode(buf *bytes.Buffer) (Value, error) {
	v, err := readString(buf)
	if err != nil {
		return nil, err
	}

	return String(v), nil
}
