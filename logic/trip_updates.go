package logic

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/bsdavidson/trimetric/trimet"
	"github.com/pkg/errors"
)

const tripUpdateDuration = 5000 * time.Millisecond

// TripUpdatesDataset provides methods to update and retrieve trip update data
type TripUpdatesDataset interface {
	UpdateTripUpdates(tus []trimet.TripUpdate) error
	UpdateTripUpdateBytes(ctx context.Context, b []byte) error
	FetchTripUpdates() ([]trimet.TripUpdate, error)
}

// TripUpdateSQLDataset wraps a DB instance that is used to store trip update data
type TripUpdateSQLDataset struct {
	DB *sql.DB
}

// FetchTripUpdates return Tripupdates from the DB
func (tuds *TripUpdateSQLDataset) FetchTripUpdates() ([]trimet.TripUpdate, error) {
	tx, err := tuds.DB.Begin()
	if err != nil {
		return nil, rollbackError(tx.Rollback(), err)
	}

	// Because trip_updates contains references to stop_time_updates, we would
	// be transferring a lot unnecessary data if we queried them both with a JOIN.
	// Instead, make two queries that are joined together after the fact. Because
	// there is a race condition between the two selects, we must ensure that both
	// queries see the same snapshot of the database.
	_, err = tuds.DB.Exec(`SET TRANSACTION ISOLATION LEVEL REPEATABLE READ`)
	if err != nil {
		return nil, rollbackError(tx.Rollback(), err)
	}

	rows, err := tuds.DB.Query(`
		SELECT
			id, trip_id, route_id, vehicle_id, vehicle_label, timestamp, delay
		FROM trip_updates
	`)
	if err != nil {
		return nil, rollbackError(tx.Rollback(), err)
	}

	defer rows.Close()

	var tripUpdates []trimet.TripUpdate
	tuIndex := map[int64]int{}
	for rows.Next() {

		var tu trimet.TripUpdate
		var id int64
		err := rows.Scan(
			&id, &tu.Trip.TripID, &tu.Trip.RouteID, &tu.Vehicle.ID, &tu.Vehicle.Label,
			&tu.Timestamp, &tu.Delay)
		if err != nil {
			return nil, rollbackError(tx.Rollback(), err)
		}

		tuIndex[id] = len(tripUpdates)
		tripUpdates = append(tripUpdates, tu)
	}

	rows, err = tuds.DB.Query(`
		SELECT
			trip_update_id, stop_sequence, stop_id, arrival_delay,
			arrival_time, arrival_uncertainty,departure_delay, departure_time,
			departure_uncertainty, schedule_relationship
		FROM stop_time_updates
		ORDER BY index ASC
	`)
	if err != nil {
		return nil, rollbackError(tx.Rollback(), err)
	}

	for rows.Next() {
		var stu trimet.StopTimeUpdate
		var id int64
		err := rows.Scan(
			&id, &stu.StopSequence, &stu.StopID, &stu.Arrival.Delay,
			&stu.Arrival.Time, &stu.Arrival.Uncertainty, &stu.Departure.Delay,
			&stu.Departure.Time, &stu.Departure.Uncertainty, &stu.ScheduleRelationship)
		if err != nil {
			return nil, rollbackError(tx.Rollback(), err)
		}
		tripUpdates[tuIndex[id]].StopTimeUpdates = append(tripUpdates[tuIndex[id]].StopTimeUpdates, stu)
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.WithStack(err)
	}
	return tripUpdates, nil
}

// UpdateTripUpdates updates trip data in the db.
func (tuds *TripUpdateSQLDataset) UpdateTripUpdates(tus []trimet.TripUpdate) error {
	tx, err := tuds.DB.Begin()
	if err != nil {
		return errors.WithStack(err)
	}

	delStmt, err := tx.Prepare(`DELETE FROM trip_updates`)
	_, err = delStmt.Exec()
	if err != nil {
		return rollbackError(tx.Rollback(), err)
	}
	if err := delStmt.Close(); err != nil {
		return rollbackError(tx.Rollback(), err)
	}

	tuStmt, err := tx.Prepare(`
		INSERT INTO trip_updates (
			trip_id, route_id, vehicle_id, vehicle_label, timestamp, delay
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`)
	if err != nil {
		return rollbackError(tx.Rollback(), err)
	}

	stuStmt, err := tx.Prepare(`
		INSERT INTO stop_time_updates (
			trip_update_id, index, stop_sequence, stop_id, arrival_delay,
			arrival_time, arrival_uncertainty,departure_delay, departure_time,
			departure_uncertainty, schedule_relationship
	  ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`)
	if err != nil {
		return rollbackError(tx.Rollback(), err)
	}

	for _, tu := range tus {
		var id int64
		err = tuStmt.QueryRow(
			tu.Trip.TripID, tu.Trip.RouteID, tu.Vehicle.ID, tu.Vehicle.Label,
			tu.Timestamp, tu.Delay).Scan(&id)
		if err != nil {
			return rollbackError(tx.Rollback(), err)
		}

		for i, stu := range tu.StopTimeUpdates {
			_, err := stuStmt.Exec(
				id, i, stu.StopSequence, stu.StopID, stu.Arrival.Delay,
				stu.Arrival.Time, stu.Arrival.Uncertainty, stu.Departure.Delay,
				stu.Departure.Time, stu.Departure.Uncertainty, stu.ScheduleRelationship,
			)
			if err != nil {
				return rollbackError(tx.Rollback(), err)
			}
		}
	}
	if err := tuStmt.Close(); err != nil {
		return rollbackError(tx.Rollback(), err)
	}
	if err := stuStmt.Close(); err != nil {
		return rollbackError(tx.Rollback(), err)
	}

	if err := tx.Commit(); err != nil {
		return rollbackError(tx.Rollback(), err)
	}
	return nil
}

// UpdateTripUpdateBytes reads bytes and updates the TripUpdates DB
func (tuds *TripUpdateSQLDataset) UpdateTripUpdateBytes(ctx context.Context, b []byte) error {
	var tu trimet.TripUpdatesMsg

	_, err := tu.UnmarshalMsg(b)
	if err != nil {
		log.Println(err)
		return err
	}

	if err = tuds.UpdateTripUpdates(tu.TripUpdates); err != nil {
		return err
	}
	return nil
}

// ProduceTripUpdates makes requests to the Trimet API and sends the results to
// a Producer.
func ProduceTripUpdates(ctx context.Context,baseURL string,  apiKey string, p Producer) error {
	ticker := time.NewTicker(tripUpdateDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}

		tripUpdates, err := trimet.RequestTripUpdate(baseURL, apiKey)
		if err != nil {
			log.Println(err)
			continue
		}
		tripUpdatesMsg := trimet.TripUpdatesMsg{
			TripUpdates: tripUpdates,
		}

		var b []byte
		msgBytes, err := tripUpdatesMsg.MarshalMsg(b)
		if err != nil {
			log.Println(err)
			continue
		}

		if err = p.Produce(msgBytes); err != nil {
			log.Println(err)
			continue
		}

	}
}
