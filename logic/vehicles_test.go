package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
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

	vehicles := loadVehicles(t, db)
	assert.Equal(t, 33, len(vehicles))

	vds := VehicleSQLDataset{DB: db}

	for _, v := range vehicles {
		v.Timestamp = uint64(time.Now().Unix())
		require.NoError(t, vds.UpsertVehiclePosition(&v))
	}

	newVehicles, err := vds.FetchVehiclePositionsByIDs([]int{201, 202, 204})
	require.NoError(t, err)

	assert.Equal(t, 3, len(newVehicles))
}

func TestProduceVehicles(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

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

	err = ProduceVehiclePositions(ctx, ts.URL, "123", mp)
	require.NoError(t, err)

	assert.Equal(t, 81, len(mp.bytes))

	vds := VehicleSQLDataset{DB: db}

	for _, b := range mp.bytes {
		err := vds.UpsertVehiclePositionBytes(ctx, b)
		require.NoError(t, err)
	}
}
