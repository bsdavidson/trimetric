package logic

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/Shopify/sarama"
	"github.com/bsdavidson/trimetric/trimet"
	postgis "github.com/cridenour/go-postgis"
	"github.com/lib/pq"
)

// VehicleDataset provides methods to update and retrieve vehicle data
type VehicleDataset interface {
	FetchByIDs(ids []int) ([]trimet.RawVehicle, error)
	FetchGTFSByIDs(ids []int) ([]trimet.GTFSVehiclePosition, error)
	Upsert(v *trimet.VehicleData) error
	UpsertGTFS(v *trimet.GTFSVehiclePosition) error
}

// VehicleSQLDataset wraps a DB instance that is used to store vehicle data
type VehicleSQLDataset struct {
	DB *sql.DB
}

// FetchByIDs makes a query against the DB and retrieves a list of vehicle data.
// If IDs are passed in, then vehicle data is restricted to those specific
// vehicle IDs. Otherwise, all vehicles with a non-expired timestamp are returned.
func (vd *VehicleSQLDataset) FetchByIDs(ids []int) ([]trimet.RawVehicle, error) {
	q := `SELECT vehicle_id, data
			  FROM raw_vehicles
				WHERE (data->>'expires')::bigint > (extract(epoch from now())*1000)::bigint`
	var args []interface{}

	if len(ids) > 0 {
		q += ` AND vehicle_id = ANY($1)`
		args = append(args, pq.Array(ids))
	}
	q += ` ORDER BY vehicle_id`
	rows, err := vd.DB.Query(q, args...)
	if err != nil {
		return nil, fmt.Errorf("VehicleSQLDataset.FetchByIDs: %v", err)
	}

	var vehicles []trimet.RawVehicle
	for rows.Next() {
		var v trimet.RawVehicle
		if err := rows.Scan(&v.VehicleID, &v.Data); err != nil {
			return nil, fmt.Errorf("VehicleSQLDataset.FetchByIDs: %v", err)
		}
		vehicles = append(vehicles, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("VehicleSQLDataset.FetchByIDs: %v", err)
	}
	return vehicles, nil
}

// FetchGTFSByIDs makes a query against the DB and retrieves a list of vehicle data.
// If IDs are passed in, then vehicle data is restricted to those specific
// vehicle IDs. Otherwise, all vehicles with a non-expired timestamp are returned.
func (vd *VehicleSQLDataset) FetchGTFSByIDs(ids []int) ([]trimet.GTFSVehiclePosition, error) {
	q := `SELECT vehicle_id, vehicle_label, trip_trip_id, trip_route_id,
							 position_bearing, position_lat_lon, current_stop_sequence,
							 stop_id, current_status, timestamp
				FROM gtfs_vehicles`
	var args []interface{}

	if len(ids) > 0 {
		q += ` AND vehicle_id = ANY($1)`
		args = append(args, pq.Array(ids))
	}

	q += ` ORDER BY vehicle_id`
	rows, err := vd.DB.Query(q, args...)
	if err != nil {
		return nil, fmt.Errorf("VehicleSQLDataset.FetchGTFSByIDs: %v", err)
	}

	var vehicles []trimet.GTFSVehiclePosition
	for rows.Next() {
		var v trimet.GTFSVehiclePosition
		var latLon postgis.PointS
		err := rows.Scan(
			&v.Vehicle.ID, &v.Vehicle.Label, &v.Trip.TripID, &v.Trip.RouteID,
			&v.Position.Bearing, &latLon, &v.CurrentStopSequence, &v.StopID,
			&v.CurrentStatus, &v.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("VehicleSQLDataset.FetchGTFSByIDs: %v", err)
		}
		v.Position.Latitude = float32(latLon.Y)
		v.Position.Longitude = float32(latLon.X)
		vehicles = append(vehicles, v)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("VehicleSQLDataset.FetchGTFSByIDs: %v", err)
	}
	return vehicles, nil
}

// Upsert updates/inserts a vehicle in the DB.
func (vd *VehicleSQLDataset) Upsert(v *trimet.VehicleData) error {
	q := `INSERT INTO raw_vehicles (vehicle_id, data)
	VALUES ($1, $2)
	ON CONFLICT (vehicle_id) DO UPDATE SET
		data = EXCLUDED.data;
 `
	_, err := vd.DB.Exec(q, v.VehicleID, v.Data)
	if err != nil {
		return fmt.Errorf("VehicleSQLDataset.Upsert: %v", err)
	}
	return nil
}

// UpsertGTFS updates/inserts a vehicle in the DB.
func (vd *VehicleSQLDataset) UpsertGTFS(v *trimet.GTFSVehiclePosition) error {
	latLon := fmt.Sprintf("SRID=4326;POINT(%f %f)", v.Position.Longitude, v.Position.Latitude)
	q := `INSERT INTO gtfs_vehicles (
					trip_trip_id, trip_route_id, vehicle_id, vehicle_label,
					position_lat_lon, position_bearing, current_stop_sequence, stop_id,
					current_status, timestamp
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
				ON CONFLICT (vehicle_id) DO UPDATE SET
					trip_trip_id = EXCLUDED.trip_trip_id,
					trip_route_id = EXCLUDED.trip_route_id,
					vehicle_label = EXCLUDED.vehicle_label,
					position_lat_lon = EXCLUDED.position_lat_lon,
					position_bearing = EXCLUDED.position_bearing,
					current_stop_sequence = EXCLUDED.current_stop_sequence,
					stop_id = EXCLUDED.stop_id,
					current_status = EXCLUDED.current_status,
					timestamp = EXCLUDED.timestamp
				`
	_, err := vd.DB.Exec(q,
		v.Trip.TripID,
		v.Trip.RouteID,
		v.Vehicle.ID,
		v.Vehicle.Label,
		latLon,
		v.Position.Bearing,
		v.CurrentStopSequence,
		v.StopID,
		v.CurrentStatus,
		v.Timestamp)
	if err != nil {
		return fmt.Errorf("VehicleSQLDataset.UpsertGTFS: %v", err)
	}
	return nil
}

// ProduceVehicles makes requests to the Trimet API and sends the results to
// Kafka.
func ProduceVehicles(ctx context.Context, apiKey string) error {
	log.Println("Starting ProduceVehicles")

	producer, err := sarama.NewAsyncProducer([]string{"kafka:9092"}, nil)
	if err != nil {
		return fmt.Errorf("ProduceVehicles: %v", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Println("ProduceVehicles:", err)
		}
	}()

	go func() {
		for e := range producer.Errors() {
			log.Println("logic.ProduceVehicles: producer error:", e.Error())
		}
	}()

	var lastQueryTime int64

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			time.Sleep(time.Second)

			tr, err := trimet.RequestVehicles(apiKey, lastQueryTime)
			if err != nil {
				log.Println("logic.ProduceVehicles:", err)
				continue
			}
			lastQueryTime = tr.ResultSet.QueryTime - 1
			for _, tv := range tr.ResultSet.Vehicles {
				producer.Input() <- &sarama.ProducerMessage{Topic: "vehicles", Value: sarama.ByteEncoder(tv.Data)}
			}
		}
	}
}

