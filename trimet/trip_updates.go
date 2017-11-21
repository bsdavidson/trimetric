package trimet

//go:generate msgp

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/google/gtfs-realtime-bindings/golang/gtfs"
	"github.com/pkg/errors"
)

// StopTimeEvent contains timing information for a single predicted event.
type StopTimeEvent struct {
	Delay       *int32     `json:"delay" msg:"delay"`
	Time        *time.Time `json:"time" msg:"time"`
	Uncertainty *int32     `json:"uncertainty" msg:"uncertainty"`
}

// StopTimeUpdate is a realtime update for arrival and/or departure events for a
// given stop on a trip.
type StopTimeUpdate struct {
	StopSequence         *uint32       `json:"stop_sequence" msg:"stop_sequence"`
	StopID               *string       `json:"stop_id" msg:"stop_id"`
	Arrival              StopTimeEvent `json:"arrival" msg:"arrival"`
	Departure            StopTimeEvent `json:"departure" msg:"departure"`
	ScheduleRelationship *int32        `json:"schedule_relationship" msg:"schedule_relationship"`
}

// TripDescriptor identifies an instance of a GTFS trip, or all instances of a
// trip along a route.
type TripDescriptor struct {
	TripID  *string `json:"trip_id" msg:"trip_id"`
	RouteID *string `json:"route_id" msg:"route_id"`
	// DirectionID          *uint32                                   `json:"direction_id"`
	// StartTime            *string                                   `json:"start_time"`
	// StartDate            *string                                   `json:"start_date"`
	// ScheduleRelationship *gtfs.TripDescriptor_ScheduleRelationship `json:"schedule_relationship"`
}

// TripUpdate is a realtime update on the progress of a vehile along a trip.
type TripUpdate struct {
	Trip            TripDescriptor    `json:"trip" msg:"trip"`
	Vehicle         VehicleDescriptor `json:"vehicle" msg:"vehicle"`
	StopTimeUpdates []StopTimeUpdate  `json:"stop_time_update" msg:"stop_time_update"`
	Timestamp       *time.Time        `json:"timestamp" msg:"timestamp"`
	Delay           *int32            `json:"delay" msg:"delay"`
}

// TripUpdatesMsg ...
type TripUpdatesMsg struct {
	TripUpdates []TripUpdate `json:"trip_updates" msg:"trip_update"`
}

// RequestTripUpdate makes a request to the trimet TripUpdate GTFS  API endpoint.
func RequestTripUpdate(apiKey string) ([]TripUpdate, error) {
	query := url.Values{}
	query.Set("appID", apiKey)

	resp, err := http.Get(fmt.Sprintf("%s?%s", TripUpdateURL, query.Encode()))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	feed := gtfs.FeedMessage{}
	err = proto.Unmarshal(b, &feed)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var tripUpdates []TripUpdate
	for _, fe := range feed.Entity {
		if fe.TripUpdate == nil {
			continue
		}

		var stopTimeUpdates []StopTimeUpdate
		for _, stu := range fe.TripUpdate.StopTimeUpdate {

			var arrival, departure StopTimeEvent
			if stu.Arrival != nil {

				arrival = StopTimeEvent{
					Delay:       stu.Arrival.Delay,
					Uncertainty: stu.Arrival.Uncertainty,
				}
				if stu.Arrival.Time != nil {
					t := time.Unix(*stu.Arrival.Time, 0)
					arrival.Time = &t
				}
			}
			if stu.Departure != nil {
				departure = StopTimeEvent{
					Delay:       stu.Departure.Delay,
					Uncertainty: stu.Departure.Uncertainty,
				}
				if stu.Departure.Time != nil {
					t := time.Unix(*stu.Departure.Time, 0)
					departure.Time = &t
				}
			}
			stopTimeUpdates = append(stopTimeUpdates, StopTimeUpdate{
				ScheduleRelationship: (*int32)(stu.ScheduleRelationship),
				StopSequence:         stu.StopSequence,
				StopID:               stu.StopId,
				Arrival:              arrival,
				Departure:            departure,
			})
		}

		var trip TripDescriptor
		var vehicle VehicleDescriptor
		if fe.TripUpdate.Trip != nil {
			trip = TripDescriptor{
				TripID:  fe.TripUpdate.Trip.TripId,
				RouteID: fe.TripUpdate.Trip.RouteId,
			}
		}
		if fe.TripUpdate.Vehicle != nil {
			vehicle = VehicleDescriptor{
				ID:    fe.TripUpdate.Vehicle.Id,
				Label: fe.TripUpdate.Vehicle.Label,
			}
		}
		tu := TripUpdate{
			Trip:            trip,
			Vehicle:         vehicle,
			StopTimeUpdates: stopTimeUpdates,
			Delay:           fe.TripUpdate.Delay,
		}
		if fe.TripUpdate.Timestamp != nil {
			t := time.Unix(int64(*fe.TripUpdate.Timestamp), 0)
			tu.Timestamp = &t
		}
		tripUpdates = append(tripUpdates, tu)
	}

	return tripUpdates, nil
}
