package nson

import (
	"bytes"
	"fmt"
)

// Float 32
func (self F32) Tag() uint8 {
	return TAG_F32
}

func (self F32) String() string {
	return fmt.Sprintf("F32(%v)", float32(self))
}

func (self F32) Encode(buf *bytes.Buffer) error {
	return writeFloat32(buf, float32(self))
}

func (self F32) Decode(buf *bytes.Buffer) (Value, error) {
	v, err := readFloat32(buf)
	if err != nil {
		return nil, err
	}

	return F32(v), nil
}

// Float 64
func (self F64) Tag() uint8 {
	return TAG_F64
}

func (self F64) String() string {
	return fmt.Sprintf("F32(%v)", float64(self))
}

func (self F64) Encode(buf *bytes.Buffer) error {
	return writeFloat64(buf, float64(self))
}

func (self F64) Decode(buf *bytes.Buffer) (Value, error) {
	v, err := readFloat64(buf)
	if err != nil {
		return nil, err
	}

	return F64(v), nil
}

// Int 32
func (self I32) Tag() uint8 {
	return TAG_I32
}

func (self I32) String() string {
	return fmt.Sprintf("I32(%v)", int32(self))
}

func (self I32) Encode(buf *bytes.Buffer) error {
	return writeInt32(buf, int32(self))
}

func (self I32) Decode(buf *bytes.Buffer) (Value, error) {
	v, err := readInt32(buf)
	if err != nil {
		return nil, err
	}

	return I32(v), nil
}

// Int 64
func (self I64) Tag() uint8 {
	return TAG_I64
}

func (self I64) String() string {
	return fmt.Sprintf("I64(%v)", int64(self))
}

func (self I64) Encode(buf *bytes.Buffer) error {
	return writeInt64(buf, int64(self))
}

func (self I64) Decode(buf *bytes.Buffer) (Value, error) {
	v, err := readInt64(buf)
	if err != nil {
		return nil, err
	}

	return I64(v), nil
}

// Uint 32
func (self U32) Tag() uint8 {
	return TAG_U32
}

func (self U32) String() string {
	return fmt.Sprintf("U32(%v)", uint32(self))
}

func (self U32) Encode(buf *bytes.Buffer) error {
	return writeUint32(buf, uint32(self))
}

func (self U32) Decode(buf *bytes.Buffer) (Value, error) {
	v, err := readUint32(buf)
	if err != nil {
		return nil, err
	}

	return U32(v), nil
}

// Uint 64
func (self U64) Tag() uint8 {
	return TAG_U64
}

func (self U64) String() string {
	return fmt.Sprintf("U64(%v)", uint64(self))
}

func (self U64) Encode(buf *bytes.Buffer) error {
	return writeUint64(buf, uint64(self))
}

func (self U64) Decode(buf *bytes.Buffer) (Value, error) {
	v, err := readUint64(buf)
	if err != nil {
		return nil, err
	}

	return U64(v), nil
}
