package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/bsdavidson/trimetric/trimet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockProducer struct {
	T     *testing.T
	topic string
	bytes [][]byte
}

func (m *mockProducer) Close() error {
	return nil
}

func (m *mockProducer) Produce(b []byte) error {
	m.bytes = append(m.bytes, b)
	return nil
}

func loadVehicles(t *testing.T, db *sql.DB) []trimet.VehiclePosition {
	f, err := os.Open("./testdata/vehicles.json")
	require.NoError(t, err)

	b, err := ioutil.ReadAll(f)
	require.NoError(t, err)

	var vehicles []trimet.VehiclePosition
	err = json.Unmarshal(b, &vehicles)
	require.NoError(t, err)

	return vehicles
}

func TestFetchByIds(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	now := uint64(time.Now().Unix())
	vehicles := loadVehicles(t, db)
	assert.Equal(t, 33, len(vehicles))
	vds := VehicleSQLDataset{DB: db}
	sv := make([]trimet.VehiclePosition, 2)
	for _, v := range vehicles {
		if *v.Vehicle.ID == "201" {
			sv[0] = v
			v.Timestamp = now
		} else if *v.Vehicle.ID == "202" {
			sv[1] = v
			v.Timestamp = now
		}

		require.NoError(t, vds.UpsertVehiclePosition(&v))
	}
	log.Println("NOW", now)
	newVehicles, err := vds.FetchVehiclePositions(int(now) - 1)
	require.NoError(t, err)
	assert.Equal(t, 2, len(newVehicles))

	for i, v := range newVehicles {
		assert.Equal(t, *sv[i].Trip.RouteID, *v.Trip.RouteID)
		assert.Equal(t, *sv[i].Trip.TripID, *v.Trip.TripID)
		assert.Equal(t, *sv[i].Vehicle.ID, *v.Vehicle.ID)
		assert.Equal(t, *sv[i].Vehicle.Label, *v.Vehicle.Label)
		assert.Equal(t, sv[i].Position, v.Position)
		assert.Equal(t, sv[i].CurrentStopSequence, v.CurrentStopSequence)
		assert.Equal(t, sv[i].StopID, v.StopID)
		assert.Equal(t, sv[i].CurrentStatus, v.CurrentStatus)
		assert.Equal(t, now, v.Timestamp)
		assert.Equal(t, sv[i].CongestionLevel, v.CongestionLevel)
		assert.Equal(t, sv[i].OccupancyStatus, v.OccupancyStatus)
	}
}

func TestProduceVehicles(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	now := uint64(time.Now().Unix())
	ctx, cancel := context.WithCancel(context.Background())
	var httpCalls int

	pb, err := ioutil.ReadFile("./testdata/vehicle_positions.pb")
	require.NoError(t, err)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpCalls++
		if httpCalls == 2 {
			cancel()
		}
		w.Write(pb)
	}))

	mp := &mockProducer{T: t}

	err = ProduceVehiclePositions(ctx, mp, ts.URL, "123", time.Millisecond)
	require.NoError(t, err)

	assert.Equal(t, 82, len(mp.bytes))

	vds := VehicleSQLDataset{DB: db}

	var expected trimet.VehiclePosition
	for _, b := range mp.bytes {
		vp := trimet.VehiclePosition{}
		_, err := vp.UnmarshalMsg(b)
		require.NoError(t, err)
		if *vp.Vehicle.ID == "3505" {
			expected = vp
			vp.Timestamp = now
		}
		ob, err := vp.MarshalMsg([]byte{})
		require.NoError(t, err)
		err = vds.UpsertVehiclePositionBytes(ctx, ob)
		require.NoError(t, err)
	}

	newVehicles, err := vds.FetchVehiclePositions(int(now) - 100)
	require.NoError(t, err)
	require.Equal(t, 1, len(newVehicles))

	assert.Equal(t, *expected.Trip.RouteID, *newVehicles[0].Trip.RouteID)
	assert.Equal(t, *expected.Trip.TripID, *newVehicles[0].Trip.TripID)
	assert.Equal(t, *expected.Vehicle.ID, *newVehicles[0].Vehicle.ID)
	assert.Equal(t, *expected.Vehicle.Label, *newVehicles[0].Vehicle.Label)
	assert.Equal(t, expected.Position, newVehicles[0].Position)
	assert.Equal(t, expected.CurrentStopSequence, newVehicles[0].CurrentStopSequence)
	assert.Equal(t, expected.StopID, newVehicles[0].StopID)
	assert.Equal(t, expected.CurrentStatus, newVehicles[0].CurrentStatus)
	assert.Equal(t, now, newVehicles[0].Timestamp)
	assert.Equal(t, expected.CongestionLevel, newVehicles[0].CongestionLevel)
	assert.Equal(t, expected.OccupancyStatus, newVehicles[0].OccupancyStatus)
}
