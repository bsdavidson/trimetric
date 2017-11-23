package logic

import (
	"database/sql"
	"flag"
	"fmt"
	"testing"

	"github.com/bsdavidson/trimetric/trimet"
	"github.com/pressly/goose"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var postgresAddr string

func init() {

	flag.StringVar(&postgresAddr, "postgres-addr", "localhost:5432", "Address of a Postgres server")
}

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://trimetric:example@%s/test_trimetric?sslmode=disable", postgresAddr))
	require.NoError(t, err)
	require.NoError(t, goose.Up(db, "../migrations"))
	tables := []string{
		"services",
		"calendar_dates",
		"routes",
		"trips",
		"stops",
		"stop_times",
		"vehicle_positions",
		"trip_updates",
		"stop_time_updates",
	}
	for _, tbl := range tables {
		_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", tbl))
		require.NoError(t, err)
	}
	return db
}

func loadStopFixtures(t *testing.T, db *sql.DB) *StopSQLDataset {
	stops, err := trimet.ReadGTFSCSV("./testdata/stops.txt")
	require.NoError(t, err)

	lds := LoaderSQLDataset{DB: db}
	sds := StopSQLDataset{DB: db}
	tx, err := db.Begin()
	require.NoError(t, err)

	require.NoError(t, lds.LoadStops(tx, stops))
	require.NoError(t, tx.Commit())
	return &sds
}

func TestFetchWithinDistance(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	sds := loadStopFixtures(t, db)
	stops, err := sds.FetchWithinDistance("45.5247402", "-122.6787931", "500")
	require.NoError(t, err)

	assert.Len(t, stops, 48)
}
