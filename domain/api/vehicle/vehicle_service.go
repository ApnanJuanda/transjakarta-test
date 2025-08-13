package vehicle

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ApnanJuanda/transjakarta/domain/api/vehicle/model"
	"github.com/ApnanJuanda/transjakarta/domain/api/vehicle/repository"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"math"
	"net/http"
	"time"
)

type VehicleServiceInterface interface {
	CreateVehicleLocation(req model.Vehiclelocations) (statusCode int, err error)
	GetLatestVehicleLocation(vehicleId string) (data model.Vehiclelocations, statusCode int, err error)
	GetHistoryVehicleLocation(req model.VehicleHistoryReq) (datas []model.Vehiclelocations, totalData int64, statusCode int, err error)
	StartPublishData(req model.VehicleLocationPublishReq)
	StopPublishData()
	MoveToDestination(startLat, startLon, destLat, destLon float64) (newLat, newLon float64)
}

type vehicleService struct {
	Repository repository.VehicleLocationRepositoryInterface
	Client     mqtt.Client
	Channel    *amqp.Channel
}

func NewVehicleService(repository repository.VehicleLocationRepositoryInterface, client mqtt.Client, channel *amqp.Channel) VehicleServiceInterface {
	return &vehicleService{
		Repository: repository,
		Client:     client,
		Channel:    channel,
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
				curLat, curLon = s.MoveToDestination(curLat, curLon, req.DestLat, req.DestLon)
				if !pubActive {
					break
				}
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

				// push to rabbitMQ
				distance := calculateDistance(data.Latitude, data.Longitude, req.DestLat, req.DestLon)
				log.Printf("INFO distance: %v meter", distance)
				if distance <= 50 {
					s.PublishEventToRabbitMQ(data)
				}

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

func (s *vehicleService) MoveToDestination(startLat, startLon, destLat, destLon float64) (newLat, newLon float64) {
	stepFraction := 0.8
	dLat := destLat - startLat
	dLon := destLon - startLon

	if dLat == 0 && dLon == 0 {
		s.StopPublishData()
		return
	}
	newLat = startLat + dLat*stepFraction
	newLon = startLon + dLon*stepFraction
	return
}

func calculateDistance(curLat, curLon, destLat, destLon float64) (distance float64) {
	const earthRaidus = 6371000 // m
	dLat := (destLat - curLat) * math.Pi / 180.0
	dLon := (destLon - curLon) * math.Pi / 180.0

	lat1Rad := curLat * math.Pi / 180.0
	lat2Rad := destLat * math.Pi / 180.0

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance = earthRaidus * c
	return
}

func (s *vehicleService) PublishEventToRabbitMQ(data model.Vehiclelocations) {
	err := s.Channel.ExchangeDeclare(
		"fleet.events",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to declare exchange: %v", err)
		return
	}

	q, err := s.Channel.QueueDeclare("geofence_alerts", true, false, false, false, nil)
	if err != nil {
		log.Println("Failed to declare queue:", err)
		return
	}
	err = s.Channel.QueueBind(q.Name, "geofence.alert", "fleet.events", false, nil)
	if err != nil {
		log.Println("Failed to bind queue:", err)
		return
	}

	// create payload
	event := model.GeofenceEvent{
		VehicleId: data.VehicleId,
		Event:     "geofence_entry",
		Location: model.Location{
			Latitude:  data.Latitude,
			Longitude: data.Longitude,
		},
		Timestamp: data.Timestamp,
	}

	// Send message
	body, _ := json.Marshal(event)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = s.Channel.PublishWithContext(
		ctx,
		"fleet.events",
		"geofence.alert",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish a message: %v", err)
		return
	}
	log.Printf("Geofence event sent to RabbitMQ: %v", string(body))
}
