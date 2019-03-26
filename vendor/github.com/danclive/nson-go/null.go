package nson

import "bytes"

func (self Null) Tag() uint8 {
	return TAG_NULL
}

func (self Null) String() string {
	return "Null"
}

func (self Null) Encode(buf *bytes.Buffer) error {
	return nil
}

func (self Null) Decode(buf *bytes.Buffer) (Value, error) {
	return Null{}, nil
}
