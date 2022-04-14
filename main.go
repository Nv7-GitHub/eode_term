package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

type Response struct {
	Error *string
}

func R(format string, args ...interface{}) Response {
	res := fmt.Sprintf(format, args...)
	return Response{Error: &res}
}
func Rg() Response { return Response{} }

func Write(cmd string, format string, args ...interface{}) {
	fmt.Printf("\u001b[32m[%s]\u001b[0m %s\n", cmd, fmt.Sprintf(format, args...))
}

func Error(cmd string, format string, args ...interface{}) {
	fmt.Printf("\u001b[31;1m[%s]\u001b[0m %s\n", cmd, fmt.Sprintf(format, args...))
}

func Clear() {
	fmt.Print("\033[H\033[2J")
}

// Syntax: cmd "arg1 value" arg2

type Command struct {
	Fn   func([]interface{})
	Typs []Type
}

var cmds = make(map[string]Command)

func Cmd(name string, fn func([]interface{}), typs ...Type) {
	cmds[name] = Command{Fn: fn, Typs: typs}
}

func Run() {
	for {
		fmt.Print("COMMAND: ")
		raw, _, err := reader.ReadLine()
		if err != nil {
			Error("io", err.Error())
			continue
		}
		vals := strings.SplitN(string(raw), " ", 2)
		name := vals[0]

		// Parse args
		argsRaw := []rune("")
		if len(vals) > 1 {
			argsRaw = []rune(vals[1])
		}
		args := make([]string, 0)
		i := 0
		curr := ""
		for i < len(argsRaw) {
			switch argsRaw[i] {
			case '"':
				i++
				for i < len(argsRaw) && argsRaw[i] != '"' {
					curr += string(argsRaw[i])
					i++
				}

			case ' ':
				args = append(args, curr)
				curr = ""

			default:
				curr += string(argsRaw[i])
			}
			i++
		}
		if curr != "" {
			args = append(args, curr)
			curr = ""
		}

		// Command
		if name == "exit" {
			Write("io", "Exiting...")
			return
		}
		if name == "clear" {
			Clear()
			continue
		}

		cmd, exists := cmds[name]
		if !exists {
			Error("io", "Unknown command: %q", name)
			continue
		}
		argV, r := Match(cmd.Typs, args)
		if r.Error != nil {
			Error(name, "%s", *r.Error)
			continue
		}
		cmd.Fn(argV)
	}
}

func main() {
	Login()
	Run()
}
