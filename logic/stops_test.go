package logic

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *sql.DB {
	setupDb, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/?user=postgres&sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer setupDb.Close()
	if err := setupDb.Ping(); err != nil {
		log.Fatal(err)
	}
	_, err = setupDb.Exec("DROP DATABASE IF EXISTS test_trimetric")
	if err != nil {
		log.Fatal(err)
	}
	_, err = setupDb.Exec("CREATE DATABASE test_trimetric")
	if err != nil {
		log.Fatal(err)
	}
	setupDb.Close()

	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/test_trimetric?user=postgres&sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	q1 := `CREATE TABLE raw_vehicles (
		vehicle_id integer NOT NULL PRIMARY KEY,
		data jsonb NOT NULL
	);`
	q2 := `CREATE EXTENSION Postgis`
	q3 := `CREATE TABLE stops (
		id text NOT NULL PRIMARY KEY,
		code text NOT NULL DEFAULT '',
		name text NOT NULL,
		"desc" text NOT NULL DEFAULT '',
		lat_lon geography(POINT) NOT NULL,
		zone_id text NOT NULL DEFAULT '',
		stop_url text NOT NULL DEFAULT '',
		location_type integer NOT NULL DEFAULT 0,
		parent_station text NOT NULL DEFAULT 0,
		direction text NOT NULL DEFAULT '',
		position text NOT NULL DEFAULT '',
		wheelchair_boarding integer NOT NULL DEFAULT 0
	);`

	_, err = db.Query(q1)
	assert.NoError(t, err)
	_, err = db.Query(q2)
	assert.NoError(t, err)
	_, err = db.Query(q3)
	assert.NoError(t, err)

	return db
}

func TestFetchWithinDistance(t *testing.T) {

	db := setupTestDB(t)
	defer db.Close()

	pwd, err := os.Getwd()
	require.NoError(t, err)
	rc, err := ioutil.ReadFile(pwd + "/fixtures/stops.txt")
	require.NoError(t, err)
	r := bytes.NewReader(rc)

	stops, err := csv.NewReader(r).ReadAll()
	require.NoError(t, err)

	sds := StopSQLDataset{DB: db}

	assert.NoError(t, sds.UpsertAll(stops))
	//lat=45.5247402&lng=-122.6787931&distance=100
	stps, err := sds.FetchWithinDistance("45.5247402", "-122.6787931", "500")
	require.NoError(t, err)

	assert.Len(t, stps, 48)

}
