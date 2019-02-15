package modbus

import (
	"github.com/danclive/nson-go"
	queen "github.com/danclive/queen-go"
)

func InitModbus(queen *queen.Queen) {
	doc := nson.Message{"driver": nson.String("modbus"), "ok": nson.Bool(true)}

	queen.Emit("init/driver/ack", doc)
}
