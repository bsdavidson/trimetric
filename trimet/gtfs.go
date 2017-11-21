package trimet

//go:generate msgp

import (
	"archive/zip"
	"bytes"
	"database/sql/driver"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tinylib/msgp/msgp"

	"github.com/pkg/errors"
)

func init() {
	msgp.RegisterExtension(99, func() msgp.Extension { return new(Time) })
}

// Stop represents a single stop from a GTFS feed.
type Stop struct {
	ID                 string  `json:"id" msg:"id"`
	Code               string  `json:"code" msg:"code"`
	Name               string  `json:"name" msg:"name"`
	Desc               string  `json:"desc" msg:"desc"`
	Lat                float64 `json:"lat" msg:"lat"`
	Lon                float64 `json:"lng" msg:"lng"`
	ZoneID             string  `json:"zone_id" msg:"zone_id"`
	URL                string  `json:"url" msg:"url"`
	LocationType       int     `json:"location_type" msg:"location_type"`
	ParentStation      string  `json:"parent_station" msg:"parent_station"`
	Direction          string  `json:"direction" msg:"direction"`
	Position           string  `json:"position" msg:"position"`
	WheelchairBoarding int     `json:"wheelchair_boarding" msg:"wheelchair_boarding"`
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

// StopTime is a single stop time from a GTFS feed.
type StopTime struct {
	TripID            string   `json:"trip_id" msg:"trip_id"`
	ArrivalTime       *Time    `json:"arrival_time" msg:"arrival_time,extension"`
	DepartureTime     *Time    `json:"departure_time" msg:"departure_time,extension"`
	StopID            string   `json:"stop_id" msg:"stop_id"`
	StopSequence      int      `json:"stop_sequence" msg:"stop_sequence"`
	StopHeadsign      *string  `json:"stop_headsign" msg:"stop_headsign"`
	PickupType        int      `json:"pickup_type" msg:"pickup_type"`
	DropOffType       int      `json:"drop_off_type" msg:"drop_off_type"`
	ShapeDistTraveled *float64 `json:"shape_dist_traveled" msg:"shape_dist_traveled"`
	Timepoint         *int     `json:"timepoint" msg:"timepoint"`
	ContinuousDropOff int      `json:"continuous_drop_off" msg:"continuous_drop_off"`
	ContinuousPickup  int      `json:"continuous_pickup" msg:"continuous_pickup"`
}

// Route represents a single route from a GTFS feed.
type Route struct {
	RouteID   string `json:"id" msg:"id"`
	AgencyID  string `json:"agency_id" msg:"agency_id"`
	ShortName string `json:"short_name" msg:"short_name"`
	LongName  string `json:"long_name" msg:"long_name"`
	Type      int    `json:"type" msg:"type"`
	URL       string `json:"url" msg:"url"`
	Color     string `json:"color" msg:"color"`
	TextColor string `json:"text_color" msg:"text_color"`
	SortOrder int    `json:"sort_order" msg:"sort_order"`
}

// Trip ...
type Trip struct {
	ID                   string  `json:"id" msg:"id"`
	RouteID              string  `json:"route_id" msg:"route_id"`
	ServiceID            string  `json:"service_id" msg:"service_id"`
	DirectionID          *int    `json:"direction_id" msg:"direction_id"`
	BlockID              *string `json:"block_id" msg:"block_id"`
	ShapeID              *string `json:"shape_id" msg:"shape_id"`
	Headsign             *string `json:"headsign" msg:"headsign"`
	ShortName            *string `json:"short_name" msg:"short_name"`
	BikesAllowed         int     `json:"bikes_allowed" msg:"bikes_allowed"`
	WheelchairAccessible int     `json:"wheelchair_accessible" msg:"wheelchair_accessible"`
}

// CalendarDate ...
type CalendarDate struct {
	ServiceID     string    `json:"service_id" msg:"service_id"`
	Date          time.Time `json:"date" msg:"date"`
	ExceptionType int       `json:"exception_type" msg:"exception_type"`
}

func parseDuration(s string) (*Time, error) {
	if s == "" {
		return nil, nil
	}
	parts := strings.Split(s, ":")
	if len(parts) != 3 {
		return nil, errors.Errorf("expected 3 parts, found %d", len(parts))
	}

	var intParts [3]time.Duration
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		intParts[i] = time.Duration(n)
	}
	dur := intParts[0]*time.Hour + intParts[1]*time.Minute + intParts[2]*time.Second
	return (*Time)(&dur), nil
}

