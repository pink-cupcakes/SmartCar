package gmapiconnector

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// gmVehicleData ... represents a flattened structure for the relevant vehicle data from GM
type gmVehicleData struct {
	Vin        string `json:"vin"`
	Color      string `json:"color"`
	IsFourDoor bool   `json:"fourDoorSedan"`
	IsTwoDoor  bool   `json:"twoDoorCoupe"`
	DriveTrain string `json:"driveTrain"`
}

/* The API response from GM is structured in nested objects type/value, represented in DataValue struct.
   MapToStruct is a method that can be extended to any struct that will map from {key: DataValue{type, value}} to a flattened {key: value}
   It reliably type checks against the types provided by GM before unmarshalling the results against the internal API struct
*/
func (gmData *gmVehicleData) MapToStruct(data map[string]DataValue) (err error) {
	transformedData := make(map[string]interface{})

	for key, val := range data {
		var convertedVal interface{}

		switch val.Type {
		case "String":
			convertedVal = val.Value
		case "Boolean":
			convertedVal, err = strconv.ParseBool(val.Value)
		case "Number":
			convertedVal, err = strconv.ParseInt(val.Value, 10, 64)
		case "Null":
			convertedVal = nil
		default:
			err = fmt.Errorf("Unsupported data type: %s", val.Type)
		}

		if err != nil {
			return err
		}

		transformedData[key] = convertedVal
	}

	jsonString, err := json.Marshal(transformedData)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(jsonString), &gmData)
	if err != nil {
		return
	}

	return
}

// gmVehicleDoorData ... represents a flattened structure for the relevant vehicle door data from GM
type GMVehicleDoorData struct {
	Location string `json:"location"`
	Locked   bool   `json:"locked"`
}

/* The API response from GM is structured in nested objects type/value, represented in DataValue struct.
   MapToStruct is a method that can be extended to any struct that will map from {key: DataValue{type, value}} to a flattened {key: value}
   It reliably type checks against the types provided by GM before unmarshalling the results against the internal API struct
*/
func (gmData *GMVehicleDoorData) MapToStruct(data map[string]DataValue) (err error) {
	transformedData := make(map[string]interface{})

	for key, val := range data {
		var convertedVal interface{}

		switch val.Type {
		case "String":
			convertedVal = val.Value
		case "Boolean":
			convertedVal, err = strconv.ParseBool(val.Value)
		case "Number":
			convertedVal, err = strconv.ParseInt(val.Value, 10, 64)
		case "Null":
			convertedVal = nil
		default:
			err = fmt.Errorf("Unsupported data type: %s", val.Type)
		}

		if err != nil {
			return err
		}

		transformedData[key] = convertedVal
	}

	jsonString, err := json.Marshal(transformedData)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(jsonString), &gmData)
	if err != nil {
		return
	}

	return
}

// gmVehicleEnergyData ... represents a flattened structure for the relevant vehicle energy data from GM
type GMVehicleEnergyData struct {
	Fuel    *float64 `json:"tankLevel"`
	Battery *float64 `json:"batteryLevel"`
}

/* The API response from GM is structured in nested objects type/value, represented in DataValue struct.
   MapToStruct is a method that can be extended to any struct that will map from {key: DataValue{type, value}} to a flattened {key: value}
   It reliably type checks against the types provided by GM before unmarshalling the results against the internal API struct
*/
func (gmData *GMVehicleEnergyData) MapToStruct(data map[string]DataValue) (err error) {
	transformedData := make(map[string]interface{})

	for key, val := range data {
		var convertedVal interface{}

		switch val.Type {
		case "String":
			convertedVal = val.Value
		case "Boolean":
			convertedVal, err = strconv.ParseBool(val.Value)
		case "Number":
			convertedVal, err = strconv.ParseFloat(val.Value, 64)
		case "Null":
			convertedVal = nil
		default:
			err = fmt.Errorf("Unsupported data type: %s", val.Type)
		}

		if err != nil {
			return err
		}

		transformedData[key] = convertedVal
	}

	jsonString, err := json.Marshal(transformedData)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(jsonString), &gmData)
	if err != nil {
		return
	}

	return
}
