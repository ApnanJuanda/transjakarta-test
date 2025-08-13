package vehicle

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ApnanJuanda/transjakarta/domain/api/vehicle/model"
	"github.com/ApnanJuanda/transjakarta/domain/api/vehicle/repository"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"net/http"
	"time"
)

type VehicleServiceInterface interface {
	CreateVehicleLocation(req model.Vehiclelocations) (statusCode int, err error)
	GetLatestVehicleLocation(vehicleId string) (data model.Vehiclelocations, statusCode int, err error)
	GetHistoryVehicleLocation(req model.VehicleHistoryReq) (datas []model.Vehiclelocations, totalData int64, statusCode int, err error)
	StartPublishData(req model.VehicleLocationPublishReq)
	StopPublishData()
}

type vehicleService struct {
	Repository repository.VehicleLocationRepositoryInterface
	Client     mqtt.Client
}

func NewVehicleService(repository repository.VehicleLocationRepositoryInterface, client mqtt.Client) VehicleServiceInterface {
	return &vehicleService{
		Repository: repository,
		Client:     client,
	}
}

func (s *vehicleService) CreateVehicleLocation(req model.Vehiclelocations) (statusCode int, err error) {
	timeStamp := req.Timestamp
	if timeStamp == 0 {
		timeStamp = time.Now().Unix()
	}
	vehicleLocation := model.Vehiclelocations{
		VehicleId: req.VehicleId,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Timestamp: timeStamp,
	}
	err = s.Repository.CreateVehicleLocation(vehicleLocation)
	if err != nil {
		statusCode = http.StatusInternalServerError
		return
	}
	statusCode = http.StatusOK
	return
}

func (s *vehicleService) GetLatestVehicleLocation(vehicleId string) (data model.Vehiclelocations, statusCode int, err error) {
	data, err = s.Repository.GetLatestVehicleLocation(vehicleId)
	if err != nil {
		statusCode = http.StatusInternalServerError
		return
	}
	statusCode = http.StatusOK
	return
}

func (s *vehicleService) GetHistoryVehicleLocation(req model.VehicleHistoryReq) (datas []model.Vehiclelocations, totalData int64, statusCode int, err error) {
	datas, totalData, err = s.Repository.GetHistoryVehicleLocation(req)
	if err != nil {
		statusCode = http.StatusInternalServerError
		return
	}
	statusCode = http.StatusOK
	return
}

var (
	pubActive bool
	pubCancel context.CancelFunc
)

func (s *vehicleService) StartPublishData(req model.VehicleLocationPublishReq) {
	if pubActive {
		log.Println("Publisher is already running")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	pubCancel = cancel
	pubActive = true

	curLat := req.CurrentLat
	curLon := req.CurrentLon

	topic := fmt.Sprintf("/fleet/vehicle/%v/location", req.VehicleId)
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("PUBLISHER STOP")
				return
			default:
				curLat, curLon = moveToDestination(curLat, curLon, req.DestLat, req.DestLon)
				data := model.Vehiclelocations{
					VehicleId: req.VehicleId,
					Latitude:  curLat,
					Longitude: curLon,
					Timestamp: time.Now().Unix(),
				}
				payload, err := json.Marshal(data)
				if err != nil {
					log.Println("ERROR Convert to json: ", err)
					return
				}
				token := s.Client.Publish(topic, 1, false, payload)
				token.Wait()
				if err := token.Error(); err != nil {
					log.Printf("ERROR Publish data: %v", err)
					continue
				}
				log.Printf("SUCCESS publish %v TO topic %s ", string(payload), topic)
				time.Sleep(2 * time.Second)
			}
		}
	}()
}

func (s *vehicleService) StopPublishData() {
	if !pubActive {
		log.Println("Publisher is already not active")
		return
	}
	pubCancel()
	pubActive = false
	return
}

func moveToDestination(startLat, startLon, destLat, destLon float64) (newLat, newLon float64) {
	dLat := destLat - startLat
	dLon := destLon - startLon

	stepFraction := 0.1
	newLat = startLat + dLat*stepFraction
	newLon = startLon + dLon*stepFraction
	return
}
