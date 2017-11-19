package logic

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/bsdavidson/trimetric/trimet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
