package gmapiconnector

import (
	"app_api/shared"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var testGMAPIConnector GMAPIConnector

func TestMain(m *testing.M) {
	testGMAPIConnector = NewGMAPIConnector()
	os.Exit(m.Run())
}

func TestGetVehicleSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testGMVehicleResponse := `{
		"service": "getVehicleInfo",
		"status": "200",
		"data": {
		  "vin": {
			"type": "String",
			"value": "123123412412"
		  },
		  "color": {
			"type": "String",
			"value": "Metallic Silver"
		  },
		  "fourDoorSedan": {
			"type": "Boolean",
			"value": "True"
		  },
		  "twoDoorCoupe": {
			"type": "Boolean",
			"value": "False"
		  },
		  "driveTrain": {
			"type": "String",
			"value": "v8"
		  }
		}
	}`

	var testGmVehicle GMVehicleResponse
	unmarshalTestResponseErr := json.Unmarshal([]byte(testGMVehicleResponse), &testGmVehicle)
	if unmarshalTestResponseErr != nil {
		t.Errorf("Failed to generate test response")
		return
	}

	// Exact URL match
	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/%s", gmAPIURL, getVehicle),
		httpmock.NewStringResponder(200, testGMVehicleResponse))

	var expectedRes gmVehicleData
	mapToStructErr := expectedRes.MapToStruct(testGmVehicle.Data)
	if mapToStructErr != nil {
		t.Errorf("Failed to generate test success flattened vehicle")
		return
	}

	res, err := testGMAPIConnector.GetVehicle(1234)
	assert.Nil(t, err, "GetVehicle success")
	assert.Equal(t, expectedRes, res)
}

func TestGetVehicleFailureGMResponse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	errCode := 400
	errReason := "Required field 'id' not found."

	requestErr := fmt.Errorf("Failed to GET vehicle from GM, non-200 response: %s Response code: %d", errReason, errCode)
	clientErr := "Failed to get vehicle"
	expectedErr := shared.NewAPIError(http.StatusInternalServerError, requestErr, clientErr)

	testGMVehicleResponse := fmt.Sprintf("{\"service\": \"getVehicleInfo\",\"status\": \"%d\",\"reason\": \"%s\"}", errCode, errReason)

	var testGmVehicle GMVehicleResponse
	unmarshalTestResponseErr := json.Unmarshal([]byte(testGMVehicleResponse), &testGmVehicle)
	if unmarshalTestResponseErr != nil {
		t.Errorf("Failed to generate test response")
		return
	}

	// Exact URL match
	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/%s", gmAPIURL, getVehicle),
		httpmock.NewStringResponder(200, testGMVehicleResponse))

	_, err := testGMAPIConnector.GetVehicle(0)
	assert.NotNil(t, err, "GetVehicle failure, GM response with missing id")

	// Ideally, this can just compare equality against the error objects, but because the error objects incorporate line numbers, it's unfeasible to compare the error objects directly.
	// This is a temporary comparison and should be refactored (perhaps even just wrapping these equalities into a shared test helper to compare error objects)
	assert.Equal(t, expectedErr.ErrorMessage, err.ErrorMessage)
	assert.Equal(t, expectedErr.ClientErrorMessage, err.ClientErrorMessage)
	assert.Equal(t, expectedErr.ErrorCode, err.ErrorCode)
}

func TestGetVehicleDoorsSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testGMVehicleDoorsResponse := `{
		"service": "getSecurityStatus",
		"status": "200",
		"data": {
			"doors": {
				"type": "Array",
				"values": [
					{
						"location": {
							"type": "String",
							"value": "frontLeft"
						},
							"locked": {
							"type": "Boolean",
							"value": "False"
						}
					},
					{
						"location": {
							"type": "String",
							"value": "frontRight"
						},
						"locked": {
							"type": "Boolean",
							"value": "True"
						}
					},
					{
						"location": {
							"type": "String",
							"value": "backLeft"
						},
							"locked": {
							"type": "Boolean",
							"value": "False"
						}
					},
					{
						"location": {
							"type": "String",
							"value": "backRight"
						},
							"locked": {
							"type": "Boolean",
							"value": "True"
						}
					}
				]
			}
		}
	}`

	var testGmVehicleDoors GMVehicleDoorsResponse
	unmarshalTestResponseErr := json.Unmarshal([]byte(testGMVehicleDoorsResponse), &testGmVehicleDoors)
	if unmarshalTestResponseErr != nil {
		t.Errorf("Failed to generate test response")
		return
	}

	// Exact URL match
	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/%s", gmAPIURL, getVehicleDoors),
		httpmock.NewStringResponder(200, testGMVehicleDoorsResponse))

	var expectedRes []GMVehicleDoorData
	for _, gmDoorResponse := range testGmVehicleDoors.Data.Doors.Values {
		var flattenedGMDoorResponse GMVehicleDoorData

		mapToStructErr := flattenedGMDoorResponse.MapToStruct(gmDoorResponse)
		if mapToStructErr != nil {
			t.Errorf("Failed to generate success doors test response")
			return
		}

		expectedRes = append(expectedRes, flattenedGMDoorResponse)
	}

	res, err := testGMAPIConnector.GetVehicleDoors(1234)
	assert.Nil(t, err, "GetVehicle doors success")
	assert.Equal(t, expectedRes, res)
}

