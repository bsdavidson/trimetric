package logic

import (
	"database/sql"
	"testing"

	"github.com/pressly/goose"
	_ "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	setupDb, err := sql.Open("postgres", "postgres://postgres:example@localhost:5432/?user=postgres&sslmode=disable")
	require.NoError(t, err)
	defer setupDb.Close()
	require.NoError(t, setupDb.Ping())
	_, err = setupDb.Exec("DROP DATABASE IF EXISTS test_trimetric")
	require.NoError(t, err)
	_, err = setupDb.Exec("CREATE DATABASE test_trimetric")
	require.NoError(t, err)
	setupDb.Close()

	db, err := sql.Open("postgres", "postgres://postgres:example@localhost:5432/test_trimetric?user=postgres&sslmode=disable")
	require.NoError(t, err)

	require.NoError(t, goose.Up(db, "../migrations"))
	return db
}

func TestFetchWithinDistance(t *testing.T) {

	// db := setupTestDB(t)
	// defer db.Close()

	// pwd, err := os.Getwd()
	// require.NoError(t, err)
	// stops, err := trimet.ReadGTFSCSV(pwd + "/fixtures/stops.txt")
	// require.NoError(t, err)

	// lds := LoaderSQLDataset{DB: db}
	// sds := StopSQLDataset{DB: db}

	// assert.NoError(t, lds.LoadStops(stops))
	// //lat=45.5247402&lng=-122.6787931&distance=100
	// stps, err := sds.FetchWithinDistance("45.5247402", "-122.6787931", "500")
	// require.NoError(t, err)

	// assert.Len(t, stps, 48)

}
