package logic

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/bsdavidson/trimetric/trimet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchByIds(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	vds := VehicleSQLDataset{DB: db}

	pwd, err := os.Getwd()
	require.NoError(t, err)
	rc, err := ioutil.ReadFile(pwd + "/fixtures/vehicles.json")
	require.NoError(t, err)

	assert.Equal(t, 530984, len(rc))

	var vrs trimet.VehicleResponse

	require.NoError(t, json.Unmarshal(rc, &vrs))

	require.Len(t, vrs.ResultSet.Vehicles, 616)
	for _, v := range vrs.ResultSet.Vehicles {
		var rv trimet.RawVehicleData

		require.NoError(t, json.Unmarshal(v.Data, &rv))
		_, ok := rv["expires"]
		if ok {
			rv["expires"] = time.Now().Add(time.Hour).Unix() * 1000
			// log.Println(time.Unix(int64(val.(float64))/1000, 0).String())
		}
		b, err := json.Marshal(rv)
		require.NoError(t, err)
		v.Data = b
		require.NoError(t, vds.Upsert(&v))
	}

	trv, err := vds.FetchByIDs([]int{101, 110, 2203})
	require.NoError(t, err)

	assert.Len(t, trv, 3)
}
