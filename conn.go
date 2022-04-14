package main

import (
	"encoding/json"
	"os"

	"github.com/gorilla/websocket"
)

const serv = "wss://http.nv7haven.com/eode"

var conn *websocket.Conn

func Send(m Method, data map[string]any) (map[string]any, Response) {
	d := Message{Method: m, Params: data}
	v, _ := json.Marshal(d) // Will never be error

	// Send
	err := conn.WriteMessage(websocket.TextMessage, v)
	if err != nil {
		Error("conn", "Send error: %s", err.Error())
		os.Exit(1)
	}

	// Recv
	_, message, err := conn.ReadMessage()
	if err != nil {
		Error("conn", "Receive error: %s", err.Error())
		os.Exit(1)
	}
	var val Resp
	err = json.Unmarshal(message, &val)
	if err != nil {
		Error("conn", "Parse error: %s", err.Error())
		os.Exit(1)
	}

	// Return
	if val.Error != nil {
		return nil, R(*val.Error)
	}
	return val.Data, Rg()
}

func Conn(id string) {
	var err error
	conn, _, err = websocket.DefaultDialer.Dial(serv, nil)
	if err != nil {
		Error("conn", "Connect error: %s", err.Error())
		os.Exit(1)
	}

	err = conn.WriteMessage(websocket.TextMessage, []byte(id))
	if err != nil {
		Error("conn", "Send error: %s", err.Error())
		os.Exit(1)
	}
}
