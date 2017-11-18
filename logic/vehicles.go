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
)

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
func ProduceVehiclePositions(ctx context.Context, apiKey string, influxClient client.Client, kafkaAddr string) error {
	log.Println("starting ProduceVehiclePositions")

	producer, err := sarama.NewAsyncProducer([]string{kafkaAddr}, nil)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Println(err)
		}
	}()

	go func() {
		for err := range producer.Errors() {
			log.Println(err)
		}
	}()

	var lastQueryTime int64
	vehicleMap := map[*string]trimet.VehiclePosition{}

	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "trimetric",
		Precision: "s",
	})
	if err != nil {
		log.Println(err)
	}
	tags := map[string]string{"trimet_vehicle": "updated_count"}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			time.Sleep(1000 * time.Millisecond)

			queryTime := time.Now()
			vehicles, err := trimet.RequestVehiclePositions(apiKey, lastQueryTime)
			if err != nil {
				log.Println(err)
				continue
			}

			// log.Printf("Retrieved %d vehicles at %v\n", len(vehicles), lastQueryTime)
			var b bytes.Buffer
			enc := gob.NewEncoder(&b)

			fields := map[string]interface{}{
				"count": len(vehicles),
			}

			pt, err := client.NewPoint("retrieved_vehicles", tags, fields, time.Now())
			if err != nil {
				log.Println(err)
			}
			bp.AddPoint(pt)
			if err := influxClient.Write(bp); err != nil {
				log.Println(err)
			}

		VEHICLE:
			for _, tv := range vehicles {
				if val, ok := vehicleMap[tv.Vehicle.ID]; ok {
					if tv.IsEqual(val) {
						continue VEHICLE
					}
				}
				vehicleMap[tv.Vehicle.ID] = tv
				if err := enc.Encode(tv); err != nil {
					log.Println(err)
					continue
				}

				producer.Input() <- &sarama.ProducerMessage{Topic: vehiclePositionsTopic, Value: sarama.ByteEncoder(b.Bytes())}
				lastQueryTime = (queryTime.Unix() * 1000)

			}
		}
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

		case <-ctx.Done():
			return nil
		}
	}

}
