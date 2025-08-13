package repository

import (
	"github.com/ApnanJuanda/transjakarta/domain/api/vehicle/model"
	"gorm.io/gorm"
)

type VehicleLocationRepositoryInterface interface {
	CreateVehicleLocation(req model.Vehiclelocations) (err error)
	GetLatestVehicleLocation(vehicleId string) (response model.Vehiclelocations, err error)
	GetHistoryVehicleLocation(req model.VehicleHistoryReq) (responses []model.Vehiclelocations, totalRow int64, err error)
}

type vehicleLocationRepository struct {
	DB *gorm.DB
}

func NewVehicleLocationRepository(db *gorm.DB) VehicleLocationRepositoryInterface {
	return &vehicleLocationRepository{DB: db}
}

func (r *vehicleLocationRepository) CreateVehicleLocation(req model.Vehiclelocations) (err error) {
	err = r.DB.Table("vehicle_locations").Create(&req).Error
	return
}

func (r *vehicleLocationRepository) GetLatestVehicleLocation(vehicleId string) (response model.Vehiclelocations, err error) {
	err = r.DB.Table("vehicle_locations").Where("vehicle_id = ?", vehicleId).
		Order("timestamp DESC").First(&response).Error
	return
}

func (r *vehicleLocationRepository) GetHistoryVehicleLocation(req model.VehicleHistoryReq) (responses []model.Vehiclelocations, totalRow int64, err error) {
	err = r.DB.Table("vehicle_locations").
		Where("vehicle_id = ? AND timestamp BETWEEN ? AND ?", req.VehicleId, req.StartTime, req.EndTime).
		Limit(req.Limit).Offset(req.Offset).Count(&totalRow).Find(&responses).Error
	return
}
