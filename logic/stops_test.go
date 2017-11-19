package logic

import (
	"database/sql"
	"testing"

	"github.com/bsdavidson/trimetric/trimet"
	"github.com/pressly/goose"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *sql.DB {
	setupDb, err := sql.Open("postgres", "postgres://trimetric:example@localhost:5432/?sslmode=disable")
	require.NoError(t, err)
	defer setupDb.Close()
	require.NoError(t, setupDb.Ping())
	_, err = setupDb.Exec("DROP DATABASE IF EXISTS test_trimetric")
	require.NoError(t, err)
	_, err = setupDb.Exec("CREATE DATABASE test_trimetric")
	require.NoError(t, err)
	setupDb.Close()

	db, err := sql.Open("postgres", "postgres://trimetric:example@localhost:5432/test_trimetric?sslmode=disable")
	require.NoError(t, err)

	require.NoError(t, goose.Up(db, "../migrations"))
	return db
}

func loadStopFixtures(t *testing.T, db *sql.DB) *StopSQLDataset {
	stops, err := trimet.ReadGTFSCSV("./fixtures/stops.txt")
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
