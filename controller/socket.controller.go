package controller

import (
	"app/config"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type socketController struct {
	redisClient *redis.Client
	upgrader    *websocket.Upgrader
	mutexSocket *sync.Mutex
	mapSocket   map[string]*websocket.Conn
}

type SocketController interface {
	AuthSocket(w http.ResponseWriter, r *http.Request)
}

func (c *socketController) AuthSocket(w http.ResponseWriter, r *http.Request) {
	// check auth with uuid
	query := r.URL.Query()
	uuid := query.Get("uuid")
	if uuid == "" {
		badRequest(w, r, errors.New("uuid not found"))
		return
	}

	// check uuid exist in redis
	infoProfile, err := c.redisClient.Get(r.Context(), uuid).Result()
	if err != nil {
		internalServerError(w, r, err)
		return
	}
	if infoProfile == "" {
		internalServerError(w, r, errors.New("uuid not found in redis"))
		return
	}

	// create connect
	conn, err := c.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()

	//connect -> map socket
	c.mutexSocket.Lock()
	c.mapSocket[uuid] = conn
	c.mutexSocket.Unlock()

	// listen connect
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			log.Println("read:", err)
			break
		}
	}

	log.Printf("Disconnect")
}

func NewSocketController() SocketController {
	return &socketController{
		mutexSocket: new(sync.Mutex),
		redisClient: config.GetRedisClient(),
		upgrader:    config.GetUpgraderSocket(),
		mapSocket:   config.GetMapSocket(),
	}
}
