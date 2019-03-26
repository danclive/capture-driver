package nson

import (
	"bytes"
	"fmt"
	"io"
)

func (self Binary) Tag() uint8 {
	return TAG_BINARY
}

func (self Binary) String() string {
	return fmt.Sprintf("Binary(%v)", []byte(self))
}

func (self Binary) Encode(buf *bytes.Buffer) error {
	if err := writeInt32(buf, int32(len(self))); err != nil {
		return nil
	}

	if _, err := buf.Write(self); err != nil {
		return err
	}

	return nil
}

func (self Binary) Decode(buf *bytes.Buffer) (Value, error) {
	l, err := readInt32(buf)
	if err != nil {
		return nil, err
	}

	b := make([]byte, l)
	if _, err := io.ReadFull(buf, b); err != nil {
		return nil, err
	}

	return Binary(b), nil
}