// ProduceGTFSVehicles makes requests to the Trimet API and sends the results to
// Kafka.
func ProduceGTFSVehicles(ctx context.Context, apiKey string) error {
	log.Println("Starting ProduceGTFSVehicles")

	producer, err := sarama.NewAsyncProducer([]string{"kafka:9092"}, nil)
	if err != nil {
		return fmt.Errorf("ProduceGTFSVehicles: %v", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Println("ProduceGTFSVehicles:", err)
		}
	}()

	go func() {
		for e := range producer.Errors() {
			log.Println("logic.ProduceGTFSVehicles: producer error:", e.Error())
		}
	}()

	var lastQueryTime int64
	vehicleMap := map[string]trimet.GTFSVehiclePosition{}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			time.Sleep(1000 * time.Millisecond)

			queryTime := time.Now()
			vehicles, err := trimet.RequestGTFSVehicles(apiKey, lastQueryTime)
			if err != nil {
				log.Println("logic.ProduceGTFSVehicles:", err)
				continue
			}

			lastQueryTime = (queryTime.Unix() * 1000)

			var b bytes.Buffer
			enc := gob.NewEncoder(&b)

		VEHICLE:
			for _, tv := range vehicles {
				if val, ok := vehicleMap[tv.Vehicle.ID]; ok {
					if tv.IsEqual(val) {
						continue VEHICLE
					}
				}
				vehicleMap[tv.Vehicle.ID] = tv
				if err := enc.Encode(tv); err != nil {
					log.Println("logic.ProduceGTFSVehicles:", err)
					continue
				}

				producer.Input() <- &sarama.ProducerMessage{Topic: "vehicles2", Value: sarama.ByteEncoder(b.Bytes())}
			}
		}
	}
}

