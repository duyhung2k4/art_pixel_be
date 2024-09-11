package config

import "flag"

func init() {

	db := flag.Bool("db", false, "")

	flag.Parse()

	loadEnv()
	connectPostgresql(*db)
	connectRedis()
}
