package util

import (
	"bytes"
	"encoding/binary"
	"math"
)

func GetBitFromBites(b byte, bit int) bool {
	return (b & (1 << uint32(bit))) > 0
}

func SetBitFromBites(b byte, bit int, io bool) byte {
	if io {
		b = b | (1 << uint32(bit))
	} else {
		b = b &^ (1 << uint32(bit))
	}

	return b
}

// bytes to int 16
func BytesToInt16(b []byte) int16 {
	buf := bytes.NewBuffer(b)
	var tmp int16
	binary.Read(buf, binary.BigEndian, &tmp)
	return tmp
}

// bytes to uint 16
func BytesToUInt16(b []byte) uint16 {
	buf := bytes.NewBuffer(b)
	var tmp uint16
	binary.Read(buf, binary.BigEndian, &tmp)
	return tmp
}

// bytes to int 32
func BytesToInt32(b []byte) int32 {
	buf := bytes.NewBuffer(b)
	var tmp int32
	binary.Read(buf, binary.BigEndian, &tmp)
	return tmp
}

// bytes to uint 32
func BytesToUInt32(b []byte) uint32 {
	buf := bytes.NewBuffer(b)
	var tmp uint32
	binary.Read(buf, binary.BigEndian, &tmp)
	return tmp
}

// bytes to float32
func BytesToFloat32(b []byte) float32 {
	buf := bytes.NewBuffer(b)
	var tmp uint32
	binary.Read(buf, binary.BigEndian, &tmp)
	return math.Float32frombits(tmp)
}

// int to 4 bytes
func IntTo4Bytes(i int) []byte {
	buf := bytes.NewBuffer([]byte{})
	tmp := uint32(i)
	binary.Write(buf, binary.BigEndian, tmp)
	return buf.Bytes()
}

// int to 2 bytes
func IntTo2Bytes(i int) []byte {
	buf := bytes.NewBuffer([]byte{})
	tmp := uint16(i)
	binary.Write(buf, binary.BigEndian, tmp)
	return buf.Bytes()
}

// Float32 to 4 bytes
func Float32To4Bytes(f float32) []byte {
	bits := math.Float32bits(f)
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, bits)
	return bytes
}
