package modbus

import (
	"strconv"
	"strings"
)

type MBAddr struct {
	Area    string
	Name    string
	Bit     int
	Address int
	Format  string
	Size    int
}

var type_size map[string]int = map[string]int{
	"BOOL":    1,
	"BOOLEAN": 1,
	"BYTE":    1,
	"CHAR":    1,
	"SHORT":   2,
	"USHORT":  2,
	"INT16":   2,
	"UINT16":  2,
	"WORD":    2,
	"DWORD":   4,
	"INT":     4,
	"UINT":    4,
	"INT32":   4,
	"UINT32":  4,
	"FLOAT":   4,
	"REAL":    4,
	"REAL32":  4,
	"REAL64":  8,
	"LONG":    8,
	"INT64":   8,
	"UINT64":  8,
	"DOUBLE":  8,
	"STRING":  1,
}

func ParseTagAddress(addr string, raw_format string) MBAddr {
	var rs MBAddr
	rs.Name = addr
	var subaddress string
	rs.Area = addr[0:1]
	subaddress = addr[1:]
	rs.Format = strings.ToUpper(raw_format)
	rs.Address, rs.Bit = ParseTagSubAddress(subaddress)
	rs.Size = type_size[rs.Format]
	return rs
}

func ParseTagSubAddress(subaddr string) (int, int) {
	bit_start := 0

	for i, s := range subaddr {
		if s == '.' {
			bit_start = i + 1
		}

	}
	address := -1
	bit := -1

	if bit_start == 0 {
		address, _ = strconv.Atoi(subaddr)
	} else {
		address, _ = strconv.Atoi(subaddr[0 : bit_start-1])
		bit, _ = strconv.Atoi(subaddr[bit_start:])

	}

	return address - 1, bit
}
