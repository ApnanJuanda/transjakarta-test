package main

import (
	"github.com/ApnanJuanda/transjakarta/config"
	"github.com/ApnanJuanda/transjakarta/db"
	"github.com/ApnanJuanda/transjakarta/lib/publish"
	_ "github.com/joho/godotenv/autoload"
	"log"
)

func main() {
	gormDB, sqlDB, err := db.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = sqlDB.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	channel, err := publish.SetupRabbitMQ()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = channel.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	client, err := publish.SetupBroker()
	if err != nil {
		log.Fatal(err)
	}

	publish.SetupSubscriber(gormDB)
	publish.ReceiveDataFromRabbitMQ(channel)

	if err = config.Router(gormDB, client, channel); err != nil {
		log.Fatal(err)
	}
}
