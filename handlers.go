package main

import (
	"net/http"
	"strconv"

	"app_api/apis/vehicle"
	"app_api/shared"
	"app_api/shared/httphelper"

	"github.com/gorilla/mux"
)

// getVehicle ... /vehicles/{vehicle_id} GET
//
// swagger:operation GET /vehicles/{vehicle_id} Vehicles GetVehicle
//
// Returns stats for the requested vehicle
//
// ---
// summary: Returns stats for the requested vehicle
// consumes:
// - application/x-www-form-urlencoded
// - application/json
// produces:
// - application/json
// schemes:
// - https
// parameters:
// - name: vehicle_id
//   in: path
//   description: The vehicle ID number
//   required: true
//   type: integer
// responses:
//   '200':
//     description: >
//       Vehicle object.
//     schema:
//       $ref: "#/definitions/Vehicle"
//   '400':
//     description: "Bad request e.g. Invalid vehicle_id"
//     schema:
//       type: "object"
//       properties:
//         message:
//           type: "string"
//           example: "Vehicle ID must be an integer"
//   '503':
//     description: "Service Unavailable"
//     schema:
//       type: "object"
//       properties:
//         message:
//           type: "string"
//           example: "Internal Error"
func (env *Env) getVehicle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vehicleID, parseErr := strconv.ParseInt(mux.Vars(r)["vehicle_id"], 10, 64)
	if parseErr != nil {
		apiError := shared.NewAPIError(http.StatusBadRequest, parseErr, "Vehicle ID must be an integer").
			SetInternalErrorMessage("Failed to parse vehicle ID")
		httphelper.NewResponse(r.Context(), w, nil, apiError)
		return
	}

	vehicleInfo, apiErr := env.Services.VehicleService.GetVehicle(vehicleID)

	// NewResponse ... is a wrapper that has error checking and logging, and sends a response to the client
	httphelper.NewResponse(ctx, w, vehicleInfo, apiErr)
	return
}

// getVehicleDoors ... /vehicles/{vehicle_id}/doors GET
//
// swagger:operation GET /vehicles/{vehicle_id}/doors Vehicles getVehicleDoors
//
// Returns status of the doors for the requested vehicle
//
// ---
// summary: Returns status of the doors for the requested vehicle
// consumes:
// - application/x-www-form-urlencoded
// - application/json
// produces:
// - application/json
// schemes:
// - https
// parameters:
// - name: vehicle_id
//   in: path
//   description: The vehicle ID number
//   required: true
//   type: integer
// responses:
//   '200':
//     description: >
//       Vehicle object.
//     schema:
//       type: "array"
//       items:
//         $ref: "#/definitions/Door"
//   '400':
//     description: "Bad request e.g. Invalid vehicle_id"
//     schema:
//       type: "object"
//       properties:
//         message:
//           type: "string"
//           example: "Vehicle ID must be an integer"
//   '503':
//     description: "Service Unavailable"
//     schema:
//       type: "object"
//       properties:
//         message:
//           type: "string"
//           example: "Internal Error"
func (env *Env) getVehicleDoors(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vehicleID, parseErr := strconv.ParseInt(mux.Vars(r)["vehicle_id"], 10, 64)
	if parseErr != nil {
		apiError := shared.NewAPIError(http.StatusBadRequest, parseErr, "Vehicle ID must be an integer").
			SetInternalErrorMessage("Failed to parse vehicle ID")
		httphelper.NewResponse(r.Context(), w, nil, apiError)
		return
	}

	vehicleDoorsInfo, apiErr := env.Services.VehicleService.GetVehicleDoors(vehicleID)

	httphelper.NewResponse(ctx, w, vehicleDoorsInfo, apiErr)
	return
}

// getVehicleFuelStatus ... /vehicles/{vehicle_id}/fuel GET
//
// swagger:operation GET /vehicles/{vehicle_id}/fuel Vehicles getVehicleFuelStatus
//
// Returns status of the fuel for the requested vehicle
//
// ---
// summary: Returns status of the fuel for the requested vehicle
// consumes:
// - application/x-www-form-urlencoded
// - application/json
// produces:
// - application/json
// schemes:
// - https
// parameters:
// - name: vehicle_id
//   in: path
//   description: The vehicle ID number
//   required: true
//   type: integer
// responses:
//   '200':
//     description: >
//       Vehicle object.
//     schema:
//       $ref: "#/definitions/Fuel"
//   '400':
//     description: "Bad request e.g. Invalid vehicle_id"
//     schema:
//       type: "object"
//       properties:
//         message:
//           type: "string"
//           example: "Vehicle ID must be an integer"
//   '503':
//     description: "Service Unavailable"
//     schema:
//       type: "object"
//       properties:
//         message:
//           type: "string"
//           example: "Internal Error"
func (env *Env) getVehicleFuelStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vehicleID, parseErr := strconv.ParseInt(mux.Vars(r)["vehicle_id"], 10, 64)
	if parseErr != nil {
		apiError := shared.NewAPIError(http.StatusBadRequest, parseErr, "Vehicle ID must be an integer").
			SetInternalErrorMessage("Failed to parse vehicle ID")
		httphelper.NewResponse(r.Context(), w, nil, apiError)
		return
	}

	vehicleFuelInfo, apiErr := env.Services.VehicleService.GetVehicleFuel(vehicleID)

	httphelper.NewResponse(ctx, w, vehicleFuelInfo, apiErr)
	return
}

