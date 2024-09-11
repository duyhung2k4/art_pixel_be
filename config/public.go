package config

import (
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
