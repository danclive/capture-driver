package capturedriver

import (
	"github.com/danclive/capture-driver/modbus"
	"github.com/danclive/capture-driver/snap7"
	queen "github.com/danclive/queen-go"
)

var (
	INIT     = "init/driver"
	INIT_ACK = "init/driver/ack"
)

func InitDriver(queen *queen.Queen) {
	snap7.InitSnap7(queen)
	modbus.InitModbus(queen)
}
