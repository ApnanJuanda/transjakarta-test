package collection

import (
	"github.com/ApnanJuanda/transjakarta/app/controllers/vehicle"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ApiRouter(db *gorm.DB, client mqtt.Client, api *gin.RouterGroup) {
	vehicleCtrl := vehicle.VehicleLocationController(db, client)
	vehicleGroup := api.Group("vehicles")
	{
		vehicleGroup.POST("", vehicleCtrl.CreateVehicleLocation)
		vehicleGroup.POST("/start-publish-data", vehicleCtrl.PublishData)
		vehicleGroup.POST("/stop-publish-data", vehicleCtrl.StopPublishData)
		vehicleGroup.GET("/:vehicle_id", vehicleCtrl.GetLatestVehicleLocation)
		vehicleGroup.GET("/:vehicle_id/history", vehicleCtrl.GetHistoryVehicleLocation)
	}
}
