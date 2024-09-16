package config

import (
	"github.com/gorilla/websocket"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func GetPsql() *gorm.DB {
	return dbPsql
}

func GetAppPort() string {
	return appPort
}

func GetAppHost() string {
	return appHost
}

func GetSocketPort() string {
	return socketPort
}

func GetRedisClient() *redis.Client {
	return redisClient
}

func GetUpgraderSocket() *websocket.Upgrader {
	return upgraderSocket
}

func GetMapSocket() map[string]*websocket.Conn {
	return mapSocket
}

func GetRabbitmq() *amqp091.Connection {
	return rabbitmq
}
