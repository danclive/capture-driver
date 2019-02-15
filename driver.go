package capturedriver

import (
	"github.com/danclive/capture-driver/modbus"
	"github.com/danclive/capture-driver/snap7"
	queen "github.com/danclive/queen-go"
)

func InitDriver(event_emiter *queen.EventEmiter) {
	snap7.InitSnap7(event_emiter)
	modbus.InitModbus(event_emiter)
}
