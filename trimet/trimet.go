// Package trimet provides methods to interact with the Trimet API.
package trimet

// Trimet API Routes
const (
	BaseTrimetURL = "https://developer.trimet.org"
	GTFS          = "/schedule/gtfs.zip"
	Stops         = "/ws/V1/stops"
	Arrivals      = "/ws/v2/arrivals"
	Vehicles      = "/ws/v2/vehicles"
	VehiclesGTFS  = "/ws/gtfs/VehiclePositions"
	Routes        = "/ws/V1/routeConfig"
	TripUpdateURL = "/ws/V1/TripUpdate"
)
