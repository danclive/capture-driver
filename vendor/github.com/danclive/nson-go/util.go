package nson

import (
	"bytes"
	"encoding/binary"
	"io"
)

func writeCstring(buf *bytes.Buffer, s string) error {
	if _, err := buf.WriteString(s); err != nil {
		return err
	}
	return buf.WriteByte(0x00)
}

func writeString(buf *bytes.Buffer, s string) error {
	if err := binary.Write(buf, binary.LittleEndian, uint32(len(s)+1)); err != nil {
		return err
	}
	if _, err := buf.WriteString(s); err != nil {
		return err
	}
	return buf.WriteByte(0x00)
}

func writeFloat32(buf *bytes.Buffer, f float32) error {
	return binary.Write(buf, binary.LittleEndian, f)
}

func writeFloat64(buf *bytes.Buffer, f float64) error {
	return binary.Write(buf, binary.LittleEndian, f)
}

func writeInt32(buf *bytes.Buffer, i int32) error {
	return binary.Write(buf, binary.LittleEndian, i)
}

func writeInt64(buf *bytes.Buffer, i int64) error {
	return binary.Write(buf, binary.LittleEndian, i)
}

func writeUint32(buf *bytes.Buffer, u uint32) error {
	return binary.Write(buf, binary.LittleEndian, u)
}

func writeUint64(buf *bytes.Buffer, u uint64) error {
	return binary.Write(buf, binary.LittleEndian, u)
}

func readCstring(rd *bytes.Buffer) (string, error) {
	s, err := rd.ReadString(0x00)
	if err != nil {
		return "", err
	}
	return s[:len(s)-1], nil
}

func readString(rd *bytes.Buffer) (string, error) {
	// Read string length.
	var sLen int32
	if err := binary.Read(rd, binary.LittleEndian, &sLen); err != nil {
		return "", err
	}
	if sLen == 0 {
		return "", nil
	}

	// Read string.
	b := make([]byte, sLen)
	if _, err := io.ReadFull(rd, b); err != nil {
		return "", err
	}
	return string(b[:len(b)-1]), nil
}

func readFloat32(rd *bytes.Buffer) (float32, error) {
	var f float32
	if err := binary.Read(rd, binary.LittleEndian, &f); err != nil {
		return 0, err
	}
	return f, nil
}

func readFloat64(rd *bytes.Buffer) (float64, error) {
	var f float64
	if err := binary.Read(rd, binary.LittleEndian, &f); err != nil {
		return 0, err
	}
	return f, nil
}

func readInt32(rd *bytes.Buffer) (int32, error) {
	var i int32
	if err := binary.Read(rd, binary.LittleEndian, &i); err != nil {
		return 0, err
	}
	return i, nil
}

func readInt64(rd *bytes.Buffer) (int64, error) {
	var i int64
	if err := binary.Read(rd, binary.LittleEndian, &i); err != nil {
		return 0, err
	}
	return i, nil
}

func readUint32(rd *bytes.Buffer) (uint32, error) {
	var u uint32
	if err := binary.Read(rd, binary.LittleEndian, &u); err != nil {
		return 0, err
	}
	return u, nil
}

func readUint64(rd *bytes.Buffer) (uint64, error) {
	var u uint64
	if err := binary.Read(rd, binary.LittleEndian, &u); err != nil {
		return 0, err
	}
	return u, nil
}
