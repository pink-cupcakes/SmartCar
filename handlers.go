package main

import (
	"net/http"
	"strconv"

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
//           example: "Bad Request: Vehicle not found"
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
//           example: "Bad Request: Vehicle not found"
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
