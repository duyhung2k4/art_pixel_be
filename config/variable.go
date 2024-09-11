package config

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	appPort    string
	appHost    string
	socketPort string

	dbHost     string
	dbPort     string
	dbName     string
	dbUser     string
	dbPassword string
	redisUrl   string

	dbPsql      *gorm.DB
	redisClient *redis.Client
)
