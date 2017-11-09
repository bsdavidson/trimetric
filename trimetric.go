package trimetric

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/bsdavidson/trimetric/api"
	"github.com/bsdavidson/trimetric/logic"
	"github.com/garyburd/redigo/redis"
	"github.com/pressly/goose"
)

// OpenDB ...
func OpenDB(url string) (*sql.DB, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("trimetric.OpenDB: %v", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("trimetric.OpenDB: %v", err)
	}
	return db, nil
}

// MigrateDB ...
func MigrateDB(db *sql.DB, path string, migrate bool) error {
	dbVer, err := goose.EnsureDBVersion(db)
	if err != nil {
		return fmt.Errorf("trimetric.MigrateDB: %v", err)
	}
	migrations, err := goose.CollectMigrations(path, 0, goose.MaxVersion)
	if err != nil {
		return fmt.Errorf("trimetric.MigrateDB: %v", err)
	}
	lm, err := migrations.Last()
	if err != nil {
		return fmt.Errorf("trimetric.MigrateDB: %v", err)
	}
	if !migrate {
		if lm.Version > dbVer {
			return fmt.Errorf("trimetric.MigrateDB: database version update required, latest version: %d, current version:%d", lm.Version, dbVer)
		}
		return nil
	}

	if err := goose.Up(db, path); err != nil {
		return fmt.Errorf("trimetric.MigrateDB: %v", err)
	}
	return nil
}

// Run ...
func Run(ctx context.Context, cancel context.CancelFunc, addr, apiKey string, db *sql.DB, redisAddr, webPath string) error {

	vds := &logic.VehicleSQLDataset{DB: db}
	sds := &logic.StopSQLDataset{DB: db}

	redisPool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", redisAddr) },
	}

	wg := &sync.WaitGroup{}
	wg.Add(3)
	defer func() {
		cancel()
		wg.Wait()
	}()

	go func() {
		defer wg.Done()
		if err := logic.ProduceVehicles(ctx, strings.TrimSpace(apiKey)); err != nil {
			log.Println(err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := logic.ConsumeVehicles(ctx, vds); err != nil {
			log.Println(err)
		}
	}()

	go func() {
		defer wg.Done()
		logic.PollStops(ctx, sds, redisPool, 24*time.Hour)
	}()

	srv := &http.Server{Addr: addr, Handler: http.DefaultServeMux}
	http.HandleFunc("/api/v1/vehicles", api.HandleVehicles(vds))
	http.HandleFunc("/api/v1/arrivals", api.HandleArrivals(apiKey))
	http.HandleFunc("/api/v1/stops", api.HandleStops(sds))
	http.Handle("/", http.FileServer(http.Dir(webPath)))
	log.Printf("Serving requests on %s", addr)

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Println(err)
		}
	}()

	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("trimetric.Run: %v", err)
	}

	return nil
}
