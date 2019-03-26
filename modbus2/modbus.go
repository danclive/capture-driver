package main

import (
	"fmt"
	"sort"

	modbus2 "github.com/danclive/capture-driver/modbus"
	"github.com/danclive/nson-go"
	"github.com/goburrow/modbus"
)

func main() {
	for true {
		aaa()
	}
}

func aaa() {
	handler := modbus.NewTCPClientHandler("ciiwater.com:27899")
	handler.SlaveId = 1

	err := handler.Connect()
	fmt.Println(err)
	defer handler.Close()

	client := modbus.NewClient(handler)

	// results, err := client.ReadHoldingRegisters(0, 14)
	// fmt.Println(results)
	// fmt.Println(err)

	//[40001, 40003, 40005, 40007, .., 40021, 40023, 40025, 40027]

	read(client)
}

type tag_t struct {
	name    string
	ntype   string
	format  string
	address string
	data    []byte
}

type tag_t2 struct {
	tag    tag_t
	mbaddr modbus2.MBAddr
	area   string
	start  int
	amount int
}

type tag_t2s []*tag_t2

func (self tag_t2s) Len() int {
	return len(self)
}

func (self tag_t2s) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func (self tag_t2s) Less(i, j int) bool {
	return self[i].start < self[j].start
}

func read(client modbus.Client) {
	tags := make(nson.Array, 0)

	tags = append(tags, nson.Message{
		"name":    nson.String("tag1"),
		"type":    nson.String("FLOAT"),
		"format":  nson.String("FLOAT"),
		"address": nson.String("40001"),
	})

	tags = append(tags, nson.Message{
		"name":    nson.String("tag2"),
		"type":    nson.String("FLOAT"),
		"format":  nson.String("FLOAT"),
		"address": nson.String("40003"),
	})

	tags = append(tags, nson.Message{
		"name":    nson.String("tag3"),
		"type":    nson.String("FLOAT"),
		"format":  nson.String("FLOAT"),
		"address": nson.String("40005"),
	})

	tags = append(tags, nson.Message{
		"name":    nson.String("tag4"),
		"type":    nson.String("FLOAT"),
		"format":  nson.String("FLOAT"),
		"address": nson.String("40007"),
	})

	tags = append(tags, nson.Message{
		"name":    nson.String("tag5"),
		"type":    nson.String("FLOAT"),
		"format":  nson.String("FLOAT"),
		"address": nson.String("4009"),
	})

	tags = append(tags, nson.Message{
		"name":    nson.String("tag6"),
		"type":    nson.String("FLOAT"),
		"format":  nson.String("FLOAT"),
		"address": nson.String("40011"),
	})

	tags = append(tags, nson.Message{
		"name":    nson.String("tag7"),
		"type":    nson.String("FLOAT"),
		"format":  nson.String("FLOAT"),
		"address": nson.String("40013"),
	})

	fmt.Println(tags)
	tags2 := make(tag_t2s, 0)

	for _, tag_item := range tags {
		tag, ok := tag_item.(nson.Message)
		if !ok {
			continue
		}

		name, err := tag.GetString("name")
		if err != nil {
			continue
		}

		ntype, err := tag.GetString("type")
		if err != nil {
			continue
		}

		format, err := tag.GetString("format")
		if err != nil {
			continue
		}

		address, err := tag.GetString("address")
		if err != nil {
			continue
		}

		tag2 := tag_t{
			name:    name,
			ntype:   ntype,
			format:  format,
			address: address,
			data:    make([]byte, 0),
		}

		mbaddr := modbus2.ParseTagAddress(address, format)

		fmt.Println(mbaddr)

		tag3 := tag_t2{
			tag:    tag2,
			mbaddr: mbaddr,
			area:   mbaddr.Area,
			start:  mbaddr.Address,
			amount: mbaddr.Size,
		}

		fmt.Println(tag3)

		tags2 = append(tags2, &tag3)
	}

	sort.Sort(tags2)

	read_tasks := cut_group(tags2)

	//fmt.Println(read_tasks)

	for _, item := range read_tasks {
		for _, task := range item {
			fmt.Println(task)

			var bits []byte
			var err error

			switch task.area {
			case "4":
				bits, err = client.ReadHoldingRegisters(uint16(task.start), uint16(task.end))
				if err != nil {
					return
				}
				fmt.Println(bits)
			}

			for _, tag := range task.tags {
				fmt.Println(tag)
				fmt.Println(bits)
				fmt.Println(tag.start - task.start)
				fmt.Println((tag.start - task.start) + tag.amount)
				data := bits[(tag.start-task.start)*2 : (tag.start-task.start)*2+tag.amount]
				fmt.Println(data)
			}
		}
	}
}

type read_task struct {
	area  string
	start int
	end   int
	tags  []*tag_t2
}

var block = 100

func cut_group(tags tag_t2s) map[string][]*read_task {
	read_tasks := make(map[string][]*read_task, 0)

	for _, item := range tags {
		cut_group_c(item.area, &read_tasks, item)
	}

	return read_tasks
}

func cut_group_c(key string, read_tasks *map[string][]*read_task, item *tag_t2) {
	if groups, ok := (*read_tasks)[key]; ok {
		groups_len := len(groups)
		if (item.start + item.amount - groups[groups_len-1].start) < block {
			groups[groups_len-1].end = (item.start + item.amount - 2)
			groups[groups_len-1].tags = append(groups[groups_len-1].tags, item)
		} else {
			rt := read_task{
				area:  item.area,
				start: item.start,
				end:   item.start + item.amount - 2,
				tags:  []*tag_t2{item},
			}

			(*read_tasks)[key] = append((*read_tasks)[key], &rt)
		}
	} else {
		rt := read_task{
			area:  item.area,
			start: item.start,
			end:   item.start + item.amount - 2,
			tags:  []*tag_t2{item},
		}

		(*read_tasks)[key] = []*read_task{&rt}
	}
}
