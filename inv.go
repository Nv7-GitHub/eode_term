package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const PageLength = 20

func RenderPage(items []int) {
	// Get info
	Write("inv", "Loading...")
	names := make([]string, len(items))
	for i, v := range items {
		info, res := Send(MethodElemInfo, map[string]any{"id": v})
		if res.Error != nil {
			Error("inv", "%s", *res.Error)
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

	page := 0

loop:
	for {
		end := PageLength * (page + 1)
		if end > len(ids) {
			end = len(ids)
		}
		RenderPage(ids[PageLength*page : end])
		Write("inv", "Page %d of %d.", page, len(ids)/PageLength)
		fmt.Print("\u001b[34m(prev, next, num, exit):\u001b[0m ")

		line, _, err := reader.ReadLine()
		if err != nil {
			Error("inv", "%s", err.Error())
			continue
		}
		switch string(line) {
		case "exit":
			break loop

		case "next":
			page++
			if page > len(ids)/PageLength {
				page = 0
			}

		case "prev":
			page--
			if page < 0 {
				page = len(ids) / PageLength
			}

		default:
			num, err := strconv.Atoi(string(line))
			if err != nil {
				Error("inv", "Invalid number")
			}
			page = num
		}
	}
}

func init() {
	Cmd("inv", "inv", Inv)
}
