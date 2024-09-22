package config

import (
	"net/smtp"

	"github.com/go-chi/jwtauth/v5"
	"github.com/gorilla/websocket"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var (
	appPort        string
	appHost        string
	socketPort     string
	pythonNodePort string

	dbHost        string
	dbPort        string
	dbName        string
	dbUser        string
	dbPassword    string
	redisUrl      string
	rabbitmqUrl   string
	smtpEmail     string
	smtpHost      string
	smtpPort      string
	smtpPassword  string
	mongodbUrl    string
	mongoDatabase string

	authSmtp       smtp.Auth
	dbPsql         *gorm.DB
	redisClient    *redis.Client
	rabbitmq       *amqp091.Connection
	upgraderSocket *websocket.Upgrader
	mongodb        *mongo.Database

	mapSocket         map[string]*websocket.Conn
	mapCheckSendEmail map[string]bool
	jwt               *jwtauth.JWTAuth
)
