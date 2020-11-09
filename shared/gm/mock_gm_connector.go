package gmapiconnector

import (
	"app_api/shared"
	"fmt"
	"net/http"
)

type mockGMAPIConnector struct{}

func NewMockGMAPIConnector() GMAPIConnector {
	return &mockGMAPIConnector{}
}

/** GetVehicle ... in this package is just a mocked response for testing purposes. It will respond with a static vehicle object on success,
and an error for any vehicleID not (1234, 1235)
*/
func (mg *mockGMAPIConnector) GetVehicle(vehicleID int64) (res gmVehicleData, err *shared.APIError) {
	if vehicleID != 1234 && vehicleID != 1235 {
		gmErrorCode := 404
		gmErrorMessage := fmt.Sprintf("Vehicle id: %d not found.", vehicleID)
		requestErr := fmt.Errorf("Failed to GET vehicle from GM, non-200 response: %s Response code: %d", gmErrorMessage, gmErrorCode)
		clientErr := "Failed to get vehicle"
		err = shared.NewAPIError(http.StatusInternalServerError, requestErr, clientErr)
		return
	}

	res = gmVehicleData{"123123412412", "Metallic Silver", true, false, "v8"}
	return
}

// GetVehicleDoors ... mocked logic for gm_connector for testing purposes
func (mg *mockGMAPIConnector) GetVehicleDoors(vehicleID int64) (res []GMVehicleDoorData, err *shared.APIError) {
	if vehicleID != 1234 && vehicleID != 1235 {
		gmErrorCode := 404
		gmErrorMessage := fmt.Sprintf("Vehicle id: %d not found.", vehicleID)
		requestErr := fmt.Errorf("Failed to GET vehicle doors from GM, non-200 response: %s Response code: %d", gmErrorMessage, gmErrorCode)
		clientErr := "Failed to get vehicle"
		err = shared.NewAPIError(http.StatusInternalServerError, requestErr, clientErr)
		return
	}

	res = append(res, GMVehicleDoorData{"frontLeft", true})
	res = append(res, GMVehicleDoorData{"frontRight", true})
	return
}

// GetVehicleEnergyStatus ... mocked logic for gm_connector for testing purposes
func (mg *mockGMAPIConnector) GetVehicleEnergyStatus(vehicleID int64) (fuelLevel, batteryLevel *float64, err *shared.APIError) {
	if vehicleID != 1234 && vehicleID != 1235 {
		gmErrorCode := 404
		gmErrorMessage := fmt.Sprintf("Vehicle id: %d not found.", vehicleID)
		requestErr := fmt.Errorf("Failed to GET vehicle energy status from GM, non-200 response: %s Response code: %d", gmErrorMessage, gmErrorCode)
		clientErr := "Failed to get vehicle"
		err = shared.NewAPIError(http.StatusInternalServerError, requestErr, clientErr)
		return
	}

	var level float64
	if vehicleID == 1234 {
		level = float64(33.5)
		fuelLevel = &level
	}

	if vehicleID == 1235 {
		level = float64(88.55)
		batteryLevel = &level
	}

	return
}

// SendVehicleEngineAction ... mocked logic for gm_connector for testing purposes
func (mg *mockGMAPIConnector) SendVehicleEngineAction(vehicleID int64, action string) (res ActionResult, err *shared.APIError) {
	if vehicleID != 1234 && vehicleID != 1235 {
		gmErrorCode := 404
		gmErrorMessage := fmt.Sprintf("Vehicle id: %d not found.", vehicleID)
		requestErr := fmt.Errorf("Failed to POST vehicle engine action to GM, non-200 response: %s Response code: %d", gmErrorMessage, gmErrorCode)
		clientErr := "Failed to get vehicle"
		err = shared.NewAPIError(http.StatusInternalServerError, requestErr, clientErr)
		return
	}

	res = ActionResult{"EXECUTED"}
	return
}
