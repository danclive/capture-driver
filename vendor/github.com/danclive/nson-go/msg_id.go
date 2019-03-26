package nson

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"sync/atomic"
	"time"
)

func (self MessageId) Tag() uint8 {
	return TAG_MESSAGE_ID
}

func (self MessageId) String() string {
	return fmt.Sprintf("MessageId(%x)", []byte(self))
}

func (self MessageId) Encode(buf *bytes.Buffer) error {
	if len(self) != 12 {
		return fmt.Errorf("MessageId must be 12 bytes: %v", self)
	}

	if _, err := buf.Write(self); err != nil {
		return err
	}

	return nil
}

func (self MessageId) Decode(buf *bytes.Buffer) (Value, error) {
	b := make([]byte, 12)
	_, err := io.ReadFull(buf, b)
	if err != nil {
		return nil, err
	}

	return MessageId(b), nil
}

var lastCount int32 = int32(time.Now().UnixNano())

// Create unique incrementing MessageId.
//
//   +---+---+---+---+---+---+---+---+---+---+---+---+
//   |       A       |     B     |   C   |     D     |
//   +---+---+---+---+---+---+---+---+---+---+---+---+
//     0   1   2   3   4   5   6   7   8   9  10  11
//   A = unix time (big endian), B = machine ID (first 3 bytes of md5 host name),
//   C = PID, D = incrementing counter (big endian)
func NewMessageId() (MessageId, error) {
	buf := new(bytes.Buffer)

	// A
	if err := binary.Write(buf, binary.BigEndian, uint32(time.Now().Unix())); err != nil {
		return nil, err
	}

	// B
	name, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	hash := md5.New()
	if _, err := hash.Write([]byte(name)); err != nil {
		return nil, err
	}
	if _, err := buf.Write(hash.Sum(nil)[:3]); err != nil {
		return nil, err
	}

	// C
	if err := binary.Write(buf, binary.BigEndian, int16(os.Getpid())); err != nil {
		return nil, err
	}

	// D
	cnt := atomic.AddInt32(&lastCount, 1) % 16777215
	cntbuf := make([]byte, 4)
	binary.BigEndian.PutUint32(cntbuf, uint32(cnt))
	if _, err := buf.Write(cntbuf[1:]); err != nil {
		return nil, err
	}

	return MessageId(buf.Bytes()), nil
}

func MessageIdFromHex(s string) (MessageId, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}

	return MessageId(b), nil
}

func (self MessageId) Hex() string {
	return hex.EncodeToString([]byte(self))
}

func (self MessageId) Timestamp() int64 {
	return int64(binary.BigEndian.Uint32([]byte(self)[:4]))
}

func (self MessageId) Time() time.Time {
	return time.Unix(self.Timestamp(), 0).UTC()
}
