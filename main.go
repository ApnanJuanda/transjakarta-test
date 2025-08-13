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

	client, err := publish.SetupBroker()
	if err != nil {
		log.Fatal(err)
	}
	publish.SetupSubscriber(gormDB)

	if err = config.Router(gormDB, client); err != nil {
		log.Fatal(err)
	}
}
