package config

import (
	"net/smtp"

	"github.com/go-chi/jwtauth/v5"
	"github.com/gorilla/websocket"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	appPort        string
	appHost        string
	socketPort     string
	pythonNodePort string

	dbHost       string
	dbPort       string
	dbName       string
	dbUser       string
	dbPassword   string
	redisUrl     string
	rabbitmqUrl  string
	smtpEmail    string
	smtpHost     string
	smtpPort     string
	smtpPassword string

	authSmtp       smtp.Auth
	dbPsql         *gorm.DB
	redisClient    *redis.Client
	rabbitmq       *amqp091.Connection
	upgraderSocket *websocket.Upgrader

	mapSocket         map[string]*websocket.Conn
	mapSocketEvent    map[string]map[string]*websocket.Conn
	mapCheckSendEmail map[string]bool
	jwt               *jwtauth.JWTAuth
)
