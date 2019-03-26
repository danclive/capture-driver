package nson

import (
	"bytes"
	"fmt"
)

func (self Bool) Tag() uint8 {
	return TAG_BOOLEAN
}

func (self Bool) String() string {
	if self {
		return fmt.Sprint("True")
	} else {
		return fmt.Sprint("False")
	}
}

func (self Bool) Encode(buf *bytes.Buffer) error {
	if self {
		if err := buf.WriteByte(0x01); err != nil {
			return err
		}
	} else {
		if err := buf.WriteByte(0x00); err != nil {
			return err
		}
	}

	return nil
}

func (self Bool) Decode(buf *bytes.Buffer) (Value, error) {
	b, err := buf.ReadByte()
	if err != nil {
		return nil, err
	}

	return Bool(b == 0x01), nil
}
