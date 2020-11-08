package vehicle

import (
	"errors"
	"fmt"
	"math"
	"net/http"

	"app_api/shared"
	gmConnector "app_api/shared/gm"
)

const (
	ENGINE_START = "START"
	ENGINE_STOP  = "STOP"
)

// Service ... represents an instance of the vehicle package service interface
type Service interface {
	GetVehicle(vehicleID int64) (res Vehicle, err *shared.APIError)
	GetVehicleDoors(vehicleID int64) (res []gmConnector.GMVehicleDoorData, err *shared.APIError)
	GetVehicleFuel(vehicleID int64) (res Fuel, err *shared.APIError)
	GetVehicleBattery(vehicleID int64) (res Battery, err *shared.APIError)
	SendEngineAction(vehicleID int64, engineAction EngineActionRequest) (engineSubmissionStatus EngineActionResponse, err *shared.APIError)
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
func (s *service) GetVehicleDoors(vehicleID int64) (res []gmConnector.GMVehicleDoorData, err *shared.APIError) {
	res, err = s.gm.GetVehicleDoors(vehicleID)
	if err != nil {
		return
	}

	return
}

// GetVehicleFuel ... returns the status of the fiel for a given car
func (s *service) GetVehicleFuel(vehicleID int64) (res Fuel, err *shared.APIError) {
	fuel, _, err := s.gm.GetVehicleEnergyStatus(vehicleID)
	if err != nil {
		return
	}

	if fuel == nil {
		return
	}

	percentage := math.Round(*fuel*100) / 100
	res.Percentage = &percentage

	return
}

// GetVehicleBattery ... returns the status of the fiel for a given car
func (s *service) GetVehicleBattery(vehicleID int64) (res Battery, err *shared.APIError) {
	_, battery, err := s.gm.GetVehicleEnergyStatus(vehicleID)
	if err != nil {
		return
	}

	if battery == nil {
		return
	}

	percentage := math.Round(*battery*100) / 100
	res.Percentage = &percentage

	return
}

// EngineActionRequest response
//
// swagger:model EngineActionRequest
type EngineActionRequest struct {
	// Action
	//
	// required: true
	// example: START
	Action string `json:"action"`
}

// EngineActionResponse response
//
// swagger:model EngineActionResponse
type EngineActionResponse struct {
	// Status
	//
	// required: true
	// example: success
	Action string `json:"status"`
}

// SendEngineAction ... attempts to send the client request to GM API /actionEngineService
func (s *service) SendEngineAction(vehicleID int64, engineAction EngineActionRequest) (engineSubmissionStatus EngineActionResponse, err *shared.APIError) {
	var action string

	switch engineAction.Action {
	case ENGINE_START:
		action = gmConnector.ENGINE_START
	case ENGINE_STOP:
		action = gmConnector.ENGINE_STOP
	default:
		errorMessage := "Unsupported engine action option"
		engineActionError := fmt.Errorf("Unsupported data type: %s", engineAction.Action)
		err = shared.NewAPIError(http.StatusInternalServerError, engineActionError, errorMessage)
		return
	}

	engineResponse, err := s.gm.SendVehicleEngineAction(vehicleID, action)
	if err != nil {
		return
	}

	switch engineResponse.Status {
	case gmConnector.EXECUTED:
		engineSubmissionStatus.Action = "success"
	case gmConnector.FAILED:
		engineSubmissionStatus.Action = "error"
	default:
		errorMessage := "Failed to read response from GM"
		engineActionError := fmt.Errorf("Unsupported data type: %s", engineResponse.Status)
		err = shared.NewAPIError(http.StatusInternalServerError, engineActionError, errorMessage)
		return
	}

	return
}
