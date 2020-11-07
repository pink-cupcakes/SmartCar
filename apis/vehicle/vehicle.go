package vehicle

import (
	"errors"
	"net/http"

	"app_api/shared"
	gmConnector "app_api/shared/gm"
)

// Service ... represents an instance of the vehicle package service interface
type Service interface {
	GetVehicle(vehicleID int64) (res Vehicle, err *shared.APIError)
	GetVehicleDoors(vehicleID int64) (res []Door, err *shared.APIError)
}

// NewService ... returns an instance of the vehicle package service
func NewService(gmAPIConnector gmConnector.GMAPIConnector) Service {
	return &service{
		gm: gmAPIConnector,
	}
}

type service struct {
	gm gmConnector.GMAPIConnector
}

// GetVehicle ... returns an overview for a given car
func (s *service) GetVehicle(vehicleID int64) (res Vehicle, err *shared.APIError) {
	gmVehicleData, err := s.gm.GetVehicle(vehicleID)
	if err != nil {
		return
	}

	var doorCount int64

	if gmVehicleData.IsFourDoor && gmVehicleData.IsTwoDoor {
		requestErr := errors.New("GM responded with both two and four door as true")
		clientErr := "Failed to get vehicle"
		err = shared.NewAPIError(http.StatusInternalServerError, requestErr, clientErr).SetInternalErrorMessage("getVehicle: GM responded the car has two and four doors")
		return
	}

	if gmVehicleData.IsFourDoor {
		doorCount = 4
	} else if gmVehicleData.IsTwoDoor {
		doorCount = 2
	} else {
		requestErr := errors.New("Invalid door count response from GM")
		clientErr := "Failed to get vehicle"
		err = shared.NewAPIError(http.StatusInternalServerError, requestErr, clientErr).SetInternalErrorMessage("getVehicle: Failed to get a valid doorCount from GM")
		return
	}

	res = Vehicle{gmVehicleData.Vin, gmVehicleData.Color, doorCount, gmVehicleData.DriveTrain}
	return
}

// GetVehicleDoors ... returns the status of the doors for a given car
func (s *service) GetVehicleDoors(vehicleID int64) (res []Door, err *shared.APIError) {
	_, err = s.gm.GetVehicleDoors(vehicleID)
	if err != nil {
		return
	}

	return
}
