package logic

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/bsdavidson/trimetric/trimet"
	"github.com/garyburd/redigo/redis"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// PollGTFSData makes periodic queries to fetch static GTFS data from Trimet.
// It stores a timestamp of the last query in Redis and will only download if it has
// been more than 24 hours since the last download.
func PollGTFSData(ctx context.Context, ld LoaderDataset, redisPool *redis.Pool, dur time.Duration) {
	rc := redisPool.Get()
	defer rc.Close()

	lastTs, err := redis.Int64(rc.Do("GET", "poll_gtfs_last_ts"))
	if err != nil {
		log.Println(err)
	}

	timer := time.NewTimer(time.Until(time.Unix(lastTs, 0).Add(dur)))

	for {
		select {
		case <-timer.C:
			timer.Reset(dur)
			_, err := rc.Do("SET", "poll_gtfs_last_ts", time.Now().Unix())
			if err != nil {
				log.Println(err)
			}

			if err := ld.LoadGTFSData(); err != nil {
				log.Println(err)
			}

		case <-ctx.Done():
			break
		}
	}
}

func bulkReplace(tx *sql.Tx, c *trimet.CSV, table string, columns []string, callback func(stmt *sql.Stmt, row []string) error) (int, error) {
	defer c.Close()

	stmt, err := tx.Prepare(pq.CopyIn(table, columns...))
	if err != nil {
		return 0, errors.WithStack(err)
	}
	var count int
	for {
		row, err := c.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			stmt.Exec()
			return 0, errors.WithStack(err)
		}
		count++
		if count == 1 {
			// Skipping first row because it contains column names
			continue
		}

		err = callback(stmt, row)
		if err != nil {
			stmt.Exec()
			return 0, err
		}
	}

	_, err = stmt.Exec()
	if err != nil {
		return 0, errors.WithStack(err)
	}

	if err := stmt.Close(); err != nil {
		return 0, errors.WithStack(err)
	}

	return count, err
}

// LoaderDataset provides a method to bulk load static GTFS data.
type LoaderDataset interface {
	LoadGTFSData() error
}

// LoaderSQLDataset implements LoaderDataset for a SQL database.
type LoaderSQLDataset struct {
	DB *sql.DB
}

type gtfsLoader struct {
	filename   string
	loaderFunc func(tx *sql.Tx, c *trimet.CSV) error
}

// LoadGTFSData downloads the lastes static GTFS data from Trimet, and
// updates the GTFS data in the database.
func (ld *LoaderSQLDataset) LoadGTFSData() error {
	log.Println("downloading GTFS data")
	gtfsFile, err := trimet.RequestGTFSFile()
	if err != nil {
		return err
	}
	log.Println("finished downloading GTFS data")

	loaders := []gtfsLoader{
		{
			filename:   "calendar_dates.txt",
			loaderFunc: ld.LoadServices,
		},
		{
			filename:   "calendar_dates.txt",
			loaderFunc: ld.LoadCalendarDates,
		},
		{
			filename:   "routes.txt",
			loaderFunc: ld.LoadRoutes,
		},
		{
			filename:   "trips.txt",
			loaderFunc: ld.LoadTrips,
		},
		{
			filename:   "stops.txt",
			loaderFunc: ld.LoadStops,
		},
		{
			filename:   "stop_times.txt",
			loaderFunc: ld.LoadStopTimes,
		},
	}

	tx, err := ld.DB.Begin()
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = tx.Exec(`
		TRUNCATE calendar_dates, routes, services, stops, stop_times, trips
	`)
	if err != nil {
		return rollbackError(tx.Rollback(), err)
	}

	for _, l := range loaders {
		data, err := trimet.ReadZippedGTFSCSV(gtfsFile, l.filename)
		if err != nil {
			return rollbackError(tx.Rollback(), err)
		}
		if err := l.loaderFunc(tx, data); err != nil {
			return rollbackError(tx.Rollback(), err)
		}
		log.Printf("loaded %s", l.filename)
	}
	if err := tx.Commit(); err != nil {
		return rollbackError(tx.Rollback(), err)
	}
	return nil
}