func TestGetVehicleDoorsFailureGMResponse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	errCode := 400
	errReason := "Required field 'id' not found."

	requestErr := fmt.Errorf("Failed to GET vehicle doors from GM, non-200 response: %s Response code: %d", errReason, errCode)
	clientErr := "Failed to get vehicle doors"
	expectedErr := shared.NewAPIError(http.StatusInternalServerError, requestErr, clientErr)

	testGMVehicleDoorsResponse := fmt.Sprintf("{\"service\": \"getSecurityStatus\",\"status\": \"%d\",\"reason\": \"%s\"}", errCode, errReason)

	var testGmVehicleDoors GMVehicleDoorsResponse
	unmarshalTestResponseErr := json.Unmarshal([]byte(testGMVehicleDoorsResponse), &testGmVehicleDoors)
	if unmarshalTestResponseErr != nil {
		t.Errorf("Failed to generate test response")
		return
	}

	// Exact URL match
	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/%s", gmAPIURL, getVehicleDoors),
		httpmock.NewStringResponder(200, testGMVehicleDoorsResponse))

	_, err := testGMAPIConnector.GetVehicleDoors(123)
	assert.NotNil(t, err, "GetVehicleDoors failure, GM response with missing id")
	assert.Equal(t, expectedErr.ErrorMessage.Error(), err.ErrorMessage.Error())
	assert.Equal(t, expectedErr.ClientErrorMessage, err.ClientErrorMessage)
	assert.Equal(t, expectedErr.ErrorCode, err.ErrorCode)
}

func TestGetVehicleEnergyStatusSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testGMVehicleEnergyResponse := `{
		"service": "getEnergy",
		"status": "200",
		"data": {
			"tankLevel": {
				"type": "Number",
				"value": "30.2"
			},
			"batteryLevel": {
				"type": "Null",
				"value": "null"
			}
		}
	}`

	var testGmVehicleEnergy GMVehicleEnergyResponse
	unmarshalTestResponseErr := json.Unmarshal([]byte(testGMVehicleEnergyResponse), &testGmVehicleEnergy)
	if unmarshalTestResponseErr != nil {
		t.Errorf("Failed to generate test response")
		return
	}

	// Exact URL match
	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/%s", gmAPIURL, getVehicleEnergyLevel),
		httpmock.NewStringResponder(200, testGMVehicleEnergyResponse))

	var expectedRes GMVehicleEnergyData
	mapToStructErr := expectedRes.MapToStruct(testGmVehicleEnergy.Data)
	if mapToStructErr != nil {
		t.Errorf("Failed to generate test success flattened vehicle energy")
		return
	}

	fuel, battery, err := testGMAPIConnector.GetVehicleEnergyStatus(1234)
	assert.Nil(t, err, "GetVehicleEnergyStatus success")
	assert.Equal(t, expectedRes.Fuel, fuel)
	assert.Equal(t, expectedRes.Battery, battery)
}

func TestGetVehicleEnergyStatusFailureGMResponse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	errCode := 400
	errReason := "Required field 'id' not found."

	requestErr := fmt.Errorf("Failed to GET vehicle energy status from GM, non-200 response: %s Response code: %d", errReason, errCode)
	clientErr := "Failed to get vehicle energy status"
	expectedErr := shared.NewAPIError(http.StatusInternalServerError, requestErr, clientErr)

	testGMVehicleEnergyResponse := fmt.Sprintf("{\"service\": \"getSecurityStatus\",\"status\": \"%d\",\"reason\": \"%s\"}", errCode, errReason)

	var testGmVehicleEnergy GMVehicleEnergyResponse
	unmarshalTestResponseErr := json.Unmarshal([]byte(testGMVehicleEnergyResponse), &testGmVehicleEnergy)
	if unmarshalTestResponseErr != nil {
		t.Errorf("Failed to generate test response")
		return
	}

	// Exact URL match
	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/%s", gmAPIURL, getVehicleEnergyLevel),
		httpmock.NewStringResponder(200, testGMVehicleEnergyResponse))

	_, _, err := testGMAPIConnector.GetVehicleEnergyStatus(123)
	assert.NotNil(t, err, "GetVehicleEnergyStatus failure, GM response with missing id")
	assert.Equal(t, expectedErr.ErrorMessage.Error(), err.ErrorMessage.Error())
	assert.Equal(t, expectedErr.ClientErrorMessage, err.ClientErrorMessage)
	assert.Equal(t, expectedErr.ErrorCode, err.ErrorCode)
}

func TestSendVehicleEngineActionSuccess(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testGMVehicleEngineResponse := `{
		"service": "actionEngine",
		"status": "200",
		"actionResult": {
			"status": "EXECUTED"
		}
	}`

	var testGmVehicleEngineResponse GMEngineActionResponse
	unmarshalTestResponseErr := json.Unmarshal([]byte(testGMVehicleEngineResponse), &testGmVehicleEngineResponse)
	if unmarshalTestResponseErr != nil {
		t.Errorf("Failed to generate test response")
		return
	}

	// Exact URL match
	httpmock.RegisterResponder("POST", fmt.Sprintf("%s/%s", gmAPIURL, postVehicleEngineAction),
		httpmock.NewStringResponder(200, testGMVehicleEngineResponse))

	res, err := testGMAPIConnector.SendVehicleEngineAction(1234, "start")
	assert.Nil(t, err, "SendVehicleEngineAction success")
	assert.Equal(t, testGmVehicleEngineResponse.Result, res)
}
