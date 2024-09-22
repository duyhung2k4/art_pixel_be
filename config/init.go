package config

import (
	"flag"

	"github.com/go-chi/jwtauth/v5"
)

func init() {
	db := flag.Bool("db", false, "")
	jwt = jwtauth.New("HS256", []byte("h-shop"), nil)

	flag.Parse()

	// connect
	loadEnv()
	makeVariable()
	connectPostgresql(*db)
	connectRedis()
	connectMongoDB()
	createFolder()
	initSocket()
	connectRabbitmq()
	initSmptAuth()
}
