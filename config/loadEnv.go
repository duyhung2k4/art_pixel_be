package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	appHost = os.Getenv(APP_HOST)
	appPort = os.Getenv(APP_PORT)
	socketPort = os.Getenv(SOCKET_PORT)

	dbHost = os.Getenv(DB_HOST)
	dbPort = os.Getenv(DB_PORT)
	dbName = os.Getenv(DB_NAME)
	dbUser = os.Getenv(DB_USER)
	dbPassword = os.Getenv(DB_PASSWORD)
	redisUrl = os.Getenv(REDIS_URL)
}
