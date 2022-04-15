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
	fmt.Printf("\u001b[32;1m[%s]\u001b[0m %s\n", cmd, fmt.Sprintf(format, args...))
}

func Error(cmd string, format string, args ...interface{}) {
	fmt.Printf("\u001b[31;1m[%s]\u001b[0m %s\n", cmd, fmt.Sprintf(format, args...))
}

func Clear() {
	fmt.Print("\033[H\033[2J")
}

// Syntax: cmd "arg1 value" arg2

type Command struct {
	Name  string
	Fn    func([]interface{})
	Typs  []Type
	Usage string
}

var cmds = make(map[string]Command)

func Cmd(name string, usage string, fn func([]interface{}), typs ...Type) {
	cmds[name] = Command{Name: name, Fn: fn, Typs: typs, Usage: usage}
}

func Run() {
	for {
		fmt.Print("\u001b[33;1m>>>\u001b[0m ")
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
			conn.Close()
			Write("io", "Exiting...")
			return
		}
		if name == "clear" {
			Clear()
			continue
		}
		if name == "gld" {
			if len(args) != 1 {
				Error("gld", "Expected 1 argument")
				continue
			}
			GuildLogin(args[0])
			continue
		}

		if guild == "" {
			Error(name, "Before playing, login in to a server using %q!", "gld <server ID>")
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
	Clear()

	//Login()
	id = "567132457820749842"    // Nv7
	guild = "705084182673621033" // Elemental on Discord

	// Connect
	Conn(id)
	if guild != "" { // Testing
		Send(MethodGuild, map[string]any{"gld": guild})
	}
	Run()
}