func parseFloat(s string, defaultValue float64) (float64, error) {
	if s == "" {
		return defaultValue, nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return f, nil
}

func parseInt(s string, defaultValue int) (int, error) {
	if s == "" {
		return defaultValue, nil
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return n, nil
}

func parseNullableFloat(s string) (*float64, error) {
	if s == "" {
		return nil, nil
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &f, nil
}

func parseNullableInt(s string) (*int, error) {
	if s == "" {
		return nil, nil
	}

	n, err := strconv.Atoi(s)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &n, nil
}

func parseNullableString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

const calendarDateLayout = "20060102"

// NewCalendarDateFromRow takes a single row from processing a
// calendar_dates.txt and creates a CalendarDate.
func NewCalendarDateFromRow(row []string) (*CalendarDate, error) {
	cd := CalendarDate{ServiceID: row[0]}
	var err error
	cd.Date, err = time.Parse(calendarDateLayout, row[1])
	if err != nil {
		return nil, errors.WithStack(err)
	}
	cd.ExceptionType, err = strconv.Atoi(row[2])
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &cd, nil
}

// NewTripFromRow takes a single row from processing a
// trips.txt file and creates a Trip.
func NewTripFromRow(row []string) (*Trip, error) {
	t := Trip{
		ID:        row[2],
		RouteID:   row[0],
		ServiceID: row[1],
		BlockID:   parseNullableString(row[4]),
		ShapeID:   parseNullableString(row[5]),
		// Commented out fields are in the spec but not provided by Trimet.
		// TripHeadsign:  parseNullableString(row[3]),
		// TripShortName: parseNullableString(row[4]),
	}

	var err error

	t.DirectionID, err = parseNullableInt(row[3])
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// t.BikesAllowed, err = parseInt(row[9], 0)
	// if err != nil {
	// 	return nil, errors.WithStack(err)
	// }

	// t.WheelchairAccessible, err = parseInt(row[8], 0)
	// if err != nil {
	// 	return nil, errors.WithStack(err)
	// }

	return &t, nil
}

// NewStopTimeFromRow takes a single row from processing a
// stop_times.txt file and creates a StopTime.
func NewStopTimeFromRow(row []string) (*StopTime, error) {
	st := StopTime{
		TripID: row[0],
		StopID: row[3],
	}
	var err error

	st.ArrivalTime, err = parseDuration(row[1])
	if err != nil {
		return nil, errors.WithStack(err)
	}
	st.DepartureTime, err = parseDuration(row[2])
	if err != nil {
		return nil, errors.WithStack(err)
	}

	st.StopSequence, err = strconv.Atoi(row[4])
	if err != nil {
		return nil, errors.WithStack(err)
	}

	st.StopHeadsign = parseNullableString(row[5])

	st.PickupType, err = parseInt(row[6], 0)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	st.DropOffType, err = parseInt(row[7], 0)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	st.ShapeDistTraveled, err = parseNullableFloat(row[8])
	if err != nil {
		return nil, errors.WithStack(err)
	}

	st.Timepoint, err = parseNullableInt(row[9])
	if err != nil {
		return nil, errors.WithStack(err)
	}

	st.ContinuousDropOff, err = parseInt(row[10], 0)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	st.ContinuousPickup, err = parseInt(row[11], 0)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &st, nil
}

// RequestGTFSFile makes a request to download the current GTFS data set from Trimet.
// It returns an array of stops from the file.
func RequestGTFSFile() (*zip.ReadCloser, error) {
	f, err := ioutil.TempFile("", "tmp")
	if err != nil {
		return nil, errors.Wrap(err, "error creating tmp file")
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()
	resp, err := http.Get(GTFS)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return nil, err
	}
	f.Close()

	z, err := zip.OpenReader(f.Name())
	if err != nil {
		return nil, err
	}
	return z, nil
}

// CSV ...
type CSV struct {
	rc io.ReadCloser `msg:"-"`
	cr *csv.Reader   `msg:"-"`
}

func (c *CSV) Read() ([]string, error) {
	return c.cr.Read()
}

// Close ...
func (c *CSV) Close() error {
	if c.rc == nil {
		return nil
	}
	return c.rc.Close()
}

// ReadGTFSCSV reads a GTFS txt file and returns a CSV object.
func ReadGTFSCSV(filename string) (*CSV, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	cr := csv.NewReader(f)

	cr.ReuseRecord = true
	return &CSV{cr: cr}, nil
}

func getLines(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

// ReadZippedGTFSCSV opens a GTFS.zip file and extracts fileName as a CSV
// object.
func ReadZippedGTFSCSV(z *zip.ReadCloser, fileName string) (*CSV, error) {
	var idx int
	for i, zf := range z.File {
		if zf.Name == fileName {
			idx = i
			break
		}
	}
	rc, err := z.File[idx].Open()
	if err != nil {
		return nil, err
	}

	cr := csv.NewReader(rc)

	cr.ReuseRecord = true
	return &CSV{rc: rc, cr: cr}, nil

}
