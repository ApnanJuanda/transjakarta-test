package config

import (
	"github.com/ApnanJuanda/transjakarta/app/controllers/root"
	"github.com/ApnanJuanda/transjakarta/config/collection"
	"github.com/ApnanJuanda/transjakarta/lib/env"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Router(DB *gorm.DB, client mqtt.Client) error {
	router := gin.Default()
	corsConfig(router)

	loyaltyGroup := router.Group(env.String("ROOT_PATH", "fleet"))

	loyaltyGroup.GET("/", root.Index)

	api := loyaltyGroup.Group("/api")
	collection.ApiRouter(DB, client, api)

	if err := router.Run(); err != nil {
		return err
	}
	return nil
}
