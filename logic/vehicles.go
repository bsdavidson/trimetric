package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
	"github.com/bsdavidson/trimetric/trimet"
	"github.com/lib/pq"
)

// VehicleDataset ...
type VehicleDataset interface {
	FetchByIDs(ids []int) ([]trimet.RawVehicle, error)
	Upsert(v *trimet.VehicleData) error
}

// VehicleSQLDataset ...
type VehicleSQLDataset struct {
	DB *sql.DB
}

// FetchByIDs ...
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
		return nil, err
	}

	var vehicles []trimet.RawVehicle
	for rows.Next() {
		var v trimet.RawVehicle
		if err := rows.Scan(&v.VehicleID, &v.Data); err != nil {
			return nil, err
		}
		vehicles = append(vehicles, v)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return vehicles, nil
}

// Upsert ...
func (vd *VehicleSQLDataset) Upsert(v *trimet.VehicleData) error {
	q := `INSERT INTO raw_vehicles (vehicle_id, data)
	VALUES ($1, $2)
	ON CONFLICT (vehicle_id) DO UPDATE SET
		data = EXCLUDED.data;
 `
	_, err := vd.DB.Exec(q, v.VehicleID, v.Data)
	if err != nil {
		return err
	}
	return nil
}

// ProduceVehicles ...
func ProduceVehicles(ctx context.Context, apiKey string) error {
	log.Println("Starting ProduceVehicles")

	producer, err := sarama.NewAsyncProducer([]string{"kafka:9092"}, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	go func() {
		for e := range producer.Errors() {
			log.Println("Producer Error!", e.Error())
		}
	}()

	var lastQueryTime int64

	for {
		select {
		case <-ctx.Done():
			log.Println("Closing")
			return nil
		default:
			time.Sleep(time.Second)

			tr, err := trimet.RequestVehicles(apiKey, lastQueryTime)
			if err != nil {
				log.Println("requestVehicles:", err)
				continue
			}
			lastQueryTime = tr.ResultSet.QueryTime - 1
			// log.Println("Got Vehicles", len(tr.ResultSet.Vehicles))
			for _, tv := range tr.ResultSet.Vehicles {
				producer.Input() <- &sarama.ProducerMessage{Topic: "vehicles", Value: sarama.ByteEncoder(tv.Data)}
			}
		}
	}
}

// ConsumeVehicles monitors the Kafka 'vehicles' topic for new messages and
// writes them to a db.
func ConsumeVehicles(ctx context.Context, vd VehicleDataset) error {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	consumer, err := sarama.NewConsumer([]string{"kafka:9092"}, config)
	if err != nil {
		return fmt.Errorf("consumeVehicles: %v", err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Println("consumeVehicles:", err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition("vehicles", 0, sarama.OffsetNewest)
	if err != nil {
		return fmt.Errorf("consumeVehicles: %v", err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Println("consumeVehicles: ", err)
		}
	}()

	go func() {
		for e := range partitionConsumer.Errors() {
			log.Println("Consumer Error", e.Error(), e.Partition, e.Topic)
		}
	}()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			var v trimet.VehicleData
			if err := json.Unmarshal(msg.Value, &v); err != nil {
				log.Println("ConsumeVehicles:", err)
				break
			}
			if err := vd.Upsert(&v); err != nil {
				log.Println("ConsumeVehicles:", err)
				break
			}

		case <-ctx.Done():
			log.Println("Closing")
			return nil
		}
	}
	return nil
}
