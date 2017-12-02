package logic

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/bsdavidson/trimetric/trimet"
	postgis "github.com/cridenour/go-postgis"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	vehicleProducerDuplicatesTotal        prometheus.Counter
	vehicleProducerDuplicateMapSize       prometheus.Gauge
	vehicleProducerEncodingErrorsTotal    prometheus.Counter
	vehicleProducerMessageErrorsTotal     prometheus.Counter
	vehicleProducerMessagesTotal          prometheus.Counter
	vehicleProducerProcessDurationSeconds prometheus.Histogram
	vehicleProducerRequestDurationSeconds prometheus.Histogram
	vehicleProducerRequestErrorsTotal     prometheus.Counter
	vehicleProducerRequestItemsTotal      prometheus.Counter
)

func init() {
	vehicleProducerDuplicatesTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "trimetric_vehicle_producer_duplicates_total",
			Help: "Total number of duplicate vehicle positions received",
		},
	)
	vehicleProducerDuplicateMapSize = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "trimetric_vehicle_producer_duplicate_map_size",
			Help: "Size of the map that tracks duplicate vehicles",
		},
	)
	vehicleProducerEncodingErrorsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "trimetric_vehicle_producer_encoding_errors_total",
			Help: "Total number of errors from encoding GTFS realtime vehicle positions",
		},
	)
	vehicleProducerMessageErrorsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "trimetric_vehicle_producer_errors_total",
			Help: "Total number of errors from producing GTFS realtime vehicle positions",
		},
	)
	vehicleProducerMessagesTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "trimetric_vehicle_producer_messages_total",
			Help: "Total messages produced for GTFS realtime vehicle positions",
		},
	)
	vehicleProducerProcessDurationSeconds = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "trimetric_vehicle_producer_duration_seconds",
			Help: "Duration of time to process a batch of GTFS realtime vehicle positions",
		},
	)
	vehicleProducerRequestDurationSeconds = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "trimetric_vehicle_producer_request_duration_seconds",
			Help: "Duration of requests for GTFS realtime vehicle positions",
		},
	)
	vehicleProducerRequestErrorsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "trimetric_vehicle_producer_request_errors_total",
			Help: "Total number of errors from requests for GTFS realtime vehicle positions",
		},
	)
	vehicleProducerRequestItemsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "trimetric_vehicle_producer_request_items_total",
			Help: "Total number of items recieved within responses for GTFS realtime vehicle positions",
		},
	)

	prometheus.MustRegister(vehicleProducerDuplicatesTotal)
	prometheus.MustRegister(vehicleProducerDuplicateMapSize)
	prometheus.MustRegister(vehicleProducerEncodingErrorsTotal)
	prometheus.MustRegister(vehicleProducerMessageErrorsTotal)
	prometheus.MustRegister(vehicleProducerMessagesTotal)
	prometheus.MustRegister(vehicleProducerProcessDurationSeconds)
	prometheus.MustRegister(vehicleProducerRequestDurationSeconds)
	prometheus.MustRegister(vehicleProducerRequestErrorsTotal)
	prometheus.MustRegister(vehicleProducerRequestItemsTotal)
}

// VehiclePositionWithRouteType adds routetype to identify the vehicle type
type VehiclePositionWithRouteType struct {
	trimet.VehiclePosition
	RouteType trimet.RouteType `json:"route_type" msg:"route_type"`
}

// VehicleDataset provides methods to update and retrieve vehicle data
type VehicleDataset interface {
	FetchVehiclePositions(since int) ([]VehiclePositionWithRouteType, error)
	UpsertVehiclePosition(v *trimet.VehiclePosition) error
	UpsertVehiclePositionBytes(ctx context.Context, b []byte) error
}

// VehicleSQLDataset wraps a DB instance that is used to store vehicle data
type VehicleSQLDataset struct {
	DB *sql.DB
}

