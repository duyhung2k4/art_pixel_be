package config

import (
	"net/smtp"

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

func GetPythonNodePort() string {
	return pythonNodePort
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

func GetMapCheckSendEmail() map[string]bool {
	return mapCheckSendEmail
}

func GetRabbitmq() *amqp091.Connection {
	return rabbitmq
}

func GetSmtpPort() string {
	return smtpPort
}

func GetSmtpHost() string {
	return smtpHost
}

func GetAuthSmtp() smtp.Auth {
	return authSmtp
}
