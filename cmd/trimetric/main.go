package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/bsdavidson/trimetric/api"
	"github.com/bsdavidson/trimetric/logic"
	"github.com/garyburd/redigo/redis"
	_ "github.com/lib/pq"
)

const apiKeyPath = "/run/secrets/trimet-api-key"

func main() {
	addr := flag.String("addr", ":80", "Address to bind to")
	webPath := flag.String("web-path", "./web/dist", "Path to website assets")
	pgUser := flag.String("pg-user", "postgres", "Postgres username")
	pgPassword := flag.String("pg-password", "example", "Postgres password")
	pgHost := flag.String("pg-host", "postgres", "Postgres hostname")
	pgDatabase := flag.String("pg-database", "trimetric", "Postgres database")
	redisAddr := flag.String("redis", "redis:6379", "Redis address")
	flag.Parse()

	b, err := ioutil.ReadFile(apiKeyPath)
	if err != nil {
		log.Fatalf("Error reading %s: %v", apiKeyPath, err)
	}
	apiKey := strings.TrimSpace(string(b))

	dbURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", *pgUser, *pgPassword, *pgHost, *pgDatabase)
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	vds := &logic.VehicleSQLDataset{DB: db}
	sds := &logic.StopSQLDataset{DB: db}

	redisPool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", *redisAddr) },
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		<-signals
		cancel()
	}()

	wg := &sync.WaitGroup{}
	wg.Add(3)

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

	srv := &http.Server{Addr: *addr, Handler: http.DefaultServeMux}
	http.HandleFunc("/api/v1/vehicles", api.HandleVehicles(vds))
	http.HandleFunc("/api/v1/arrivals", api.HandleArrivals(apiKey))
	http.HandleFunc("/api/v1/stops", api.HandleStops(sds))
	// http.HandleFunc("/api/v1/gtfs", api.HandleGTFS())
	http.Handle("/", http.FileServer(http.Dir(*webPath)))
	log.Printf("Serving requests on %s", *addr)

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Println(err)
		}
	}()

	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
	wg.Wait()
}
