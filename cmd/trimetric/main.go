package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func handleVehicles(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vehicles := []Vehicle{}
		q := `SELECT id, vehicle_id, type, sign_message
					FROM vehicles ORDER BY vehicle_id`
		rows, err := db.Query(q)
		if err != nil {
			log.Println("db.Query:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		for rows.Next() {
			var v Vehicle
			if err := rows.Scan(&v.ID, &v.VehicleID, &v.Type, &v.SignMessage); err != nil {
				log.Println("rows.Scan:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			vehicles = append(vehicles, v)
		}
		if err := rows.Err(); err != nil {
			log.Println("rows.Err:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(vehicles)
		if err != nil {
			log.Println("json.Marshal:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := w.Write(b); err != nil {
			log.Println("w.Write:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// TrimetVehicleResponse ...
type TrimetVehicleResponse struct {
	ResultSet TrimetVehicleResultSet `json:"resultSet"`
}

// TrimetVehicleResultSet ...
type TrimetVehicleResultSet struct {
	QueryTime int             `json:"queryTime"`
	Vehicles  []TrimetVehicle `json:"vehicle"`
}

// TrimetVehicle ...
type TrimetVehicle struct {
	VehicleID   int    `json:"vehicleID"`
	Type        string `json:"type"`
	SignMessage string `json:"signMessage"`
}

// Vehicle ...
type Vehicle struct {
	ID          int    `json:"id"`
	VehicleID   int    `json:"vehicle_id"`
	Type        string `json:"type"`
	SignMessage string `json:"sign_message"`
}

// https://developer.trimet.org/ws/v2/vehicles?appID=65795DCAB40706D335474B716&json=true

func requestVehicles(apiKey string) (*TrimetVehicleResponse, error) {
	query := url.Values{}
	query.Set("appID", apiKey)
	query.Set("json", "true")
	resp, err := http.Get("https://developer.trimet.org/ws/v2/vehicles?" + query.Encode())
	if err != nil {
		return nil, fmt.Errorf("http.Get: %s", err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ioutil.ReadAll: %s", err)
	}
	var tr TrimetVehicleResponse
	err = json.Unmarshal(b, &tr)
	if err != nil {
		return nil, fmt.Errorf("json.Unmarshal: %s", err)
	}
	return &tr, nil
}

func pollVehicles(db *sql.DB, apiKey string) {
	for {
		time.Sleep(time.Second)
		tr, err := requestVehicles(apiKey)
		if err != nil {
			log.Println("requestVehicles:", err)
			continue
		}
		for _, tv := range tr.ResultSet.Vehicles {
			v := Vehicle{
				VehicleID:   tv.VehicleID,
				Type:        tv.Type,
				SignMessage: tv.SignMessage,
			}

			q := `INSERT INTO vehicles (vehicle_id, type, sign_message)
						VALUES ($1, $2, $3)
						ON CONFLICT (vehicle_id) DO UPDATE SET
							type = EXCLUDED.type,
							sign_message = EXCLUDED.sign_message;
					 `
			_, err := db.Exec(q, v.VehicleID, v.Type, v.SignMessage)
			if err != nil {
				log.Println("db.Exec:", err)
			}
		}
	}
}

func main() {
	addr := flag.String("addr", ":80", "Address to bind to")
	webPath := flag.String("web-path", "./web/dist", "Path to website assets")
	pgUser := flag.String("pg-user", "postgres", "Postgres username")
	pgPassword := flag.String("pg-password", "example", "Postgres password")
	pgHost := flag.String("pg-host", "postgres", "Postgres hostname")
	pgDatabase := flag.String("pg-database", "trimetric", "Postgres database")
	flag.Parse()

	apiKey := os.Getenv("TRIMET_API_KEY")
	if apiKey == "" {
		log.Fatal("TRIMET_API_KEY must be set")
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", *pgUser, *pgPassword, *pgHost, *pgDatabase)
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	go pollVehicles(db, apiKey)

	http.HandleFunc("/api/v1/vehicles", handleVehicles(db))
	http.Handle("/", http.FileServer(http.Dir(*webPath)))
	log.Printf("Serving requests on %s", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
