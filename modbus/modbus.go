package modbus

import (
	"github.com/danclive/nson-go"
	queen "github.com/danclive/queen-go"
)

func InitModbus(event_emiter *queen.EventEmiter) {
	doc := nson.Message{"driver": nson.String("modbus"), "ok": nson.Bool(true)}

	event_emiter.Emit("init/driver/ack", doc)
}
