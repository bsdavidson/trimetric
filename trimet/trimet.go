// Package trimet provides methods to interact with the Trimet API.
package trimet

// Trimet API Routes
const (
	GTFS          = "https://developer.trimet.org/schedule/gtfs.zip"
	Stops         = "https://developer.trimet.org/ws/V1/stops"
	Arrivals      = "https://developer.trimet.org/ws/v2/arrivals"
	Vehicles      = "https://developer.trimet.org/ws/v2/vehicles"
	VehiclesGTFS  = "http://developer.trimet.org/ws/gtfs/VehiclePositions"
	Routes        = "https://developer.trimet.org/ws/V1/routeConfig"
	TripUpdateURL = "https://developer.trimet.org/ws/V1/TripUpdate"
)
