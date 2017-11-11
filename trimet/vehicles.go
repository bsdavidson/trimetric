package trimet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gogo/protobuf/proto"
	"github.com/google/gtfs-realtime-bindings/golang/gtfs"
)

// Trimet API Routes
const (
	Stops        = "https://developer.trimet.org/ws/V1/stops"
	Arrivals     = "https://developer.trimet.org/ws/v2/arrivals"
	Vehicles     = "https://developer.trimet.org/ws/v2/vehicles"
	VehiclesGTFS = "http://developer.trimet.org/ws/gtfs/VehiclePositions"
	Routes       = "https://developer.trimet.org/ws/V1/routeConfig"
)

// VehicleResponse is the top level of a Trimet API vehicle response
type VehicleResponse struct {
	ResultSet VehicleResultSet `json:"resultSet"`
}

// VehicleResultSet is the inner wrapper of a Trimet API vehicle response.
type VehicleResultSet struct {
	QueryTime int64         `json:"queryTime"`
	Vehicles  []VehicleData `json:"vehicle"`
}

// VehicleData is a single vehicle in a vehicle response
type VehicleData struct {
	VehicleID int    `json:"vehicleID"`
	Data      []byte `json:"-"`
}

type vehicleDataAlias VehicleData

// UnmarshalJSON sets the VehicleID and stores the rest of the data on the Data field.
// This allows us to pass the bytes to the client for processing.
func (t *VehicleData) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, (*vehicleDataAlias)(t)); err != nil {
		return err
	}
	t.Data = make([]byte, len(b))
	copy(t.Data, b)
	return nil
}

// RawVehicleData represents the raw JSON data for a single vehicle.
type RawVehicleData map[string]interface{}

type GTFSVehiclePosition struct {
	Trip                GTFSTrip                               `json:"trip"`
	Vehicle             GTFSVehicle                            `json:"vehicle"`
	Position            GTFSPosition                           `json:"position"`
	CurrentStopSequence uint32                                 `json:"current_stop_sequence"`
	StopID              string                                 `json:"stop_id"`
	CurrentStatus       gtfs.VehiclePosition_VehicleStopStatus `json:"current_status"`
	Timestamp           uint64                                 `json:"timestamp"`
	CongestionLevel     gtfs.VehiclePosition_CongestionLevel   `json:"congestion_level"`
	OccupancyStatus     gtfs.VehiclePosition_OccupancyStatus   `json:"occupancy_status"`
}

type GTFSTrip struct {
	TripID  string `json:"trip_id"`
	RouteID string `json:"route_id"`
}

type GTFSVehicle struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type GTFSPosition struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Bearing   float32 `json:"bearing"`
	Odometer  float64 `json:"odometer"`
	Speed     float32 `json:"speed"`
}

// IsEqual ...
func (pvp *GTFSVehiclePosition) IsEqual(cvp GTFSVehiclePosition) bool {

	if pvp.Trip != cvp.Trip {
		return false
	}
	if pvp.CurrentStatus != cvp.CurrentStatus {
		return false
	}
	if pvp.CurrentStopSequence != cvp.CurrentStopSequence {
		return false
	}
	if pvp.Position != cvp.Position {
		return false
	}
	if pvp.StopID != cvp.StopID {
		return false
	}
	if pvp.Vehicle != cvp.Vehicle {
		return false
	}
	return true
}

// Scan unmarshals the raw JSON bytes stored in the DB into a map.
func (vd *RawVehicleData) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("RawVehicleData.Scan: src must be []byte, got %T", src)
	}
	if err := json.Unmarshal(b, vd); err != nil {
		return fmt.Errorf("RawVehicleData.Scan: %v", err)
	}
	return nil
}

// RawVehicle wraps raw vehicle data in a struct that contains
// the VehicleID.
type RawVehicle struct {
	VehicleID int            `json:"vehicle_id"`
	Data      RawVehicleData `json:"data"`
}

// RequestVehicles contacts the Trimet Vehicles API and retrieves all vehicles
// updated after the 'since' value. If no 'since' value is specified, it defaults
// to retrieving them all since midnight of the service day.
func RequestVehicles(appID string, since int64) (*VehicleResponse, error) {
	query := url.Values{}
	query.Set("appID", appID)
	query.Set("json", "true")
	if since > 0 {
		query.Set("since", strconv.FormatInt(since, 10))
	}
	resp, err := http.Get(fmt.Sprintf("%s?%s", Vehicles, query.Encode()))
	if err != nil {
		return nil, fmt.Errorf("http.Get: %s", err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %s", err)
	}
	var tr VehicleResponse
	err = json.Unmarshal(b, &tr)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %s", err)
	}
	return &tr, nil
}

// RequestGTFSVehicles contacts the Trimet Vehicles GTFS API and retrieves all vehicles
// updated after the 'since' value. If no 'since' value is specified, it defaults
// to retrieving them all since midnight of the service day.
func RequestGTFSVehicles(appID string, since int64) ([]GTFSVehiclePosition, error) {
	query := url.Values{}
	query.Set("appID", appID)
	if since > 0 {
		query.Set("since", strconv.FormatInt(since, 10))
	}
	resp, err := http.Get(fmt.Sprintf("%s?%s", VehiclesGTFS, query.Encode()))
	if err != nil {
		return nil, fmt.Errorf("trimet.RequestGTFSVehicles: %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("trimet.RequestGTFSVehicles: %s", err)
	}

	feed := gtfs.FeedMessage{}

	err = proto.Unmarshal(body, &feed)
	if err != nil {
		return nil, fmt.Errorf("trimet.RequestGTFSVehicles: %s", err)
	}

	var vp []GTFSVehiclePosition
	for _, entity := range feed.Entity {
		vehicle := entity.GetVehicle()
		v := GTFSVehiclePosition{
			Trip: GTFSTrip{
				TripID:  vehicle.GetTrip().GetTripId(),
				RouteID: vehicle.GetTrip().GetRouteId(),
			},
			Vehicle: GTFSVehicle{
				ID:    vehicle.GetVehicle().GetId(),
				Label: vehicle.GetVehicle().GetLabel(),
			},
			Position: GTFSPosition{
				Latitude:  vehicle.GetPosition().GetLatitude(),
				Longitude: vehicle.GetPosition().GetLongitude(),
				Bearing:   vehicle.GetPosition().GetBearing(),
				Odometer:  vehicle.GetPosition().GetOdometer(),
				Speed:     vehicle.GetPosition().GetSpeed(),
			},
			CurrentStopSequence: vehicle.GetCurrentStopSequence(),
			StopID:              vehicle.GetStopId(),
			CurrentStatus:       vehicle.GetCurrentStatus(),
			Timestamp:           vehicle.GetTimestamp(),
			CongestionLevel:     vehicle.GetCongestionLevel(),
			OccupancyStatus:     vehicle.GetOccupancyStatus(),
		}
		vp = append(vp, v)
	}

	return vp, nil
}
