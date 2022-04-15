package main

import (
	"os"
	"sort"

	"github.com/jedib0t/go-pretty/v6/table"
)

func Help(args []interface{}) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Command", "Usage"})
	vals := make([]Command, len(cmds))
	i := 0
	for _, v := range cmds {
		vals[i] = v
		i++
	}
	sort.Slice(vals, func(i, j int) bool {
		return vals[i].Name < vals[j].Name
	})
	for _, cmd := range cmds {
		t.AppendRow(table.Row{cmd.Name, cmd.Usage})
	}
	t.Render()
}

func init() {
	Cmd("help", "help", Help)
}
