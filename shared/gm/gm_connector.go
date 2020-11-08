package gmapiconnector

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"app_api/shared"
)

// GMAPIConnector ... is an interface of appapi methods called
type GMAPIConnector interface {
	GetVehicle(vehicleID int64) (res gmVehicleData, err *shared.APIError)
	GetVehicleDoors(vehicleID int64) (res []GMVehicleDoorData, err *shared.APIError)
	GetVehicleEnergyStatus(vehicleID int64) (fuelLevel, batteryLevel *float64, err *shared.APIError)
	SendVehicleEngineAction(vehicleID int64, action string) (res ActionResult, err *shared.APIError)
}

type gmAPIConnector struct{}

const (
	gmAPIURL         = "http://gmapi.azurewebsites.net"
	jsonResponseType = "JSON"

	getVehicle              = "getVehicleInfoService"
	getVehicleDoors         = "getSecurityStatusService"
	getVehicleEnergyLevel   = "getEnergyService"
	postVehicleEngineAction = "actionEngineService"

	ENGINE_START = "START_VEHICLE"
	ENGINE_STOP  = "STOP_VEHICLE"
	EXECUTED     = "EXECUTED"
	FAILED       = "FAILED"
)

// NewGMAPIConnector ... returns an interface of GMAPIConnector
func NewGMAPIConnector() GMAPIConnector {
	return &gmAPIConnector{}
}

type GMVehicleResponse struct {
	StatusString string               `json:"status"`
	ErrorMessage string               `json:"reason"`
	Data         map[string]DataValue `json:"data"`
}

type DataValue struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// GetVehicle ... returns an overview for a given car from GM API
func (gm *gmAPIConnector) GetVehicle(vehicleID int64) (res gmVehicleData, err *shared.APIError) {
	requestBody, requestBodyErr := json.Marshal(map[string]interface{}{
		"id":           fmt.Sprintf("%d", vehicleID),
		"responseType": jsonResponseType,
	})
	if requestBodyErr != nil {
		err = shared.NewAPIError(http.StatusInternalServerError, requestBodyErr, "Internal Error").
			SetInternalErrorMessage("GetVehicle: Failed to marshal request body")
		return
	}

	// Make the request to GM to get vehicle information
	resp, requestErr := gm.makeRequest(getVehicle, "POST", requestBody, nil)
	if requestErr != nil {
		clientErr := "Internal Error"
		err = shared.NewAPIError(http.StatusInternalServerError, requestErr, clientErr).SetInternalErrorMessage("GetVehicle: Failed to send request")
		return
	}

	defer resp.Body.Close()

	b, requestErr := ioutil.ReadAll(resp.Body)
	if requestErr != nil {
		clientErr := "Failed to get vehicle"
		err = shared.NewAPIError(http.StatusInternalServerError, requestErr, clientErr).SetInternalErrorMessage("GetVehicle: Failed to read response body from get vehicles")
		return
	}

	// Check if the response coming from GM for failed requests
	if resp.StatusCode != 200 {
		requestErr = errors.New("Failed to GET vehicle from GM, non-200 response: " + string(b))
		clientErr := "Failed to get vehicle"
		err = shared.NewAPIError(http.StatusInternalServerError, requestErr, clientErr)
		return
	}

	var gmVehicleResponse GMVehicleResponse

	requestErr = json.Unmarshal(b, &gmVehicleResponse)
	if requestErr != nil {
		clientErr := "Failed to get vehicle"
		err = shared.NewAPIError(http.StatusInternalServerError, requestErr, clientErr)
		err.SetInternalErrorMessage("GetVehicle: Failed to unmarshal GM vehicle result")
		return
	}

	/* GM's API responses don't surface errors from the API sercice in the response code. Errors messages are surfaced in the API response.
	   This will parse the error code in GM's API response and validate the request for errors
	*/
	gmStatusCode, parseCodeError := strconv.ParseInt(gmVehicleResponse.StatusString, 10, 64)
	if parseCodeError != nil {
		err = shared.NewAPIError(http.StatusInternalServerError, parseCodeError, "Failed to parse status code")
		return
	}
	if gmStatusCode != 200 {
		requestErr = errors.New(fmt.Sprintf("Failed to GET vehicle from GM, non-200 response: %s Response code: %d", gmVehicleResponse.ErrorMessage, gmStatusCode))
		clientErr := "Failed to get vehicle"
		err = shared.NewAPIError(http.StatusInternalServerError, requestErr, clientErr)
		return
	}

	mapToStructErr := res.MapToStruct(gmVehicleResponse.Data)
	if mapToStructErr != nil {
		clientErr := "Failed to get vehicle"
		err = shared.NewAPIError(http.StatusInternalServerError, mapToStructErr, clientErr).SetInternalErrorMessage(fmt.Sprintf("GetVehicle: Failed to parse GM API structured data from GM response. GM response is: %s", string(b)))
		return
	}

	return
}

