// Package logic provides the business logic for Trimetric. It handles interacting with
// the database and processing data.
package logic

import "context"

// ConsumerFunc defines a generic functioin type to read messages from a
// Producer
type ConsumerFunc func(ctx context.Context, b []byte) error

// Producer implements methods to Produce Kafka messages
type Producer interface {
	Close() error
	Produce(b []byte) error
}
