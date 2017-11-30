package trimetric

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	_ "net/http/pprof" // Import the pprof package to add the HTTP handlers
	"strings"
	"sync"
	"time"

	"github.com/bsdavidson/trimetric/api"
	"github.com/bsdavidson/trimetric/logic"
	"github.com/bsdavidson/trimetric/trimet"
	"github.com/garyburd/redigo/redis"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/pkg/errors"
	"github.com/pressly/goose"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// OpenDB connects to the database.
func OpenDB(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err := db.Ping(); err != nil {
		return nil, errors.WithStack(err)
	}
	return db, nil
}

// OpenInfluxDB will open a connection and return an InfluxDB client
func OpenInfluxDB(host, username, password string) (client.Client, error) {
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     host,
		Username: username,
		Password: password,
		Timeout:  100 * time.Millisecond,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return c, nil
}

// MigrateDB applies any updates needed to the database to bring it up to the
// latest version of the schema. If migrate is false, then it won't modify the DB
// but will instead return an error if the database is not up-to-date.
func MigrateDB(db *sql.DB, path string, migrate bool) error {
	dbVer, err := goose.EnsureDBVersion(db)
	if err != nil {
		return errors.WithStack(err)
	}
	migrations, err := goose.CollectMigrations(path, 0, goose.MaxVersion)
	if err != nil {
		return errors.WithStack(err)
	}
	lm, err := migrations.Last()
	if err != nil {
		return errors.WithStack(err)
	}
	if !migrate {
		if lm.Version > dbVer {
			return errors.Errorf("database version update required, latest version: %d, current version:%d", lm.Version, dbVer)
		}
		return nil
	}

	if err := goose.Up(db, path); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Run starts all the processes for Trimetric.
func Run(ctx context.Context, cancel context.CancelFunc, debug bool, addr, apiKey string, db *sql.DB, influxClient client.Client, kafkaAddr, redisAddr, webPath string) error {
	vds := &logic.VehicleSQLDataset{DB: db}
	sds := &logic.StopSQLDataset{DB: db}
	shds := &logic.ShapeSQLDataset{DB: db}
	rds := &logic.RouteSQLDataset{DB: db}
	lds := &logic.LoaderSQLDataset{DB: db}
	tuds := &logic.TripUpdateSQLDataset{DB: db}

	redisPool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", redisAddr) },
	}

	wg := &sync.WaitGroup{}
	wg.Add(4)
	defer func() {
		cancel()
		wg.Wait()
	}()

	go func() {
		defer cancel()
		defer wg.Done()
		producer, err := logic.NewKafkaProducer(logic.VehiclePositionsTopic, []string{kafkaAddr})
		if err != nil {
			log.Println(err)
			return
		}
		defer func() {
			if err := producer.Close(); err != nil {
				log.Println(err)
			}
		}()

		if err := logic.ProduceVehiclePositions(ctx, producer, trimet.BaseTrimetURL, strings.TrimSpace(apiKey), time.Second); err != nil {
			log.Println(err)
		}
	}()

	go func() {
		defer cancel()
		defer wg.Done()
		producer, err := logic.NewKafkaProducer(logic.TripUpdatesTopic, []string{kafkaAddr})
		if err != nil {
			log.Println(err)
			return
		}
		defer func() {
			if err := producer.Close(); err != nil {
				log.Println(err)
			}
		}()

		if err := logic.ProduceTripUpdates(ctx, trimet.BaseTrimetURL, strings.TrimSpace(apiKey), producer); err != nil {
			log.Println(err)
		}
	}()

	go func() {
		defer cancel()
		defer wg.Done()
		if err := logic.ConsumeKafkaTopic(ctx, vds.UpsertVehiclePositionBytes, logic.VehiclePositionsTopic, []string{kafkaAddr}); err != nil {
			log.Println(err)
		}
	}()

	go func() {
		defer cancel()
		defer wg.Done()
		if err := logic.ConsumeKafkaTopic(ctx, tuds.UpdateTripUpdateBytes, logic.TripUpdatesTopic, []string{kafkaAddr}); err != nil {
			log.Println(err)
		}
	}()

	go func() {
		defer cancel()
		defer wg.Done()
		logic.PollGTFSData(ctx, lds, redisPool, 24*time.Hour)
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Println("debug listening on :9876")
	internalSrv := &http.Server{Addr: ":9876", Handler: http.DefaultServeMux}
	go func() {
		if err := internalSrv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/vehicles", api.HandleVehiclePositions(vds))
	mux.HandleFunc("/api/v1/trimet/arrivals", api.HandleTrimetArrivals(apiKey))
	mux.HandleFunc("/api/v1/arrivals", api.HandleArrivals(sds))
	mux.HandleFunc("/api/v1/routes", api.HandleRoutes(rds))
	mux.HandleFunc("/api/v1/routes/lines", api.HandleRouteLines(shds))
	mux.HandleFunc("/api/v1/shapes", api.HandleShapes(shds))
	mux.HandleFunc("/api/v1/stops", api.HandleStops(sds))
	mux.HandleFunc("/api/v1/trip", api.HandleTripUpdates(tuds))
	mux.HandleFunc("/api/ws", api.HandleWSVehicles(vds, shds, sds, rds))
	mux.Handle("/", http.FileServer(http.Dir(webPath)))
	srv := &http.Server{Addr: addr, Handler: mux}
	log.Printf("serving requests on %s", addr)

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Println(err)
		}
		if err := internalSrv.Shutdown(shutdownCtx); err != nil {
			log.Println(err)
		}
	}()

	if err := srv.ListenAndServe(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