type GMVehicleDoorsResponse struct {
	StatusString string `json:"status"`
	ErrorMessage string `json:"reason"`
	Data         Doors  `json:"data"`
}

type Doors struct {
	Doors DoorsArrayDataValue `json:"doors"`
}

type DoorsArrayDataValue struct {
	Type   string                 `json:"type"`
	Values []map[string]DataValue `json:"values"`
}

// GetVehicleDoors ... returns the status of the doors for a given car from GM API
func (gm *gmAPIConnector) GetVehicleDoors(vehicleID int64) (res []GMVehicleDoorData, err *shared.APIError) {
	requestBody, requestBodyErr := json.Marshal(map[string]interface{}{
		"id":           fmt.Sprintf("%d", vehicleID),
		"responseType": jsonResponseType,
	})
	if requestBodyErr != nil {
		err = shared.NewAPIError(http.StatusInternalServerError, requestBodyErr, "Internal Error").
			SetInternalErrorMessage("GetVehicleDoors: Failed to marshal request body")
		return
	}

	resp, requestErr := gm.makeRequest(getVehicleDoors, "POST", requestBody, nil)
	if requestErr != nil {
		clientErr := "Internal Error"
		err = shared.NewAPIError(http.StatusInternalServerError, requestErr, clientErr).SetInternalErrorMessage("GetVehicleDoors: Failed to send request")
		return
	}

	defer resp.Body.Close()

	b, requestErr := ioutil.ReadAll(resp.Body)
	if requestErr != nil {
		clientErr := "Failed to get vehicle doors"
		err = shared.NewAPIError(http.StatusInternalServerError, requestErr, clientErr).SetInternalErrorMessage("GetVehicleDoors: Failed to read response body from get vehicles")
		return
	}

	if resp.StatusCode != 200 {
		requestErr = errors.New("Failed to GET vehicle from GM, non-200 response: " + string(b))
		clientErr := "Failed to get vehicle doors"
		err = shared.NewAPIError(http.StatusInternalServerError, requestErr, clientErr)
		return
	}

	var gmVehicleDoorsResponse GMVehicleDoorsResponse

	requestErr = json.Unmarshal(b, &gmVehicleDoorsResponse)
	if requestErr != nil {
		clientErr := "Failed to get vehicle doors"
		err = shared.NewAPIError(http.StatusInternalServerError, requestErr, clientErr)
		err.SetInternalErrorMessage("GetVehicleDoors: Failed to unmarshal GM vehicle result")
		return
	}

	gmStatusCode, parseCodeError := strconv.ParseInt(gmVehicleDoorsResponse.StatusString, 10, 64)
	if parseCodeError != nil {
		err = shared.NewAPIError(http.StatusInternalServerError, parseCodeError, "Failed to parse status code")
		return
	}

	if gmStatusCode != 200 {
		requestErr = fmt.Errorf("Failed to GET vehicle doors from GM, non-200 response: %s Response code: %d", gmVehicleDoorsResponse.ErrorMessage, gmStatusCode)
		clientErr := "Failed to get vehicle doors"
		err = shared.NewAPIError(http.StatusInternalServerError, requestErr, clientErr)
		return
	}

	data := gmVehicleDoorsResponse.Data.Doors

	if data.Type != "Array" {
		requestErr = fmt.Errorf("Incorrect data type from GM API for vehicle doors. Response is \n%s", string(b))
		clientErr := "Failed to get vehicle doors"
		err = shared.NewAPIError(http.StatusInternalServerError, requestErr, clientErr)
		return
	}

	for _, gmDoorResponse := range data.Values {
		var flattenedGMDoorResponse GMVehicleDoorData

		mapToStructErr := flattenedGMDoorResponse.MapToStruct(gmDoorResponse)
		if mapToStructErr != nil {
			clientErr := "Failed to get vehicle doors"
			err = shared.NewAPIError(http.StatusInternalServerError, mapToStructErr, clientErr).SetInternalErrorMessage(fmt.Sprintf("GetVehicleDoors: Failed to parse GM API structured data from GM response. GM response is: %s", string(b)))
			return
		}

		res = append(res, flattenedGMDoorResponse)
	}

	return
}

