package config

import (
	"github.com/gorilla/websocket"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	appPort    string
	appHost    string
	socketPort string

	dbHost      string
	dbPort      string
	dbName      string
	dbUser      string
	dbPassword  string
	redisUrl    string
	rabbitmqUrl string

	dbPsql         *gorm.DB
	redisClient    *redis.Client
	rabbitmq       *amqp091.Connection
	upgraderSocket *websocket.Upgrader

	mapSocket map[string]*websocket.Conn
)
