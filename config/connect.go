package config

import (
	"app/model"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectPostgresql(migrate bool) error {
	var err error
	dns := fmt.Sprintf(
		`
			host=%s
			user=%s
			password=%s
			dbname=%s
			port=%s
			sslmode=disable`,
		dbHost,
		dbUser,
		dbPassword,
		dbName,
		dbPort,
	)

	dbPsql, err = gorm.Open(postgres.Open(dns), &gorm.Config{})

	if migrate {
		errMigrate := dbPsql.AutoMigrate(
			&model.Profile{},
			&model.Face{},
		)

		if errMigrate != nil {
			return errMigrate
		}
	}

	return err
}

func connectRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: redisUrl,
	})
}

func connectRabbitmq() error {
	var err error
	rabbitmq, err = amqp091.Dial(rabbitmqUrl)
	if err != nil {
		rabbitmq.Close()
	}
	return err
}
