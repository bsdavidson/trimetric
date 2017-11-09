package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"

	"github.com/bsdavidson/trimetric"
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
	migrate := flag.Bool("migrate", false, "Perform database migrations if true.")
	migratePath := flag.String("migrate-path", "./migrations", "Path to migration files")
	redisAddr := flag.String("redis", "redis:6379", "Redis address")
	flag.Parse()

	b, err := ioutil.ReadFile(apiKeyPath)
	if err != nil {
		log.Fatalf("Error reading %s: %v", apiKeyPath, err)
	}
	apiKey := strings.TrimSpace(string(b))

	db, err := trimetric.OpenDB(fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", *pgUser, *pgPassword, *pgHost, *pgDatabase))
	if err != nil {
		log.Fatal(err)
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

	if err := trimetric.Run(ctx, cancel, *addr, apiKey, db, *redisAddr, *webPath); err != nil {
		log.Fatal(err)
	}

}