// getVehicleBatteryStatus ... /vehicles/{vehicle_id}/battery GET
//
// swagger:operation GET /vehicles/{vehicle_id}/battery Vehicles getVehicleBatteryStatus
//
// Returns status of the battery for the requested vehicle
//
// ---
// summary: Returns status of the battery for the requested vehicle
// consumes:
// - application/x-www-form-urlencoded
// - application/json
// produces:
// - application/json
// schemes:
// - https
// parameters:
// - name: vehicle_id
//   in: path
//   description: The vehicle ID number
//   required: true
//   type: integer
// responses:
//   '200':
//     description: >
//       Vehicle object.
//     schema:
//       $ref: "#/definitions/Battery"
//   '400':
//     description: "Bad request e.g. Invalid vehicle_id"
//     schema:
//       type: "object"
//       properties:
//         message:
//           type: "string"
//           example: "Vehicle ID must be an integer"
//   '503':
//     description: "Service Unavailable"
//     schema:
//       type: "object"
//       properties:
//         message:
//           type: "string"
//           example: "Internal Error"
func (env *Env) getVehicleBatteryStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vehicleID, parseErr := strconv.ParseInt(mux.Vars(r)["vehicle_id"], 10, 64)
	if parseErr != nil {
		apiError := shared.NewAPIError(http.StatusBadRequest, parseErr, "Vehicle ID must be an integer").
			SetInternalErrorMessage("Failed to parse vehicle ID")
		httphelper.NewResponse(r.Context(), w, nil, apiError)
		return
	}

	vehicleBatteryInfo, apiErr := env.Services.VehicleService.GetVehicleBattery(vehicleID)

	httphelper.NewResponse(ctx, w, vehicleBatteryInfo, apiErr)
	return
}

// actionEngine ... /vehicles/{vehicle_id}/engine POST
//
// swagger:operation POST /vehicles/{vehicle_id}/engine Vehicles actionEngine
//
// Submits commands to the engine for the requested vehicle
//
// ---
// summary: Submits commands to the engine for the requested vehicle
// consumes:
// - application/x-www-form-urlencoded
// - application/json
// produces:
// - application/json
// schemes:
// - https
// parameters:
// - name: vehicle_id
//   in: path
//   description: The vehicle ID number
//   required: true
//   type: integer
// - name: body
//   in: body
//   description: body parameters
//   schema:
//     "$ref": "#/definitions/EngineActionRequest"
//   required: true
// responses:
//   '200':
//     description: >
//       Vehicle object.
//     schema:
//       $ref: "#/definitions/EngineActionResponse"
//   '400':
//     description: "Bad request e.g. Invalid vehicle_id"
//     schema:
//       type: "object"
//       properties:
//         message:
//           type: "string"
//           example: "Vehicle ID must be an integer"
//   '503':
//     description: "Service Unavailable"
//     schema:
//       type: "object"
//       properties:
//         message:
//           type: "string"
//           example: "Internal Error"
func (env *Env) actionEngine(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	vehicleID, parseErr := strconv.ParseInt(mux.Vars(r)["vehicle_id"], 10, 64)
	if parseErr != nil {
		apiError := shared.NewAPIError(http.StatusBadRequest, parseErr, "Vehicle ID must be an integer").
			SetInternalErrorMessage("Failed to parse vehicle ID")
		httphelper.NewResponse(r.Context(), w, nil, apiError)
		return
	}

	ea := vehicle.EngineActionRequest{}

	// validate json body
	err := httphelper.DecodeJSONBody(w, r, &ea)
	if err != nil {
		httphelper.NewResponse(r.Context(), w, nil, err)
		return
	}

	engineSubmissionStatus, apiErr := env.Services.VehicleService.SendEngineAction(vehicleID, ea)

	httphelper.NewResponse(ctx, w, engineSubmissionStatus, apiErr)
	return
}
