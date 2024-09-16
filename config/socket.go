package config

import (
	"net/http"

	"github.com/gorilla/websocket"
)

func initSocket() {
	upgraderSocket = &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}