type GMVehicleEnergyResponse struct {
	StatusString string               `json:"status"`
	ErrorMessage string               `json:"reason"`
	Data         map[string]DataValue `json:"data"`
}

// GetVehicleEnergyStatus ... returns the status of the remaining energy for a given car from GM API
func (gm *gmAPIConnector) GetVehicleEnergyStatus(vehicleID int64) (fuelLevel, batteryLevel *float64, err *shared.APIError) {
	requestBody, requestBodyErr := json.Marshal(map[string]interface{}{
		"id":           fmt.Sprintf("%d", vehicleID),
		"responseType": jsonResponseType,
	})
	if requestBodyErr != nil {
		err = shared.NewAPIError(http.StatusInternalServerError, requestBodyErr, "Internal Error").
			SetInternalErrorMessage("GetVehicleEnergyStatus: Failed to marshal request body")
		return
	}

	resp, requestErr := gm.makeRequest(getVehicleEnergyLevel, "POST", requestBody, nil)
	if requestErr != nil {
		clientErr := "Internal Error"
		err = shared.NewAPIError(http.StatusInternalServerError, requestErr, clientErr).SetInternalErrorMessage("GetVehicleEnergyStatus: Failed to send request")
		return
	}

	defer resp.Body.Close()

	b, requestErr := ioutil.ReadAll(resp.Body)
	if requestErr != nil {
		clientErr := "Failed to get vehicle energy status"
		err = shared.NewAPIError(http.StatusInternalServerError, requestErr, clientErr).SetInternalErrorMessage("GetVehicleEnergyStatus: Failed to read response body from get vehicle energy")
		return
	}

	if resp.StatusCode != 200 {
		requestErr = errors.New("Failed to GET vehicle from GM, non-200 response: " + string(b))
		clientErr := "Failed to get vehicle energy status"
		err = shared.NewAPIError(http.StatusInternalServerError, requestErr, clientErr)
		return
	}

	var gmVehicleEnergyResponse GMVehicleEnergyResponse

	requestErr = json.Unmarshal(b, &gmVehicleEnergyResponse)
	if requestErr != nil {
		clientErr := "Failed to get vehicle energy status"
		err = shared.NewAPIError(http.StatusInternalServerError, requestErr, clientErr)
		err.SetInternalErrorMessage("GetVehicleEnergyStatus: Failed to unmarshal GM vehicle result")
		return
	}

	gmStatusCode, parseCodeError := strconv.ParseInt(gmVehicleEnergyResponse.StatusString, 10, 64)
	if parseCodeError != nil {
		err = shared.NewAPIError(http.StatusInternalServerError, parseCodeError, "Failed to parse status code")
		return
	}

	if gmStatusCode != 200 {
		requestErr = fmt.Errorf("Failed to GET vehicle energy status from GM, non-200 response: %s Response code: %d", gmVehicleEnergyResponse.ErrorMessage, gmStatusCode)
		clientErr := "Failed to get vehicle energy status"
		err = shared.NewAPIError(http.StatusInternalServerError, requestErr, clientErr)
		return
	}

	var flattenedGMDoorResponse GMVehicleEnergyData

	mapToStructErr := flattenedGMDoorResponse.MapToStruct(gmVehicleEnergyResponse.Data)
	if mapToStructErr != nil {
		clientErr := "Failed to get vehicle energy levels"
		err = shared.NewAPIError(http.StatusInternalServerError, mapToStructErr, clientErr).SetInternalErrorMessage(fmt.Sprintf("GetVehicleEnergyStatus: Failed to parse GM API structured data from GM response. GM response is: %s", string(b)))
		return
	}

	return flattenedGMDoorResponse.Fuel, flattenedGMDoorResponse.Battery, nil
}