// LoadCalendarDates loads calendar_dates.txt.
func (ld *LoaderSQLDataset) LoadCalendarDates(tx *sql.Tx, c *trimet.CSV) error {
	cols := []string{"service_id", "date", "exception_type"}
	n, err := bulkReplace(tx, c, "calendar_dates", cols, func(stmt *sql.Stmt, row []string) error {
		cd, err := trimet.NewCalendarDateFromRow(row)
		if err != nil {
			return err
		}
		_, err = stmt.Exec(cd.ServiceID, cd.Date, cd.ExceptionType)
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	log.Printf("processed %d rows", n)
	return nil
}

// LoadRoutes loads routes.txt.
func (ld *LoaderSQLDataset) LoadRoutes(tx *sql.Tx, c *trimet.CSV) error {
	cols := []string{
		"id", "agency_id", "short_name", "long_name",
		"type", "url", "color", "text_color",
		"sort_order",
	}
	n, err := bulkReplace(tx, c, "routes", cols, func(stmt *sql.Stmt, row []string) error {

		routeType, err := strconv.Atoi(row[4])
		if err != nil {
			return errors.WithStack(err)
		}

		sortOrder, err := strconv.Atoi(row[8])
		if err != nil {
			return errors.WithStack(err)
		}

		_, err = stmt.Exec(
			row[0], row[1], row[2], row[3], row[4], routeType, row[6], row[7],
			sortOrder)
		if err != nil {
			return errors.WithStack(err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	log.Printf("processed %d rows", n)
	return nil
}

// LoadServices populates the services table with the unique set of
// service_ids from the calendar_dates.txt file.
func (ld *LoaderSQLDataset) LoadServices(tx *sql.Tx, c *trimet.CSV) error {
	cols := []string{"id"}
	ids := map[string]struct{}{}

	n, err := bulkReplace(tx, c, "services", cols, func(stmt *sql.Stmt, row []string) error {
		if _, ok := ids[row[0]]; ok {
			return nil
		}
		ids[row[0]] = struct{}{}
		_, err := stmt.Exec(row[0])
		if err != nil {
			return errors.WithStack(err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	log.Printf("processed %d rows", n)
	return nil
}

// LoadStops loads stops.txt.
func (ld *LoaderSQLDataset) LoadStops(tx *sql.Tx, c *trimet.CSV) error {
	cols := []string{
		"id", "code", "name", "desc", "lat_lon", "zone_id", "url",
		"location_type", "parent_station", "direction", "position",
	}
	n, err := bulkReplace(tx, c, "stops", cols, func(stmt *sql.Stmt, row []string) error {
		locType, err := strconv.Atoi(row[8])
		if err != nil {
			return err
		}

		lonLat := fmt.Sprintf("SRID=4326;POINT(%s %s)", row[5], row[4])
		_, err = stmt.Exec(
			row[0], row[1], row[2], row[3], lonLat, row[6], row[7], locType,
			row[9], row[10], row[11])
		if err != nil {
			return errors.WithStack(err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	log.Printf("processed %d rows", n)
	return nil
}

// LoadStopTimes loads stop_times.txt.
func (ld *LoaderSQLDataset) LoadStopTimes(tx *sql.Tx, c *trimet.CSV) error {
	cols := []string{
		"trip_id", "arrival_time", "departure_time", "stop_id",
		"stop_sequence", "stop_headsign", "pickup_type", "drop_off_type",
		"shape_dist_traveled", "timepoint", "continuous_drop_off",
		"continuous_pickup",
	}
	n, err := bulkReplace(tx, c, "stop_times", cols, func(stmt *sql.Stmt, row []string) error {

		st, err := trimet.NewStopTimeFromRow(row)
		if err != nil {
			return err
		}

		_, err = stmt.Exec(
			st.TripID, st.ArrivalTime, st.DepartureTime, st.StopID, st.StopSequence,
			st.StopHeadsign, st.PickupType, st.DropOffType, st.ShapeDistTraveled,
			st.Timepoint, st.ContinuousDropOff, st.ContinuousPickup)
		if err != nil {
			return errors.WithStack(err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	log.Printf("processed %d rows", n)
	return nil
}

// LoadTrips loads trips.txt.
func (ld *LoaderSQLDataset) LoadTrips(tx *sql.Tx, c *trimet.CSV) error {
	cols := []string{
		"route_id", "service_id", "id", "headsign",
		"short_name", "direction_id", "block_id", "shape_id",
		"wheelchair_accessible", "bikes_allowed",
	}
	n, err := bulkReplace(tx, c, "trips", cols, func(stmt *sql.Stmt, row []string) error {
		t, err := trimet.NewTripFromRow(row)
		if err != nil {
			return err
		}

		_, err = stmt.Exec(
			t.RouteID, t.ServiceID, t.ID, t.Headsign, t.ShortName,
			t.DirectionID, t.BlockID, t.ShapeID, t.WheelchairAccessible,
			t.BikesAllowed)
		if err != nil {
			return errors.WithStack(err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	log.Printf("processed %d rows", n)
	return nil
}
