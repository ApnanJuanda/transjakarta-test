package publish

import (
	"encoding/json"
	"fmt"
	"github.com/ApnanJuanda/transjakarta/domain/api/vehicle/model"
	"github.com/ApnanJuanda/transjakarta/lib/env"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
	"log"
	"regexp"
)

func SetupBroker() (client mqtt.Client, err error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(env.String("MQTT_BROKER_ADDRESS", "tcp://localhost:1883"))
	opts.SetClientID(env.String("MQTT_PUBLISHER_ID", "vehicle-location-publisher"))

	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		err = token.Error()
		log.Println("ERROR CONNECT TO MQTT: ", token.Error())
		return
	}
	return
}

func SetupRabbitMQ() (ch *amqp.Channel, err error) {
	conn, err := amqp.Dial(env.String("RABBITMQ_URL", "amqp://guest:guest@localhost:5672"))
	if err != nil {
		log.Println("Failed to connect RabbitMQ:", err)
		return
	}
	ch, err = conn.Channel()
	return
}

func SetupSubscriber(DB *gorm.DB) {
	var dataChan = make(chan model.Vehiclelocations, 100)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(env.String("MQTT_BROKER_ADDRESS", "tcp://localhost:1883"))
	opts.SetClientID(env.String("MQTT_SUBSCRIBER_ID", "vehicle-location-subscriber"))

	opts.SetDefaultPublishHandler(func(client mqtt.Client, message mqtt.Message) {
		var vehicleLocation model.Vehiclelocations
		if err := json.Unmarshal(message.Payload(), &vehicleLocation); err != nil {
			log.Println("ERROR Unmarshal payload: ", err)
			return
		}
		if !validateDataPayload(vehicleLocation) {
			log.Println("ERROR validateDataPayload")
			return
		}

		dataChan <- vehicleLocation
	})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Println("ERROR Connect MQTT Subscriber: ", token.Error())
	}
	topic := fmt.Sprintf("/fleet/vehicle/+/location")
	if token := client.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
		log.Println("ERROR Subscribe MQTT Subscriber: ", token.Error())
	}

	go func() {
		for location := range dataChan {
			err := DB.Table("vehicle_locations").Create(&location).Error
			if err != nil {
				loc, _ := json.Marshal(location)
				log.Printf("ERROR save payload %v: %v", string(loc), err)
				continue
			}
		}
	}()
}

func validateDataPayload(data model.Vehiclelocations) (result bool) {
	result = false
	match, _ := regexp.MatchString(`^[A-Z0-9]+$`, data.VehicleId)
	if !match {
		return
	}
	if data.Latitude < -90 || data.Latitude > 90 {
		return
	}
	if data.Longitude < -180 || data.Longitude > 180 {
		return
	}
	if data.Timestamp <= 0 {
		return
	}
	result = true
	return
}

func ReceiveDataFromRabbitMQ(ch *amqp.Channel) {
	err := ch.ExchangeDeclare(
		"fleet.events",
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("Failed to declare exchange:", err)
	}

	q, err := ch.QueueDeclare(
		"geofence_alerts",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("Failed to declare queue:", err)
	}

	err = ch.QueueBind(
		q.Name,
		"geofence.alert",
		"fleet.events",
		false,
		nil,
	)
	if err != nil {
		log.Println("Failed to bind queue:", err)
	}

	// consume message
	msgs, err := ch.Consume(
		q.Name,
		"",
		true, // auto-ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println("Failed to register consumer:", err)
		return
	}
	go func() {
		for d := range msgs {
			var event model.GeofenceEvent
			err := json.Unmarshal(d.Body, &event)
			if err != nil {
				log.Println("Failed to unmarshal message: ", err)
				continue
			}
			log.Printf("Received geofence event: %v\n", string(d.Body))
		}
	}()
}
