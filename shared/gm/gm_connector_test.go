package gmapiconnector

import (
	"encoding/json"
	"fmt"
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
