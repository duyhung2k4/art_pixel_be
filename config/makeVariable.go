package config

import "github.com/gorilla/websocket"

func makeVariable() {
	mapSocket = make(map[string]*websocket.Conn)
	mapCheckSendEmail = make(map[string]bool)
}
