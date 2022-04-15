package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

const MaxWidth = 100

func Wrap(val string) string {
	out := &strings.Builder{}
	width := 0
	for _, char := range val {
		out.WriteRune(char)
		width++
		if width > MaxWidth {
			width = 0
			out.WriteRune('\n')
		}
	}
	return out.String()
}

func Info(pars []interface{}) {
	name := pars[0].(string)
	id, res := Send(MethodElem, map[string]any{"name": name})
	if res.Error != nil {
		Error("combo", "%s", *res.Error)
		return
	}
	info, res := Send(MethodElemInfo, map[string]any{"id": int(id["id"].(float64))})
	if res.Error != nil {
		Error("combo", "%s", *res.Error)
		return
	}
	var el map[string]interface{}
	json.Unmarshal([]byte(info["data"].(string)), &el)

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Key", "Value"})
	for k, v := range el {
		if k == "Parents" {
			continue
		}

		t.AppendRow(table.Row{"\u001b[34m" + k + "\u001b[0m", Wrap(fmt.Sprintf("%v", v))})
		t.AppendSeparator()
	}
	t.Render()
}

func init() {
	Cmd("info", "info [element]", Info, STRING)
}