type GMEngineActionResponse struct {
	StatusString string       `json:"status"`
	ErrorMessage string       `json:"reason"`
	Result       ActionResult `json:"actionResult"`
}

type ActionResult struct {
	Status string `json:"status"`
}

// SendVehicleEngineAction ... returns the status of the remaining energy for a given car from GM API
func (gm *gmAPIConnector) SendVehicleEngineAction(vehicleID int64, action string) (res ActionResult, err *shared.APIError) {
	requestBody, requestBodyErr := json.Marshal(map[string]interface{}{
		"id":           fmt.Sprintf("%d", vehicleID),
		"command":      action,
		"responseType": jsonResponseType,
	})
	if requestBodyErr != nil {
		err = shared.NewAPIError(http.StatusInternalServerError, requestBodyErr, "Internal Error").
			SetInternalErrorMessage("SendVehicleEngineAction: Failed to marshal request body")
		return
	}

	resp, requestErr := gm.makeRequest(postVehicleEngineAction, "POST", requestBody, nil)
	if requestErr != nil {
		clientErr := "Internal Error"
		err = shared.NewAPIError(http.StatusInternalServerError, requestErr, clientErr).SetInternalErrorMessage("SendVehicleEngineAction: Failed to send request")
		return
	}

	defer resp.Body.Close()

	b, requestErr := ioutil.ReadAll(resp.Body)
	if requestErr != nil {
		clientErr := "Failed to send engine action"
		err = shared.NewAPIError(http.StatusInternalServerError, requestErr, clientErr).SetInternalErrorMessage("SendVehicleEngineAction: Failed to read response body from get vehicle energy")
		return
	}

	if resp.StatusCode != 200 {
		requestErr = errors.New("Failed to POST vehicle engine action to GM, non-200 response: " + string(b))
		clientErr := "Failed to send engine action"
		err = shared.NewAPIError(http.StatusInternalServerError, requestErr, clientErr)
		return
	}

	var gmVehicleEngineResponse GMEngineActionResponse

	requestErr = json.Unmarshal(b, &gmVehicleEngineResponse)
	if requestErr != nil {
		clientErr := "Failed to send engine action"
		err = shared.NewAPIError(http.StatusInternalServerError, requestErr, clientErr)
		err.SetInternalErrorMessage("SendVehicleEngineAction: Failed to unmarshal GM vehicle result")
		return
	}

	gmStatusCode, parseCodeError := strconv.ParseInt(gmVehicleEngineResponse.StatusString, 10, 64)
	if parseCodeError != nil {
		err = shared.NewAPIError(http.StatusInternalServerError, parseCodeError, "Failed to parse status code")
		return
	}

	if gmStatusCode != 200 {
		requestErr = fmt.Errorf("Failed to POST vehicle engine action to GM, non-200 response: %s Response code: %d", gmVehicleEngineResponse.ErrorMessage, gmStatusCode)
		clientErr := "Failed to send engine action"
		err = shared.NewAPIError(http.StatusInternalServerError, requestErr, clientErr)
		return
	}

	return gmVehicleEngineResponse.Result, nil
}

func (gm *gmAPIConnector) makeRequest(endpoint, method string, body []byte, params url.Values) (resp *http.Response, err error) {
	client := &http.Client{}
	URL, err := url.Parse(fmt.Sprintf("%s/%s", gmAPIURL, endpoint))
	if err != nil {
		return nil, err
	}

	if params != nil {
		URL.RawQuery = params.Encode()
	}

	req, err := http.NewRequest(method, URL.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
