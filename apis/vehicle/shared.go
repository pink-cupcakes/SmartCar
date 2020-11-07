package vehicle

// Vehicle response
//
// swagger:model Vehicle
type Vehicle struct {
	// Vin
	//
	// required: true
	// example: 1213231
	Vin string `json:"vin"`

	// Color
	//
	// required: true
	// example: Metallic Silver
	Color string `json:"color"`

	// DoorCount
	//
	// required: true
	// example: 4
	DoorCount int64 `json:"doorCount"`

	// DriveTrain
	//
	// required: true
	// example: v8
	DriveTrain string `json:"driveTrain"`
}

// Door doors response
//
// swagger:model Door
type Door struct {
	// Location
	//
	// required: true
	// example: frontLeft
	Location string `json:"location"`

	// Locked
	//
	// required: true
	// example: Metallic Silver
	Locked bool `json:"locked"`
}
