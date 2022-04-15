package main

import (
	"encoding/json"
	"strconv"
)

func Combo(args []interface{}) {
	// Get ids
	ids := make([]int, len(args))
	for i, arg := range args {
		id, res := Send(MethodElem, map[string]any{"name": arg})
		if res.Error != nil {
			Error("combo", "%s", *res.Error)
			return
		}
		ids[i] = int(id["id"].(float64))
	}

	// Combine
	d, res := Send(MethodCombo, map[string]any{"elems": ids})
	if res.Error != nil {
		Error("combo", "%s", *res.Error)
		return
	}
	el := int(d["id"].(float64)) // Elem3 ID

	// Get info
	info, res := Send(MethodElemInfo, map[string]any{"ids": []int{el}})
	if res.Error != nil {
		Error("combo", "%s", *res.Error)
		return
	}
	var elem Element
	json.Unmarshal([]byte(info[strconv.Itoa(el)].(string)), &elem)

	// Result
	exists := d["exists"].(bool)
	if exists {
		Write("combo", "You made %s, but already have it.", elem.Name)
	} else {
		Write("combo", "You made %s!", elem.Name)
	}
}

func init() {
	Cmd("combo", "combo [elem] [elem2] (elem3...)", Combo, STRING, VARIADIC)
}
