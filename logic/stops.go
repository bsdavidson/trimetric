package logic

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bsdavidson/trimetric/trimet"
	postgis "github.com/cridenour/go-postgis"
	"github.com/garyburd/redigo/redis"
)

func rollbackError(fn string, rberr error, err error) error {
	if rberr == nil {
		return err
	}
	return fmt.Errorf("%s: error rolling back: %s (for %s)", fn, rberr, err)
}

type StopWithDistance struct {
	trimet.Stop
	Distance float64 `json:"distance"`
}

type StopDataset interface {
	FetchWithinDistance(lat, lng, dist string) ([]StopWithDistance, error)
	UpsertAll(stops [][]string) error
}

type StopSQLDataset struct {
	DB *sql.DB
}

func (sd *StopSQLDataset) FetchWithinDistance(lat, lng, dist string) ([]StopWithDistance, error) {
	q := `SELECT id, code, name, "desc", lat_lon, zone_id, stop_url,
							 location_type, parent_station, direction, position,
							 ST_Distance(ST_GeogFromText($1), lat_lon) as distance
				FROM stops
				WHERE ST_DWithin(ST_GeogFromText($1), lat_lon, $2)
				ORDER BY distance ASC`
	rows, err := sd.DB.Query(q, fmt.Sprintf("SRID=4326;POINT(%s %s)", lng, lat), dist)
	if err != nil {
		return nil, err
	}
	stops := []StopWithDistance{}

	for rows.Next() {
		var s StopWithDistance
		var latLon postgis.PointS
		if err := rows.Scan(&s.ID, &s.Code, &s.Name, &s.Desc, &latLon, &s.ZoneID, &s.StopURL, &s.LocationType, &s.ParentStation, &s.Direction, &s.Position, &s.Distance); err != nil {
			return nil, err
		}
		s.Lat = latLon.Y
		s.Lon = latLon.X
		stops = append(stops, s)
	}
	return stops, nil
}

func (sd *StopSQLDataset) UpsertAll(stops [][]string) error {
	q := `INSERT INTO stops
					(id, code, name, "desc", lat_lon, zone_id, stop_url, location_type,
					parent_station, direction, position)
				VALUES ($1, $2, $3, $4, ST_GeogFromText($5), $6, $7, $8, $9, $10, $11)
				ON CONFLICT (id) DO UPDATE SET
					code = EXCLUDED.code,
					name = EXCLUDED.name,
					"desc" = EXCLUDED.desc,
					lat_lon = EXCLUDED.lat_lon,
					zone_id = EXCLUDED.zone_id,
					stop_url = EXCLUDED.stop_url,
					location_type = EXCLUDED.location_type,
					parent_station = EXCLUDED.parent_station,
					direction = EXCLUDED.direction,
					position = EXCLUDED.position
			 `
	tx, err := sd.DB.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(q)
	if err != nil {
		return rollbackError("StopSQLDataset.UpsertAll", tx.Rollback(), err)
	}

	for i, stop := range stops {
		if i == 0 {
			continue
		}

		locType, err := strconv.Atoi(stop[8])
		if err != nil {
			return rollbackError("StopSQLDataset.UpsertAll", tx.Rollback(), err)
		}

		latLon := fmt.Sprintf("SRID=4326;POINT(%s %s)", stop[5], stop[4])
		_, err = stmt.Exec(stop[0], stop[1], stop[2], stop[3], latLon, stop[6], stop[7], locType, stop[9], stop[10], stop[11])
		if err != nil {
			return rollbackError("StopSQLDataset.UpsertAll", tx.Rollback(), err)
		}
	}
	if err := stmt.Close(); err != nil {
		return rollbackError("StopSQLDataset.UpsertAll", tx.Rollback(), err)
	}

	if err := tx.Commit(); err != nil {
		return rollbackError("StopSQLDataset.UpsertAll", tx.Rollback(), err)
	}
	log.Printf("Wrote %d stops", len(stops))
	return nil
}

// PollStops ...
func PollStops(ctx context.Context, sd StopDataset, redisPool *redis.Pool, dur time.Duration) {
	rc := redisPool.Get()
	defer rc.Close()

	lastTs, err := redis.Int64(rc.Do("GET", "poll_stops_last_ts"))
	if err != nil {
		log.Println(err)
	}

	timer := time.NewTimer(time.Until(time.Unix(lastTs, 0).Add(dur)))

	for {
		select {
		case <-timer.C:
			timer.Reset(dur)
			_, err := rc.Do("SET", "poll_stops_last_ts", time.Now().Unix())
			if err != nil {
				log.Printf("PollStops: %v", err)
			}

			stops, err := trimet.RequestStops()
			if err != nil {
				log.Println("PollStops:", err)
				continue
			}
			log.Printf("Fetched %d stops", len(stops))

			if err := sd.UpsertAll(stops); err != nil {
				log.Println(err)
				continue
			}

			log.Println("########### Fetch and Write Stops ############")

		case <-ctx.Done():
			break
		}
	}
}
