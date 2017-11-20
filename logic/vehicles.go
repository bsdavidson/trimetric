package logic

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/Shopify/sarama"
	"github.com/bsdavidson/trimetric/trimet"
	postgis "github.com/cridenour/go-postgis"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	vehicleProducerDuplicatesTotal        prometheus.Counter
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
	prometheus.MustRegister(vehicleProducerEncodingErrorsTotal)
	prometheus.MustRegister(vehicleProducerMessageErrorsTotal)
	prometheus.MustRegister(vehicleProducerMessagesTotal)
	prometheus.MustRegister(vehicleProducerProcessDurationSeconds)
	prometheus.MustRegister(vehicleProducerRequestDurationSeconds)
	prometheus.MustRegister(vehicleProducerRequestErrorsTotal)
	prometheus.MustRegister(vehicleProducerRequestItemsTotal)
}

const vehiclePositionsTopic = "vehicle_positions"

// VehicleDataset provides methods to update and retrieve vehicle data
type VehicleDataset interface {
	FetchVehiclePositionsByIDs(ids []int) ([]trimet.VehiclePosition, error)
	UpsertVehiclePosition(v *trimet.VehiclePosition) error
}

// VehicleSQLDataset wraps a DB instance that is used to store vehicle data
type VehicleSQLDataset struct {
	DB *sql.DB
}

// FetchVehiclePositionsByIDs makes a query against the DB and retrieves a list of vehicle data.
// If IDs are passed in, then vehicle data is restricted to those specific
// vehicle IDs. Otherwise, all vehicles with a non-expired timestamp are returned.
func (vd *VehicleSQLDataset) FetchVehiclePositionsByIDs(ids []int) ([]trimet.VehiclePosition, error) {
	q := `
		SELECT
			v.vehicle_id, v.vehicle_label, v.trip_id, v.route_id, v.position_bearing,
			v.position_lon_lat, v.current_stop_sequence, v.stop_id, v.current_status,
			v.timestamp, COALESCE(r.type, 3) as route_type
		FROM vehicle_positions v
		LEFT OUTER JOIN routes r ON v.route_id = r.id
		WHERE v.timestamp > extract(epoch from now() - interval '5 minute')::bigint
	`
	var args []interface{}

	if len(ids) > 0 {
		q += ` AND vehicle_id = ANY($1)`
		args = append(args, pq.Array(ids))
	}

	q += ` ORDER BY vehicle_id`
	rows, err := vd.DB.Query(q, args...)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	var vehicles []trimet.VehiclePosition
	for rows.Next() {
		var v trimet.VehiclePosition
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
	if vehicles == nil {
		vehicles = []trimet.VehiclePosition{}
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

// ProduceVehiclePositions makes requests to the Trimet API and sends the results to
// Kafka.
func ProduceVehiclePositions(ctx context.Context, apiKey string, kafkaAddr string) error {
	producer, err := sarama.NewSyncProducer([]string{kafkaAddr}, nil)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Println(err)
		}
	}()

	var lastQueryTime time.Time
	vehicleMap := map[*string]trimet.VehiclePosition{}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
REQUEST_LOOP:
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}

		queryTime := time.Now()
		vehicles, err := trimet.RequestVehiclePositions(apiKey, lastQueryTime.Unix()*1000)
		vehicleProducerRequestDurationSeconds.Observe(time.Since(queryTime).Seconds())
		if err != nil {
			vehicleProducerRequestErrorsTotal.Add(1)
			log.Println(err)
			continue
		}
		vehicleProducerRequestItemsTotal.Add(float64(len(vehicles)))

		var b bytes.Buffer
		enc := gob.NewEncoder(&b)
		t := time.Now()
		for _, tv := range vehicles {
			if val, ok := vehicleMap[tv.Vehicle.ID]; ok {
				if tv.IsEqual(val) {
					vehicleProducerDuplicatesTotal.Add(1)
					continue
				}
			}
			vehicleMap[tv.Vehicle.ID] = tv
			if err := enc.Encode(tv); err != nil {
				vehicleProducerEncodingErrorsTotal.Add(1)
				log.Println(err)
				continue REQUEST_LOOP
			}
			vehicleProducerMessagesTotal.Add(1)
			_, _, err := producer.SendMessage(&sarama.ProducerMessage{
				Topic: vehiclePositionsTopic,
				Value: sarama.ByteEncoder(b.Bytes()),
			})
			if err != nil {
				vehicleProducerMessageErrorsTotal.Add(1)
				log.Println(err)
				continue REQUEST_LOOP
			}
		}
		vehicleProducerProcessDurationSeconds.Observe(time.Since(t).Seconds())
		lastQueryTime = queryTime
	}
}

// ConsumeVehiclePositions monitors the Kafka 'vehicles' topic for new messages and
// writes them to a DB.
func ConsumeVehiclePositions(ctx context.Context, vd VehicleDataset, influxClient client.Client, kafkaAddr string) error {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	consumer, err := sarama.NewConsumer([]string{kafkaAddr}, config)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Println(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition(vehiclePositionsTopic, 0, sarama.OffsetNewest)
	if err != nil {
		return errors.WithStack(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Println(err)
		}
	}()

	go func() {
		for err := range partitionConsumer.Errors() {
			log.Println(err, err.Partition, err.Topic)
		}
	}()

MESSAGE:
	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-partitionConsumer.Messages():
			var b bytes.Buffer
			var v trimet.VehiclePosition

			decoder := gob.NewDecoder(&b)
			_, err := b.Write(msg.Value)
			if err != nil {
				log.Println(err)
				break
			}
		DECODE:
			for {
				err = decoder.Decode(&v)
				if err == io.EOF {
					break DECODE
				} else if err != nil {
					log.Println(err)
					continue MESSAGE
				}
			}
			if err := vd.UpsertVehiclePosition(&v); err != nil {
				log.Println(err)
				break
			}

		}

	}
}
