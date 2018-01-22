package logic

import (
	"context"
	"log"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
)

// Kafka topics
const (
	TripUpdatesTopic      = "trip_updates"
	VehiclePositionsTopic = "vehicle_positions"
)

// KafkaProducer provides a struct with methods that implement a Producer.
type KafkaProducer struct {
	producer sarama.SyncProducer
	topic    string
}

// NewKafkaProducer returns a new KafkaProducer
func NewKafkaProducer(topic string, addrs []string) (*KafkaProducer, error) {
	p, err := sarama.NewSyncProducer(addrs, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &KafkaProducer{producer: p, topic: topic}, nil
}

// Close wraps a sarama Close method to provide a generic method that satisifies
// the Producer interface.
func (k *KafkaProducer) Close() error {
	return errors.WithStack(k.producer.Close())
}

// Produce wraps a sarama Producer to provide a generic method that satisifies
// the Producer interface.
func (k *KafkaProducer) Produce(b []byte) error {
	_, _, err := k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: k.topic,
		Value: sarama.ByteEncoder(b),
	})
	return errors.WithStack(err)
}

// ConsumeKafkaTopic reads from a Kafka partition and writes the messages
// into a ConsumerFunc.
func ConsumeKafkaTopic(ctx context.Context, c ConsumerFunc, topic string, addrs []string) error {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	client, err := sarama.NewClient(addrs, config)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Println(err)
		}
	}()
	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Println(err)
		}
	}()

	offsetManager, err := sarama.NewOffsetManagerFromClient("trimetric", client)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		if err := offsetManager.Close(); err != nil {
			log.Println(err)
		}
	}()

	partManager, err := offsetManager.ManagePartition(topic, 0)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		if err := partManager.Close(); err != nil {
			log.Println(err)
		}
	}()

	offset, _ := partManager.NextOffset()
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, offset)
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

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-partitionConsumer.Messages():
			if err := c(ctx, msg.Value); err != nil {
				log.Println(err)
				break
			}
			partManager.MarkOffset(msg.Offset+1, "")
		}

	}

}
