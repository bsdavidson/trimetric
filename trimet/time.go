package trimet

import (
	"database/sql/driver"
	"fmt"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/tinylib/msgp/msgp"
)

func init() {
	msgp.RegisterExtension(99, func() msgp.Extension { return new(Time) })
}

// Time is represented in the GTFS feeds as a duration of time since midnight.
// Note that for trips that start the previous day and end past midnight, Time
// can go past 24:00:00.
type Time time.Duration

func (i *Time) ExtensionType() int8 {
	log.Println("Extension called")
	return 99
}
func (i *Time) Len() int {
	return i.Len()
}
func (i *Time) MarshalBinaryTo(b []byte) error {
	msgp.AppendInt64(b, int64(*i))
	return nil
}

func (i *Time) UnmarshalBinary(b []byte) error {
	n, _, err := msgp.ReadInt64Bytes(b)
	if err != nil {
		return errors.WithStack(err)
	}
	*i = Time(n)
	return nil
}

// Scan converts a SQL interval into a Time object.
func (i *Time) Scan(src interface{}) error {
	b, ok := src.([]byte)
	if !ok {
		return errors.Errorf("expected []byte, got %T", src)
	}

	ni, err := parseDuration(string(b))
	if err != nil {
		return err
	}
	*i = *ni
	return nil
}

// Value converts a Time object into a SQL interval string value.
func (i *Time) Value() (driver.Value, error) {
	if i == nil {
		return nil, nil
	}
	b, err := i.MarshalText()
	if err != nil {
		return nil, err
	}
	return driver.Value(string(b)), nil
}

// MarshalText converts a time into text.
func (i *Time) MarshalText() ([]byte, error) {
	d := time.Duration(*i)
	h := int(d / time.Hour)
	m := int((d % time.Hour) / time.Minute)
	s := int((d % time.Minute) / time.Second)
	return []byte(fmt.Sprintf("%02d:%02d:%02d", h, m, s)), nil
}

// UnmarshalText converts text into a Time value.
// It expects the time to be in HH:MM:SS format.
func (i *Time) UnmarshalText(b []byte) error {
	ni, err := parseDuration(string(b))
	if err != nil {
		return err
	}
	*i = *ni
	return nil
}
