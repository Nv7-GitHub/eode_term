package main

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
)

const PageLength = 20

func RenderPage(items []int, cmd string) {
	// Get info
	Write(cmd, "Loading...")
	names := make([]string, len(items))
	for i, v := range items {
		info, res := Send(MethodElemInfo, map[string]any{"id": v})
		if res.Error != nil {
			Error(cmd, "%s", *res.Error)
			return
		}
		var el Element
		json.Unmarshal([]byte(info["data"].(string)), &el)
		names[i] = el.Name
	}

	Clear()
	for _, name := range names {
		fmt.Println("\u001b[35m" + name + "\u001b[0m")
	}
}

func PageSwitcher(elems []int, cmd string) {
	page := 0

loop:
	for {
		end := PageLength * (page + 1)
		if end > len(elems) {
			end = len(elems)
		}
		RenderPage(elems[PageLength*page:end], cmd)
		Write(cmd, "Page %d of %d.", page, len(elems)/PageLength)
		fmt.Print("\u001b[34m(prev, next, num, exit):\u001b[0m ")

		line, _, err := reader.ReadLine()
		if err != nil {
			Error(cmd, "%s", err.Error())
			continue
		}
		switch string(line) {
		case "exit":
			break loop

		case "next":
			page++
			if page > len(elems)/PageLength {
				page = 0
			}

		case "prev":
			page--
			if page < 0 {
				page = len(elems) / PageLength
			}

		default:
			num, err := strconv.Atoi(string(line))
			if err != nil {
				Error(cmd, "Invalid number")
			}
			page = num
		}
	}
}

func Inv(_ []interface{}) {
	inv, res := Send(MethodInv, map[string]any{})
	if res.Error != nil {
		Error("inv", "%s", *res.Error)
		return
	}
	vals := inv["elems"].([]interface{})
	ids := make([]int, len(vals))
	for i, val := range vals {
		ids[i] = int(val.(float64))
	}
	sort.Ints(ids)

	PageSwitcher(ids, "inv")
}

func Cat(args []interface{}) {
	name := args[0].(string)
	cat, res := Send(MethodCategory, map[string]any{"name": name})
	if res.Error != nil {
		Error("cat", "%s", *res.Error)
		return
	}
	vals := cat["elems"].([]interface{})
	ids := make([]int, len(vals))
	for i, val := range vals {
		ids[i] = int(val.(float64))
	}
	sort.Ints(ids)

	PageSwitcher(ids, "cat")
}

func init() {
	Cmd("inv", "inv", Inv)
	Cmd("cat", "cat [name]", Cat, STRING)
}