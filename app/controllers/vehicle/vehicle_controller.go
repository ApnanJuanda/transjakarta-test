package vehicle

import (
	"github.com/ApnanJuanda/transjakarta/domain/api/vehicle"
	"github.com/ApnanJuanda/transjakarta/domain/api/vehicle/model"
	"github.com/ApnanJuanda/transjakarta/domain/api/vehicle/repository"
	"github.com/ApnanJuanda/transjakarta/lib/response"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type vehicleLocationController struct {
	vehicleLocationService vehicle.VehicleServiceInterface
}

func VehicleLocationController(DB *gorm.DB, client mqtt.Client) *vehicleLocationController {
	return &vehicleLocationController{
		vehicleLocationService: vehicle.NewVehicleService(repository.NewVehicleLocationRepository(DB), client),
	}
}

func (c *vehicleLocationController) PublishData(ctx *gin.Context) {
	var req model.VehicleLocationPublishReq
	if err := ctx.BindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	go c.vehicleLocationService.StartPublishData(req)
	response.Json(ctx, http.StatusOK, nil)
}

func (c *vehicleLocationController) StopPublishData(ctx *gin.Context) {
	go c.vehicleLocationService.StopPublishData()
	response.Json(ctx, http.StatusOK, nil)
}

func (c *vehicleLocationController) CreateVehicleLocation(ctx *gin.Context) {
	var req model.Vehiclelocations
	if err := ctx.BindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	statusCode, err := c.vehicleLocationService.CreateVehicleLocation(req)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Json(ctx, statusCode, nil)
}

func (c *vehicleLocationController) GetLatestVehicleLocation(ctx *gin.Context) {
	vehicleId := ctx.Param("vehicle_id")
	data, statusCode, err := c.vehicleLocationService.GetLatestVehicleLocation(vehicleId)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.Json(ctx, statusCode, data)
}

func (c *vehicleLocationController) GetHistoryVehicleLocation(ctx *gin.Context) {
	vehicleId := ctx.Param("vehicle_id")
	startTime, _ := strconv.ParseInt(ctx.Query("start"), 10, 64)
	endTime, _ := strconv.ParseInt(ctx.Query("end"), 10, 64)

	page := ctx.DefaultQuery("page", "1")
	perPage := ctx.DefaultQuery("limit", "10")
	convPage, _ := strconv.Atoi(page)
	convPerPage, _ := strconv.Atoi(perPage)
	offset := (convPage - 1) * convPerPage

	req := model.VehicleHistoryReq{
		VehicleId: vehicleId,
		StartTime: startTime,
		EndTime:   endTime,
		Limit:     convPerPage,
		Offset:    offset,
	}
	datas, totalData, statusCode, err := c.vehicleLocationService.GetHistoryVehicleLocation(req)
	if err != nil {
		response.Error(ctx, statusCode, err.Error())
		return
	}
	response.JsonPagination(ctx, statusCode, datas, convPage, convPerPage, totalData)
}
