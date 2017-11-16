package trimet

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTime(t *testing.T) {

	testTime1 := Time(time.Duration(time.Hour))

	sqlInterval, err := testTime1.Value()
	require.NoError(t, err)
	assert.NotNil(t, sqlInterval)

	b, err := testTime1.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "01:00:00", string(b))

	var testTime2 Time
	require.NoError(t, testTime2.UnmarshalText(b))

	dv, err := testTime2.Value()
	assert.Equal(t, "01:00:00", dv)

}

func TestParseDuration(t *testing.T) {
	nilDur := time.Duration(-1)

	tests := []struct {
		input string
		dur   time.Duration
		err   error
	}{
		{
			input: "10:11:12",
			dur:   10*time.Hour + 11*time.Minute + 12*time.Second,
			err:   nil,
		},
		{
			input: "",
			dur:   nilDur,
			err:   nil,
		},
		{
			input: "12:23",
			dur:   nilDur,
			err:   errors.New("gtfs.parseDuration: expected 3 parts, found 2"),
		},
		{
			input: "12:BAD:TEXT",
			dur:   nilDur,
			err:   errors.New("gtfs.parseDuration: strconv.Atoi: parsing \"BAD\": invalid syntax"),
		},
	}

	for _, test := range tests {
		dur, err := parseDuration(test.input)
		assert.Equal(t, err, test.err)
		if test.dur == nilDur {
			assert.Nil(t, dur)
		} else if assert.NotNil(t, dur) {
			assert.Equal(t, *dur, test.dur)
		}
	}

}

func TestNewCalendarDateFromRow(t *testing.T) {
	cd, err := NewCalendarDateFromRow([]string{"U.497", "20180225", "1"})
	require.NoError(t, err)
	assert.Equal(t, "2018-02-25 00:00:00 +0000 UTC", cd.Date.String())
	assert.Equal(t, 1, cd.ExceptionType)
	assert.Equal(t, "U.497", cd.ServiceID)
}

func TestNewTripFromRow(t *testing.T) {
	tr, err := NewTripFromRow([]string{"290", "C.500", "7735337", "0", "9068", "353218"})
	require.NoError(t, err)
	assert.Equal(t, "7735337", tr.ID)
	assert.Equal(t, 0, *tr.DirectionID)
	assert.Equal(t, "353218", *tr.ShapeID)
}

func TestStopTimeFromRow(t *testing.T) {
	st, err := NewStopTimeFromRow([]string{"7718078", "06:56:50", "06:56:50", "198", "10", "45th Ave", "0", "0", "9289.4", "0", "", ""})
	require.NoError(t, err)
	arrTime, err := st.ArrivalTime.MarshalText()
	require.NoError(t, err)
	depTime, err := st.DepartureTime.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "06:56:50", string(arrTime))
	assert.Equal(t, "06:56:50", string(depTime))

	assert.Equal(t, 0, st.PickupType)
	assert.Equal(t, "45th Ave", *st.StopHeadsign)
	assert.Equal(t, 9289.4, *st.ShapeDistTraveled)
}

func TestReadGTFSCSV(t *testing.T) {
	f, err := ReadGTFSCSV("./fixtures/agency.txt")
	assert.NoError(t, err)
	rows, err := f.Read()
	assert.NoError(t, err)

	assert.Len(t, rows, 9)
}
