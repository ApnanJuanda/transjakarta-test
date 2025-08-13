package model

type (
	Vehiclelocations struct {
		VehicleId string  `json:"vehicle_id"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Timestamp int64   `json:"timestamp"`
	}

	VehicleLocationPayload struct {
		VehicleId string  `json:"vehicle_id"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Timestamp int64   `json:"timestamp"`
		DestLat   float64 `json:"dest_lat"`
		DestLon   float64 `json:"dest_lon"`
	}

	VehicleHistoryReq struct {
		VehicleId string `json:"vehicle_id"`
		StartTime int64  `json:"start_time"`
		EndTime   int64  `json:"end_time"`
		Limit     int    `json:"limit"`
		Offset    int    `json:"offset"`
	}

	VehicleLocationPublishReq struct {
		VehicleId  string  `json:"vehicle_id"`
		CurrentLat float64 `json:"current_lat"`
		CurrentLon float64 `json:"current_lon"`
		DestLat    float64 `json:"dest_lat"`
		DestLon    float64 `json:"dest_lon"`
	}

	Location struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	GeofenceEvent struct {
		VehicleId string   `json:"vehicle_id"`
		Event     string   `json:"event"`
		Location  Location `json:"location"`
		Timestamp int64    `json:"timestamp"`
	}
)
