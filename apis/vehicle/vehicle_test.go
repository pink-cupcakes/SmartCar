package vehicle

import (
	"os"
	"testing"

	gmConnector "app_api/shared/gm"

	"github.com/stretchr/testify/assert"
)

var vehicleService Service
var testGMAPIConnector gmConnector.GMAPIConnector

func TestMain(m *testing.M) {
	testGMAPIConnector = gmConnector.NewMockGMAPIConnector()
	vehicleService = NewService(testGMAPIConnector)
	os.Exit(m.Run())
}

func TestGetVehicleSuccess(t *testing.T) {
	expectedRes := Vehicle{"123123412412", "Metallic Silver", 4, "v8"}

	res, err := vehicleService.GetVehicle(1234)
	assert.Nil(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestGetVehicleFailureInvalidVehicleID(t *testing.T) {
	_, err := vehicleService.GetVehicle(1236)
	assert.NotNil(t, err)
}

func TestGetVehicleDoorsSuccess(t *testing.T) {
	expectedRes := []gmConnector.GMVehicleDoorData{}
	expectedRes = append(expectedRes, gmConnector.GMVehicleDoorData{"frontLeft", true})
	expectedRes = append(expectedRes, gmConnector.GMVehicleDoorData{"frontRight", true})

	res, err := vehicleService.GetVehicleDoors(1234)
	assert.Nil(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestGetVehicleDoorsFailureInvalidVehicleID(t *testing.T) {
	_, err := vehicleService.GetVehicleDoors(1236)
	assert.NotNil(t, err)
}

func TestGetVehicleFuelSuccess(t *testing.T) {
	level := 33.5
	expectedFuelRes := Fuel{&level}

	fuelRes, err := vehicleService.GetVehicleFuel(1234)
	assert.Nil(t, err)
	assert.Equal(t, expectedFuelRes, fuelRes)
}

func TestGetVehicleFuelNilSuccess(t *testing.T) {
	expectedFuelRes := Fuel{}

	fuelRes, err := vehicleService.GetVehicleFuel(1235)
	assert.Nil(t, err)
	assert.Equal(t, expectedFuelRes, fuelRes)
}

func TestGetVehicleFuelFailureInvalidVehicleID(t *testing.T) {
	_, err := vehicleService.GetVehicleFuel(1236)
	assert.NotNil(t, err)
}

func TestGetVehicleBatterySuccess(t *testing.T) {
	level := 88.55
	expectedBatteryRes := Battery{&level}

	batteryRes, err := vehicleService.GetVehicleBattery(1235)
	assert.Nil(t, err)
	assert.Equal(t, expectedBatteryRes, batteryRes)
}

func TestGetVehicleBatteryNilSuccess(t *testing.T) {
	expectedBatteryRes := Battery{}

	batteryRes, err := vehicleService.GetVehicleBattery(1234)
	assert.Nil(t, err)
	assert.Equal(t, expectedBatteryRes, batteryRes)
}

func TestGetVehicleBatteryFailureInvalidVehicleID(t *testing.T) {
	_, err := vehicleService.GetVehicleBattery(1236)
	assert.NotNil(t, err)
}

func TestSendEngineActionNilSuccess(t *testing.T) {
	expectedRes := EngineActionResponse{"success"}

	engineAction := EngineActionRequest{"START"}
	res, err := vehicleService.SendEngineAction(1234, engineAction)
	assert.Nil(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestSendEngineActionFailureInvalidVehicleID(t *testing.T) {
	engineAction := EngineActionRequest{"START"}

	_, err := vehicleService.SendEngineAction(1236, engineAction)
	assert.NotNil(t, err)
}

func TestSendEngineActionFailureInvalidAction(t *testing.T) {
	engineAction := EngineActionRequest{"FOOBAR"}

	_, err := vehicleService.SendEngineAction(1235, engineAction)
	assert.NotNil(t, err)
}
