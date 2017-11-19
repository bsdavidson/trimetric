package trimet

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
	Trip                TripDescriptor                         `json:"trip"`
	Vehicle             VehicleDescriptor                      `json:"vehicle"`
	Position            Position                               `json:"position"`
	CurrentStopSequence uint32                                 `json:"current_stop_sequence"`
	StopID              string                                 `json:"stop_id"`
	CurrentStatus       gtfs.VehiclePosition_VehicleStopStatus `json:"current_status"`
	Timestamp           uint64                                 `json:"timestamp"`
	CongestionLevel     gtfs.VehiclePosition_CongestionLevel   `json:"congestion_level"`
	OccupancyStatus     gtfs.VehiclePosition_OccupancyStatus   `json:"occupancy_status"`
	RouteType           RouteType                              `json:"route_type"`
}

// VehicleDescriptor contains identification information for a vehicle
// performing a trip.
type VehicleDescriptor struct {
	ID    *string `json:"id"`
	Label *string `json:"label"`
	// LicensePlate *string `json:"license_plate"`
}

// Position is a geographic position of a vehicle.
type Position struct {
	Latitude  float32 `json:"lat"`
	Longitude float32 `json:"lng"`
	Bearing   float32 `json:"bearing"`
	Odometer  float64 `json:"odometer"`
	Speed     float32 `json:"speed"`
}

// IsEqual returns true if the two vehicle positions are the same.
func (pvp *VehiclePosition) IsEqual(cvp VehiclePosition) bool {

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

// RequestVehiclePositions contacts the Trimet Vehicles GTFS API and retrieves all vehicles
// updated after the 'since' value. If no 'since' value is specified, it defaults
// to retrieving them all since midnight of the service day.
func RequestVehiclePositions(appID string, since int64) ([]VehiclePosition, error) {
	query := url.Values{}
	query.Set("appID", appID)
	if since > 0 {
		query.Set("since", strconv.FormatInt(since, 10))
	}
	resp, err := http.Get(fmt.Sprintf("%s?%s", VehiclesGTFS, query.Encode()))
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

	var vp []VehiclePosition
	for _, entity := range feed.Entity {
		vehicle := entity.GetVehicle()
		if vehicle.GetTrip() == nil {
			continue
		}
		v := VehiclePosition{
			Trip: TripDescriptor{
				TripID:  vehicle.GetTrip().TripId,
				RouteID: vehicle.GetTrip().RouteId,
			},
			Vehicle: VehicleDescriptor{
				ID:    vehicle.GetVehicle().Id,
				Label: vehicle.GetVehicle().Label,
			},
			Position: Position{
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
