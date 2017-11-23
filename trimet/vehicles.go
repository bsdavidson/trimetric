package trimet

//go:generate msgp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gogo/protobuf/proto"
	"github.com/google/gtfs-realtime-bindings/golang/gtfs"
	"github.com/pkg/errors"
)

// RouteType indicates the type of vehicle serving the route
type RouteType int

// Defines the types of vehicles serving a route.
const (
	RouteTypeTram RouteType = iota
	RouteTypeSubway
	RouteTypeRail
	RouteTypeBus
	RouteTypeFerry
	RouteTypeCableCar
	RouteTypeGondola
	RouteTypeFunicular
)

// VehiclePosition is the realtime position information for a given vehicle.
type VehiclePosition struct {
	Trip                TripDescriptor    `json:"trip" msg:"trip"`
	Vehicle             VehicleDescriptor `json:"vehicle" msg:"vehicle"`
	Position            Position          `json:"position" msg:"position"`
	CurrentStopSequence uint32            `json:"current_stop_sequence" msg:"current_stop_sequence"`
	StopID              string            `json:"stop_id" msg:"stop_id"`
	CurrentStatus       int32             `json:"current_status" msg:"current_status"`
	Timestamp           uint64            `json:"timestamp" msg:"timestamp"`
	CongestionLevel     int32             `json:"congestion_level" msg:"congestion_level"`
	OccupancyStatus     int32             `json:"occupancy_status" msg:"occupancy_status"`
}

// VehicleDescriptor contains identification information for a vehicle
// performing a trip.
type VehicleDescriptor struct {
	ID    *string `json:"id" msg:"id"`
	Label *string `json:"label" msg:"label"`
	// LicensePlate *string `json:"license_plate"`
}

// Position is a geographic position of a vehicle.
type Position struct {
	Latitude  float32 `json:"lat"  msg:"lat"`
	Longitude float32 `json:"lng"  msg:"lng"`
	Bearing   float32 `json:"bearing"  msg:"bearing"`
	Odometer  float64 `json:"odometer"  msg:"odometer"`
	Speed     float32 `json:"speed"  msg:"speed"`
}

// RequestVehiclePositions contacts the Trimet Vehicles GTFS API and retrieves all vehicles
// updated after the 'since' value. If no 'since' value is specified, it defaults
// to retrieving them all since midnight of the service day.
func RequestVehiclePositions(baseURL string, appID string, since int64) ([]VehiclePosition, error) {
	query := url.Values{}
	query.Set("appID", appID)
	if since > 0 {
		query.Set("since", strconv.FormatInt(since, 10))
	}
	resp, err := http.Get(fmt.Sprintf("%s?%s", baseURL+VehiclesGTFS, query.Encode()))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	feed := gtfs.FeedMessage{}
	err = proto.Unmarshal(body, &feed)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var vps []VehiclePosition
	for _, e := range feed.Entity {
		vp := VehiclePosition{
			CurrentStopSequence: e.Vehicle.GetCurrentStopSequence(),
			StopID:              e.Vehicle.GetStopId(),
			CurrentStatus:       (int32)(e.Vehicle.GetCurrentStatus()),
			Timestamp:           e.Vehicle.GetTimestamp(),
			CongestionLevel:     (int32)(e.Vehicle.GetCongestionLevel()),
			OccupancyStatus:     (int32)(e.Vehicle.GetOccupancyStatus()),
		}
		if e.Vehicle.Trip != nil {
			vp.Trip = TripDescriptor{
				TripID:  e.Vehicle.Trip.TripId,
				RouteID: e.Vehicle.Trip.RouteId,
			}
		}
		if e.Vehicle.Vehicle != nil {
			vp.Vehicle = VehicleDescriptor{
				ID:    e.Vehicle.Vehicle.Id,
				Label: e.Vehicle.Vehicle.Label,
			}
		}
		if e.Vehicle.Position != nil {
			vp.Position = Position{
				Latitude:  e.Vehicle.Position.GetLatitude(),
				Longitude: e.Vehicle.Position.GetLongitude(),
				Bearing:   e.Vehicle.Position.GetBearing(),
				Odometer:  e.Vehicle.Position.GetOdometer(),
				Speed:     e.Vehicle.Position.GetSpeed(),
			}
		}
		vps = append(vps, vp)
	}
	return vps, nil
}
