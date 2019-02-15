package snap7

import (
	"strconv"
	"strings"
)

func Substr(str string, start int, end int) string {
	return string([]rune(str)[start:end])
}

type S7Addr struct {
	Area    string
	Name    string
	Bit     int
	Address int
	Cmd     [4]int
	Format  string
	Size    int
}

var type_size map[string]int = map[string]int{
	"B":      1,
	"DBB":    1,
	"C":      1,
	"DBC":    1,
	"X":      1,
	"DBX":    1,
	"I":      2,
	"DBI":    2,
	"W":      2,
	"DBW":    2,
	"D":      4,
	"DBD":    4,
	"DI":     4,
	"DBDI":   4,
	"REAL":   4,
	"DBREAL": 4,
}

var S7Area map[string]int = map[string]int{
	"PE": 0x81,
	"PA": 0x82,
	"MK": 0x83,
	"DB": 0x84,
	"CT": 0x1C,
	"TM": 0x1D,
}

func ParseTagAddress(addr string) S7Addr {
	var rs S7Addr
	rs.Name = addr

	var subaddress string
	var DB_No = 0
	if strings.HasPrefix(addr, "DB") {
		rs.Area = "DB"
		//Parse DBnn,xxnn or DBnn.xxnn
		block_end := -1
		for i, s := range addr {
			if s == '.' || s == ',' {
				block_end = i
				break
			}
		}

		DB_No, _ = strconv.Atoi(addr[2:block_end])

		subaddress = addr[block_end+1:]

	} else if strings.HasPrefix(addr, "PI") || strings.HasPrefix(addr, "PE") {
		rs.Area = "PE"
		subaddress = addr[2:]
		//PE area
	} else if strings.HasPrefix(addr, "PQ") || strings.HasPrefix(addr, "PA") {
		rs.Area = "PA"
		subaddress = addr[2:]
		//PA area
	} else {
		if strings.HasPrefix(addr, "I") || strings.HasPrefix(addr, "E") {
			rs.Area = "PE"
		} else if strings.HasPrefix(addr, "Q") || strings.HasPrefix(addr, "A") {
			rs.Area = "PA"
		} else if strings.HasPrefix(addr, "M") || strings.HasPrefix(addr, "F") {
			rs.Area = "MK"
		} else if strings.HasPrefix(addr, "C") || strings.HasPrefix(addr, "Z") {
			rs.Area = "CT"
		} else if strings.HasPrefix(addr, "T") {
			rs.Area = "TM"
		}

		subaddress = addr[1:]

	}
	rs.Format, rs.Address, rs.Bit = ParseTagSubAddress(subaddress)
	rs.Size = type_size[rs.Format]
	rs.Cmd = [4]int{S7Area[rs.Area], DB_No, rs.Address, rs.Size}
	return rs
}

func ParseTagSubAddress(subaddr string) (string, int, int) {

	type_name := ""

	digit_start := -1
	bit_start := 0
	for i, s := range subaddr {
		if s >= 48 && s <= 57 && digit_start == -1 {
			digit_start = i
		}
		if s == '.' {
			bit_start = i + 1
		}

	}
	if digit_start == 0 {
		type_name = "X"
	} else {
		type_name = strings.ToUpper(subaddr[0:digit_start])
	}
	address := -1
	bit := -1

	if bit_start == 0 {
		address, _ = strconv.Atoi(subaddr[digit_start:])
	} else {
		address, _ = strconv.Atoi(subaddr[digit_start : bit_start-1])
		bit, _ = strconv.Atoi(subaddr[bit_start:])

	}

	return type_name, address, bit
}

/*
func main() {
	s := "DB1,X0.1"
	res := ParseTagAddress("DB1,X0.1")
	fmt.Println(res)
	fmt.Println(Substr("DB1,X0.1", 2, 6))
	fmt.Println(s[2:3])
	fmt.Println(ParseTagSubAddress("111.0"))
	fmt.Println(ParseTagSubAddress("B0.2"))
	fmt.Println(ParseTagSubAddress("REAL23"))
	fmt.Println(ParseTagAddress("DB1,X0.1"))
	fmt.Println(ParseTagAddress("I0.1"))
	fmt.Println(ParseTagAddress("PIW12"))
	fmt.Println(ParseTagAddress("MD23"))
	fmt.Println(ParseTagAddress("M2.7"))
	fmt.Println(ParseTagAddress("MW0.14"))
	fmt.Println(ParseTagAddress("MREAL44"))
	fmt.Println(ParseTagAddress("DB1,REAL24"))
}
*/