// ConsumeVehicles monitors the Kafka 'vehicles' topic for new messages and
// writes them to a DB.
func ConsumeVehicles(ctx context.Context, vd VehicleDataset) error {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	consumer, err := sarama.NewConsumer([]string{"kafka:9092"}, config)
	if err != nil {
		return fmt.Errorf("logic.ConsumeVehicles: %v", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Println("logic.ConsumeVehicles:", err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition("vehicles", 0, sarama.OffsetNewest)
	if err != nil {
		return fmt.Errorf("logic.ConsumeVehicles: %v", err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Println("logic.ConsumeVehicles: ", err)
		}
	}()

	go func() {
		for e := range partitionConsumer.Errors() {
			log.Println("logic.ConsumeVehicles: consumer error:", e.Error(), e.Partition, e.Topic)
		}
	}()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var v trimet.VehicleData
			if err := json.Unmarshal(msg.Value, &v); err != nil {
				log.Println("logic.ConsumeVehicles:", err)
				break
			}
			if err := vd.Upsert(&v); err != nil {
				log.Println("logic.ConsumeVehicles:", err)
				break
			}

		case <-ctx.Done():
			return nil
		}
	}
	return nil
}

// ConsumeGTFSVehicles monitors the Kafka 'vehicles' topic for new messages and
// writes them to a DB.
func ConsumeGTFSVehicles(ctx context.Context, vd VehicleDataset) error {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	consumer, err := sarama.NewConsumer([]string{"kafka:9092"}, config)
	if err != nil {
		return fmt.Errorf("logic.ConsumeGTFSVehicles: %v", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Println("logic.ConsumeGTFSVehicles:", err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition("vehicles2", 0, sarama.OffsetNewest)
	if err != nil {
		return fmt.Errorf("logic.ConsumeGTFSVehicles %v", err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Println("logic.ConsumeGTFSVehicles: ", err)
		}
	}()

	go func() {
		for e := range partitionConsumer.Errors() {
			log.Println("logic.ConsumeGTFSVehicles: consumer error:", e.Error(), e.Partition, e.Topic)
		}
	}()

MESSAGE:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var b bytes.Buffer
			var v trimet.GTFSVehiclePosition

			decoder := gob.NewDecoder(&b)
			_, err := b.Write(msg.Value)
			if err != nil {
				log.Println("logic.ConsumeGTFSVehicles: consumer error:", err.Error())
				break
			}
		DECODE:
			for {
				err = decoder.Decode(&v)
				if err == io.EOF {
					break DECODE
				} else if err != nil {
					log.Println("logic.ConsumeGTFSVehicles: consumer error:", err.Error())
					continue MESSAGE
				}
			}

			if err := vd.UpsertGTFS(&v); err != nil {
				log.Println("logic.ConsumeGTFSVehicles:", err)
				break
			}

		case <-ctx.Done():
			return nil
		}
	}

}
