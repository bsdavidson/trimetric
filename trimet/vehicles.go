package trimet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// Trimet API Routes
const (
	Stops    = "https://developer.trimet.org/ws/V1/stops"
	Arrivals = "https://developer.trimet.org/ws/v2/arrivals"
	Vehicles = "https://developer.trimet.org/ws/v2/vehicles"
	Routes   = "https://developer.trimet.org/ws/V1/routeConfig"
)

// VehicleResponse ...
type VehicleResponse struct {
	ResultSet VehicleResultSet `json:"resultSet"`
}

// VehicleResultSet ...
type VehicleResultSet struct {
	QueryTime int64         `json:"queryTime"`
	Vehicles  []VehicleData `json:"vehicle"`
}

// VehicleData ...
type VehicleData struct {
	VehicleID int    `json:"vehicleID"`
	Data      []byte `json:"-"`
}

type VehicleDataAlias VehicleData

// UnmarshalJSON ...
func (t *VehicleData) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, (*VehicleDataAlias)(t)); err != nil {
		return err
	}
	t.Data = make([]byte, len(b))
	copy(t.Data, b)
	return nil
}

// RawVehicleData ...
type RawVehicleData map[string]interface{}

// Scan ...
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

// RawVehicle ...
type RawVehicle struct {
	VehicleID int            `json:"vehicle_id"`
	Data      RawVehicleData `json:"data"`
}

// https://developer.trimet.org/ws/v2/vehicles?appID=65795DCAB40706D335474B716&json=true

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

// RequestGTFS ...
// func RequestGTFS() error {
// 	client := &http.Client{}
// 	url, err := url.Parse("http://developer.trimet.org/ws/gtfs/VehiclePositions?appID=65795DCAB40706D335474B716")
// 	if err != nil {
// 		return err
// 	}
// 	req, err := http.NewRequest("GET", url.String(), nil)
// 	if err != nil {
// 		return err
// 	}
// 		resp, err := client.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return err
// 	}

// 	feed := gtfs.FeedMessage{}
// 	err = proto.Unmarshal(body, &feed)
// 	if err != nil {
// 		return err
// 	}
// 	log.Println("GTFS Entry Count:", len(feed.Entity))
// 	for _, entity := range feed.Entity {

// 		fmt.Printf("----- %+v\n\n", entity.GetVehicle().)
// 	}

// 	return nil
// }
