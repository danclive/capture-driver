package main

import (
	"fmt"
	"time"

	"github.com/danclive/capture-driver"
	"github.com/danclive/capture-driver/snap7"
	nson "github.com/danclive/nson-go"
	queen "github.com/danclive/queen-go"
)

func main() {
	run()
	block()
}

func run() {
	event_emiter := queen.NewEventEmiter()

	event_emiter.On("run", func(context queen.Context) {
		event_emiter.Emit("init/driver", nil)
	})

	event_emiter.On("init/driver/ack", func(context queen.Context) {
		fmt.Printf("event: %s, message: %s\n", context.Event, context.Message)

		message := context.Message.(nson.Message)

		driver, err := message.GetString("driver")
		if err != nil {
			return
		}

		if driver == "snap7" {
			context.Event_emiter.Emit("driver/snap7/connect", nson.Message{
				"id":     nson.I32(123),
				"config": nson.String("S7-TCP://127.0.0.1?rank=0&slot=0&isBIGEndian=true"),
				"retry":  nson.I32(1),
			})
		} else if driver == "modbus" {

		}
	})

	event_emiter.On("init/driver", func(context queen.Context) {
		capturedriver.InitDriver(&event_emiter)
	})

	event_emiter.On(snap7.CONNECT_ACK, func(context queen.Context) {
		fmt.Println(context)

		msg, ok := context.Message.(nson.Message)
		if !ok {

		}

		if ok, _ := msg.GetBool("ok"); ok {
			fmt.Println(msg)
			read_msg(context.Event_emiter)
		}
	})

	event_emiter.On(snap7.RECONNECT_ACK, func(context queen.Context) {
		fmt.Println(context)

		msg, ok := context.Message.(nson.Message)
		if !ok {

		}

		if ok, _ := msg.GetBool("ok"); ok {
			fmt.Println(msg)
			read_msg(context.Event_emiter)
		}
	})

	event_emiter.On(snap7.READ_ACK, func(context queen.Context) {
		fmt.Println(context)

		if msg, ok := context.Message.(nson.Message); ok {
			if ok, _ := msg.GetBool("ok"); ok {
				// fmt.Println(msg)

				if tags, err := msg.GetArray("tags"); err == nil {
					tags2 := make(nson.Array, 0)

					for _, tag_i := range tags {
						tag, ok := tag_i.(nson.Message)
						if !ok {
							continue
						}

						value, err := tag.GetI32("value")
						if err != nil {
							continue
						}

						fmt.Println("//////", value)
						value += 1

						tag.Insert("value", nson.I32(value))

						tags2 = append(tags2, tag)
					}

					msg.Insert("tags", tags2)
				}

				context.Event_emiter.Emit(snap7.WRITE, msg)
			}
		}
	})

	event_emiter.Emit("run", nil)
}

func read_msg(event_emiter *queen.EventEmiter) {
	tags := make(nson.Array, 0)

	tags = append(tags, nson.Message{
		"name":    nson.String("tag1"),
		"type":    nson.String("INT"),
		"format":  nson.String("WORD"),
		"address": nson.String("DB1,W0"),
	})

	msg := nson.Message{"id": nson.I32(123), "tags": tags, "tick": nson.I32(2)}

	event_emiter.Emit(snap7.READ, msg)
}

// 连接 driver/snap7/connect driver/snap7/connect/ack
// { "id": 123, "config": "S7-TCP://127.0.0.1?rank=0&slot=0&isBIGEndian=true", "options": {}}
// 读 driver/snap7/read driver/snap7/read/ack // 紧急数据 {"urgent": true}
// { "id": 123, tags: [{"name": "name", "type": "", "format": "", "address": "", "value": ""}] }
// 能够定时读, tick ?
// 写 driver/snap7/write driver/snap7/write/ack
// { "id": 123, tags: [{"name": "name", "type": "", "format": "", "address": "", "value": ""}] }
// 关闭 driver/snap7/close driver/snap7/close/ack
// { "id": 123 }
// 状态 driver/snap7/status driver/snap7/status/ack
// { "id": 123 }

func block() {

	for {
		time.Sleep(1 * time.Second)
	}
}
