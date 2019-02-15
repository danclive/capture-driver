package snap7

import (
	"fmt"
	"sort"

	nson "github.com/danclive/nson-go"
	queen "github.com/danclive/queen-go"
	snap7 "github.com/danclive/snap7-go"
)

var block int = 100

type tag_t struct {
	name    string
	ntype   string
	format  string
	address string
	data    []byte
}

type tag_t2 struct {
	tag    tag_t
	s7addr S7Addr
	area   int
	db     int
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

func read_ext(context queen.Context, msg *nson.Message, conn *conn_t, tags nson.Array) {
	tags2 := make(tag_t2s, 0)

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

		tag2 := tag_t{
			name:    name,
			ntype:   ntype,
			format:  format,
			address: address,
			data:    make([]byte, 0),
		}

		s7addr := ParseTagAddress(address)

		tag3 := tag_t2{
			tag:    tag2,
			s7addr: s7addr,
			area:   s7addr.Cmd[0],
			db:     s7addr.Cmd[1],
			start:  s7addr.Cmd[2],
			amount: s7addr.Cmd[3],
		}

		tags2 = append(tags2, &tag3)
	}

	sort.Sort(tags2)

	read_tasks := cut_group(tags2)

	for _, item := range read_tasks {
		for _, task := range item {
			bits, err := conn.client.ReadArea(task.area, task.db, task.start, task.end-task.start+1)
			if err != nil {
				conn.islink = false
				// 连接断开时，尝试重新连接
				context.Queen.Emit(
					RECONNECT,
					nson.Message{"id": nson.I32(conn.id), "retry": nson.I32(conn.retry)})
				return
			}

			for _, tag := range task.tags {
				tag.tag.data = bits[tag.start-task.start : (tag.start-task.start)+tag.amount]
			}
		}
	}

	tags_array := make(nson.Array, 0)

	for _, item := range tags2 {
		//fmt.Println(item)
		bs, ok := convert_raw_to_nson(item)
		if ok {
			//fmt.Println(bs)
			tag := nson.Message{
				"name":    nson.String(item.tag.name),
				"type":    nson.String(item.tag.ntype),
				"format":  nson.String(item.tag.format),
				"address": nson.String(item.tag.address),
				"value":   bs,
			}

			tags_array = append(tags_array, tag)
		}
	}

	msg.Insert("tags", nson.Array(tags_array))
	msg.Insert("ok", nson.Bool(true))
}

type read_task struct {
	area  int
	db    int
	start int
	end   int
	tags  []*tag_t2
}

func cut_group(tags tag_t2s) map[string][]*read_task {

	read_tasks := make(map[string][]*read_task, 0)

	for _, item := range tags {
		switch item.area {
		case snap7.S7AreaPE:
			cut_group_c("PE", &read_tasks, item)
		case snap7.S7AreaPA:
			cut_group_c("PA", &read_tasks, item)
		case snap7.S7AreaMK:
			cut_group_c("MK", &read_tasks, item)
		case snap7.S7AreaCT:
			cut_group_c("CT", &read_tasks, item)
		case snap7.S7AreaTM:
			cut_group_c("TM", &read_tasks, item)
		case snap7.S7AreaDB:
			key := fmt.Sprintf("DB%v", item.db)
			cut_group_c(key, &read_tasks, item)
		}
	}

	return read_tasks
}

func cut_group_c(key string, read_tasks *map[string][]*read_task, item *tag_t2) {
	if groups, ok := (*read_tasks)[key]; ok {
		groups_len := len(groups)
		if (item.start + item.amount - groups[groups_len-1].start) < block {
			groups[groups_len-1].end = (item.start + item.amount - 1)
			groups[groups_len-1].tags = append(groups[groups_len-1].tags, item)
		} else {
			rt := read_task{
				area:  item.area,
				db:    item.db,
				start: item.start,
				end:   item.start + item.amount - 1,
				tags:  []*tag_t2{item},
			}

			(*read_tasks)[key] = append((*read_tasks)[key], &rt)
		}
	} else {
		rt := read_task{
			area:  item.area,
			db:    item.db,
			start: item.start,
			end:   item.start + item.amount - 1,
			tags:  []*tag_t2{item},
		}

		(*read_tasks)[key] = []*read_task{&rt}
	}
}