// FetchVehiclePositions makes a query against the DB and retrieves a list of vehicle data.
// If IDs are passed in, then vehicle data is restricted to those specific
// vehicle IDs. Otherwise, all vehicles with a non-expired timestamp are returned.
func (vd *VehicleSQLDataset) FetchVehiclePositions(since int) ([]VehiclePositionWithRouteType, error) {

	q := `
		SELECT
			v.vehicle_id, v.vehicle_label, v.trip_id, v.route_id, v.position_bearing,
			v.position_lon_lat, v.current_stop_sequence, v.stop_id, v.current_status,
			v.timestamp, COALESCE(r.type, 3) as route_type
		FROM vehicle_positions v
		LEFT OUTER JOIN routes r ON v.route_id = r.id
		WHERE v.timestamp > $1::bigint
		AND v.timestamp > extract(epoch from now() - interval '5 minute')::bigint
	`
	//

	q += ` ORDER BY vehicle_id`
	rows, err := vd.DB.Query(q, since)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	vehicles := []VehiclePositionWithRouteType{}
	for rows.Next() {
		var v VehiclePositionWithRouteType
		var lonLat postgis.PointS
		err := rows.Scan(
			&v.Vehicle.ID, &v.Vehicle.Label, &v.Trip.TripID, &v.Trip.RouteID,
			&v.Position.Bearing, &lonLat, &v.CurrentStopSequence, &v.StopID,
			&v.CurrentStatus, &v.Timestamp, &v.RouteType)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		v.Position.Latitude = float32(lonLat.Y)
		v.Position.Longitude = float32(lonLat.X)
		vehicles = append(vehicles, v)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.WithStack(err)
	}
	return vehicles, nil
}

// UpsertVehiclePosition updates/inserts a vehicle in the DB.
func (vd *VehicleSQLDataset) UpsertVehiclePosition(v *trimet.VehiclePosition) error {
	lonLat := fmt.Sprintf("SRID=4326;POINT(%f %f)", v.Position.Longitude, v.Position.Latitude)
	q := `
		INSERT INTO vehicle_positions (
			trip_id, route_id, vehicle_id, vehicle_label,
			position_lon_lat, position_bearing, current_stop_sequence, stop_id,
			current_status, timestamp
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (vehicle_id) DO UPDATE SET
			trip_id = EXCLUDED.trip_id,
			route_id = EXCLUDED.route_id,
			vehicle_label = EXCLUDED.vehicle_label,
			position_lon_lat = EXCLUDED.position_lon_lat,
			position_bearing = EXCLUDED.position_bearing,
			current_stop_sequence = EXCLUDED.current_stop_sequence,
			stop_id = EXCLUDED.stop_id,
			current_status = EXCLUDED.current_status,
			timestamp = EXCLUDED.timestamp
	`
	_, err := vd.DB.Exec(
		q,
		v.Trip.TripID,
		v.Trip.RouteID,
		v.Vehicle.ID,
		v.Vehicle.Label,
		lonLat,
		v.Position.Bearing,
		v.CurrentStopSequence,
		v.StopID,
		v.CurrentStatus,
		v.Timestamp)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// UpsertVehiclePositionBytes writes decodes bytes into VehiclePositions and
// updates the DB.
func (vd *VehicleSQLDataset) UpsertVehiclePositionBytes(ctx context.Context, b []byte) error {
	var v trimet.VehiclePosition

	o, err := v.UnmarshalMsg(b)
	if err != nil {
		log.Println(err, o)
		return errors.WithStack(err)
	}

	if err := vd.UpsertVehiclePosition(&v); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// ProduceVehiclePositions makes requests to the Trimet API and passes the result
// to a Producer
func ProduceVehiclePositions(ctx context.Context, p Producer, baseURL, apiKey string, delay time.Duration) error {
	var lastQueryTimeMs int64
	vehicleMap := map[string]uint64{}
	ticker := time.NewTicker(delay)
	defer ticker.Stop()
REQUEST_LOOP:
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}

		queryTime := time.Now()
		vehicles, err := trimet.RequestVehiclePositions(baseURL, apiKey, lastQueryTimeMs)
		vehicleProducerRequestDurationSeconds.Observe(time.Since(queryTime).Seconds())
		if err != nil {
			vehicleProducerRequestErrorsTotal.Add(1)
			log.Println(err)
			continue
		}
		vehicleProducerRequestItemsTotal.Add(float64(len(vehicles)))

		t := time.Now()
		for _, tv := range vehicles {
			if val, ok := vehicleMap[*tv.Vehicle.ID]; ok {
				if tv.Timestamp == val {
					vehicleProducerDuplicatesTotal.Add(1)
					continue
				}
			}

			vehicleMap[*tv.Vehicle.ID] = tv.Timestamp
			vehicleProducerDuplicateMapSize.Set(float64(len(vehicleMap)))

			var b []byte
			msgBytes, err := tv.MarshalMsg(b)
			if err != nil {
				vehicleProducerEncodingErrorsTotal.Add(1)
				log.Println(err)
				continue REQUEST_LOOP
			}
			vehicleProducerMessagesTotal.Add(1)
			err = p.Produce(msgBytes)
			if err != nil {
				vehicleProducerMessageErrorsTotal.Add(1)
				log.Println(err)
				continue REQUEST_LOOP
			}

		}

		vehicleProducerProcessDurationSeconds.Observe(time.Since(t).Seconds())
		lastQueryTimeMs = queryTime.Unix() * 1000
	}
}
