package modbus

import (
	"sort"

	"github.com/danclive/capture-driver/util"
	nson "github.com/danclive/nson-go"
	queen "github.com/danclive/queen-go"
)

/*
func read_ext(context queen.Context, msg *nson.Message, conn *conn_t, tags nson.Array) {
	tags_array := make(nson.Array, 0)

	for _, tag_i := range tags {
		tag, ok := tag_i.(nson.Message)
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

		_ = name
		_ = ntype

		mbaddr := ParseTagAddress(address, format)

		var bits []byte

		switch mbaddr.Area {
		case "0":
			bits, err = conn.client.ReadCoils(uint16(mbaddr.Address), uint16(mbaddr.Size))
		case "1":
			bits, err = conn.client.ReadDiscreteInputs(uint16(mbaddr.Address), uint16(mbaddr.Size))
		case "3":
			bits, err = conn.client.ReadInputRegisters(uint16(mbaddr.Address), uint16((mbaddr.Size+1)/2))
		case "4":
			bits, err = conn.client.ReadHoldingRegisters(uint16(mbaddr.Address), uint16((mbaddr.Size+1)/2))
		}

		if err != nil {

			util.ErrorWithFields("read modbus err: ",
				util.Fields{
					"err":     err,
					"tag":     name,
					"address": address,
					"type":    ntype,
					"format":  format,
					"size":    mbaddr.Size,
				})

			conn.islink = false
			// 连接断开时，尝试重新连接
			context.Queen.Emit(
				RECONNECT,
				nson.Message{"id": nson.I64(conn.id), "retry": nson.I32(conn.retry)})
			return
		}

		var value nson.Value

		switch mbaddr.Area {
		case "0", "1":
			if len(bits) > 0 {
				v := false

				if bits[0] > 0 {
					v = true
				}

				value = nson.Bool(v)
			}
		case "3", "4":
			v, o := convert_raw_to_nson(conn, name, mbaddr, bits)
			if o {
				value = v
			}
		}

		if value != nil {
			rtag := nson.Message{
				"name":    nson.String(name),
				"type":    nson.String(ntype),
				"format":  nson.String(format),
				"address": nson.String(address),
				"value":   value,
			}

			tags_array = append(tags_array, rtag)
		}
	}

	msg.Insert("tags", nson.Array(tags_array))
	msg.Insert("ok", nson.Bool(true))
}
*/
func read_ext(context queen.Context, msg *nson.Message, conn *conn_t, tags nson.Array) {
	tags_array := make(nson.Array, 0)

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

		mbaddr := ParseTagAddress(address, format)

		tag3 := tag_t2{
			tag:    tag2,
			mbaddr: mbaddr,
			area:   mbaddr.Area,
			start:  mbaddr.Address,
			amount: mbaddr.Size,
		}

		tags2 = append(tags2, &tag3)
	}

	sort.Sort(tags2)

	read_tasks := cut_group(tags2)

	for _, item := range read_tasks {
		for _, task := range item {
			// fmt.Println(task)

			var bits []byte
			var err error

			switch task.area {
			case "3":
				bits, err = conn.client.ReadInputRegisters(uint16(task.start), uint16(task.end))
			case "4":
				bits, err = conn.client.ReadHoldingRegisters(uint16(task.start), uint16(task.end))
			}

			if err != nil {

				util.ErrorWithFields("read modbus err: ",
					util.Fields{
						"err":   err,
						"area":  task.area,
						"start": task.start,
						"end":   task.end / 2,
					})

				conn.islink = false
				// 连接断开时，尝试重新连接
				context.Queen.Emit(
					RECONNECT,
					nson.Message{"id": nson.I64(conn.id), "retry": nson.I32(conn.retry)})
				return
			}

			for _, tag := range task.tags {
				data := bits[(tag.start-task.start)*2 : (tag.start-task.start)*2+tag.amount]

				value, ok := convert_raw_to_nson(conn, tag.tag.name, tag.mbaddr, data)
				if !ok {
					continue
				}

				rtag := nson.Message{
					"name":    nson.String(tag.tag.name),
					"type":    nson.String(tag.tag.ntype),
					"format":  nson.String(tag.tag.format),
					"address": nson.String(tag.tag.address),
					"value":   value,
				}

				tags_array = append(tags_array, rtag)
			}
		}
	}

	msg.Insert("tags", nson.Array(tags_array))
	msg.Insert("ok", nson.Bool(true))
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
	mbaddr MBAddr
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
