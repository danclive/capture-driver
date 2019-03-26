package nson

import (
	"bytes"
	"fmt"
	"time"
)

// Timestamp
func (self Timestamp) Tag() uint8 {
	return TAG_TIMESTAMP
}

func (self Timestamp) String() string {
	return fmt.Sprintf("Timestamp(%v)", int64(self))
}

func (self Timestamp) Encode(buf *bytes.Buffer) error {
	return writeInt64(buf, int64(self))
}

func (self Timestamp) Decode(buf *bytes.Buffer) (Value, error) {
	v, err := readInt64(buf)
	if err != nil {
		return nil, err
	}

	return Timestamp(v), nil
}

// UTC datatime
func (self UTCDateTime) Tag() uint8 {
	return TAG_UTC_DATETIME
}

func (self UTCDateTime) String() string {
	return fmt.Sprintf("UTCDateTime(%v)", time.Time(self))
}

func (self UTCDateTime) Encode(buf *bytes.Buffer) error {
	t := time.Time(self).UnixNano() / 1000 / 1000

	return writeInt64(buf, t)
}

func (self UTCDateTime) Decode(buf *bytes.Buffer) (Value, error) {
	v, err := readInt64(buf)
	if err != nil {
		return nil, err
	}

	return UTCDateTime(time.Unix(v/1000, (v%1000)*1000000).UTC()), nil
}
