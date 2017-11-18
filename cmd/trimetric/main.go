package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/bsdavidson/trimetric"
	_ "github.com/lib/pq"
)

const apiKeyPath = "/run/secrets/trimet-api-key"

func main() {
	addr := flag.String("addr", ":80", "Address to bind to")
	webPath := flag.String("web-path", "./web/dist", "Path to website assets")
	pgUser := flag.String("pg-user", "trimetric", "Postgres username")
	pgPassword := flag.String("pg-password", "example", "Postgres password")
	pgHost := flag.String("pg-host", "postgres", "Postgres hostname")
	pgDatabase := flag.String("pg-database", "trimetric", "Postgres database")
	migrate := flag.Bool("migrate", false, "Perform database migrations if true.")
	migratePath := flag.String("migrate-path", "./migrations", "Path to migration files")
	redisAddr := flag.String("redis-addr", "redis:6379", "Redis address")
	influxURL := flag.String("influx-url", "http://influxdb:8086", "InfluxDB URL")
	influxUser := flag.String("influx-user", "trimetric", "InfluxDB username")
	influxPassword := flag.String("influx-password", "example", "InfluxDB password")
	kafkaAddr := flag.String("kafka-addr", "kafka:9092", "Kafka broker address")

	flag.Parse()

	log.SetFlags(log.Lshortfile)

	b, err := ioutil.ReadFile(apiKeyPath)
	if err != nil {
		log.Fatalf("error reading %s: %v", apiKeyPath, err)
	}
	apiKey := strings.TrimSpace(string(b))

	var db *sql.DB
	for {
		db, err = trimetric.OpenDB(fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", *pgUser, *pgPassword, *pgHost, *pgDatabase))
		if err == nil {
			break
		}
		log.Println("error connecting to postgres:", err)
		log.Println("retying in 1 second")
		time.Sleep(time.Second)
	}

	influxClient, err := trimetric.OpenInfluxDB(*influxURL, *influxUser, *influxPassword)
	if err != nil {
		log.Println("InfluxDB:", err)
	}

	if err := trimetric.MigrateDB(db, *migratePath, *migrate); err != nil {
		log.Fatal(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		<-signals
		cancel()
	}()

	if err := trimetric.Run(ctx, cancel, *addr, apiKey, db, influxClient, *kafkaAddr, *redisAddr, *webPath); err != nil {
		log.Fatal(err)
	}

}
