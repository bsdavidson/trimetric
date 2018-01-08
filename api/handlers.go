package api

import (
	"compress/zlib"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/bsdavidson/trimetric/logic"
	"github.com/bsdavidson/trimetric/trimet"
	"github.com/gorilla/websocket"
)

func httpError(w http.ResponseWriter, prefix string, err error, code int) {
	log.Println(prefix, err)
	http.Error(w, http.StatusText(code), code)
}

func commaSplit(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}

func commaSplitInts(s string) ([]int, error) {
	var nums []int
	for _, sn := range commaSplit(s) {
		n, err := strconv.Atoi(sn)
		if err != nil {
			return nil, err
		}
		nums = append(nums, n)
	}
	return nums, nil
}

type stopsWithDistanceResponse struct {
	Stops []logic.StopWithDistance `json:"stops"`
}

// HandleStops provides responses for the /api/v1/stops endpoint.
// It searches for stops within a specified distance from a lat/lon.
func HandleStops(sd logic.StopDataset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lat := r.URL.Query().Get("lat")
		lng := r.URL.Query().Get("lng")
		dist := r.URL.Query().Get("distance")
		south := r.URL.Query().Get("south")
		north := r.URL.Query().Get("north")
		east := r.URL.Query().Get("east")
		west := r.URL.Query().Get("west")

		var stops []logic.StopWithDistance
		var err error
		if south != "" && north != "" && east != "" && west != "" {
			stops, err = sd.FetchWithinBox(west, south, east, north)
			if err != nil {
				httpError(w, "HandleStops:", err, http.StatusInternalServerError)
				return
			}
		} else if lat != "" && lng != "" && dist != "" {
			stops, err = sd.FetchWithinDistance(lat, lng, dist)
			if err != nil {
				httpError(w, "HandleStops:", err, http.StatusInternalServerError)
				return
			}
		} else {
			stops, err = sd.FetchAllStops()
			if err != nil {
				httpError(w, "HandleStops:", err, http.StatusInternalServerError)
				return
			}
		}

		b, err := json.Marshal(stopsWithDistanceResponse{Stops: stops})
		if err != nil {
			httpError(w, "HandleStops:", err, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(b); err != nil {
			httpError(w, "HandleStops:", err, http.StatusInternalServerError)
			return
		}
	}
}

// HandleTrimetArrivals provides responses for the /api/v1/arrivals endpoint.
// It proxies requests mostly untouched to the trimet API and returns a list of
// arrivals for the specified location IDs.
func HandleTrimetArrivals(apiKey string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ids, err := commaSplitInts(r.URL.Query().Get("locIDs"))
		if err != nil {
			http.Error(w, fmt.Sprintf("error parsing ids: %v", err), http.StatusBadRequest)
			return
		}

		b, err := trimet.RequestArrivals(apiKey, ids)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

	}
}

// HandleVehiclePositions provides responses for the /api/v1/vehicles endpoint.
// It returns a list of vehicles pulled from a local DB populated from a GTFS feed.
func HandleVehiclePositions(vd logic.VehicleDataset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		sinceStr := r.URL.Query().Get("since")
		var since int
		if sinceStr != "" {
			since, err = strconv.Atoi(sinceStr)
			if err != nil {
				http.Error(w, fmt.Sprintf("error parsing since: %v", err), http.StatusBadRequest)
				return
			}
		}

		vehicles, err := vd.FetchVehiclePositions(since)
		if err != nil {
			httpError(w, "HandleVehiclePositions:", err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		b, err := json.Marshal(vehicles)
		if err != nil {
			httpError(w, "HandleVehiclePositions:", err, http.StatusInternalServerError)
			return
		}
		w.Write(b)

	}
}

// HandleTripUpdates returns a list of trip updates as json
func HandleTripUpdates(tds logic.TripUpdatesDataset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tu, err := tds.FetchTripUpdates()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(tu)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

	}
}

func parseArgs(argStr, sep string) []string {
	argsSplit := strings.Split(argStr, sep)
	var argsArr []string
	if len(argsSplit) == 1 && argsSplit[0] == "" {
		return nil
	}

	for _, a := range argsSplit {
		if a != "" {
			argsArr = append(argsArr, a)
		}
	}
	if len(argsArr) == 0 {
		return nil
	}
	return argsArr
}

// ArrivalWithTrip ...
type ArrivalWithTrip struct {
	logic.Arrival
	TripShape *logic.TripShape `json:"trip_shape"`
}

// HandleArrivals returns a list of upcoming arrivals for a list of stops.
func HandleArrivals(sds logic.StopDataset, shds logic.ShapeDataset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ids := parseArgs(r.URL.Query().Get("stop_ids"), ",")
		if ids == nil {
			http.Error(w, "error: a list of stop IDs are required", http.StatusBadRequest)
			return
		}
		ar, err := sds.FetchArrivals(ids)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		awts := []ArrivalWithTrip{}

		trips := []string{}
		for _, a := range ar {
			trips = append(trips, a.TripID)
		}
		shapes, err := shds.FetchTripShapes(trips)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, a := range ar {
			awt := ArrivalWithTrip{
				Arrival:   a,
				TripShape: shapes[a.TripID],
			}
			awts = append(awts, awt)
		}

		b, err := json.Marshal(awts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(ar) == 0 {
			log.Println("arrivals empty")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

// HandleRoutes returns a list of routes
func HandleRoutes(rds logic.RouteDataset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		routes, err := rds.FetchRoutes()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(routes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

	}
}

// HandleShapes returns a list of shapes
func HandleShapes(sds logic.ShapeDataset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		shapeIDs := parseArgs(r.URL.Query().Get("shape_ids"), ",")
		routeIDs := parseArgs(r.URL.Query().Get("route_ids"), ",")
		if shapeIDs == nil && routeIDs == nil {
			http.Error(w, "must specify route_id or shape_id", http.StatusBadRequest)
			return
		}

		shapes, err := sds.FetchShapes(routeIDs, shapeIDs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(shapes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

	}
}

// HandleRouteLines returns a list of routelines
func HandleRouteLines(sds logic.ShapeDataset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		lines, err := sds.FetchRouteShapes()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(lines)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)

	}
}

type wsMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type writeFunc func(c *websocket.Conn, message interface{}) error

// writeJSON writes data to the websocket conn
func writeJSON(c *websocket.Conn, v interface{}) error {
	c.EnableWriteCompression(true)
	w, err := c.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	err1 := json.NewEncoder(w).Encode(v)
	err2 := w.Close()
	if err1 != nil {
		return err1
	}
	return err2
}

// writeCompressedJSON compresses the data using zlib before sending it to the Conn
func writeCompressedJSON(c *websocket.Conn, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	w, err := c.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return err
	}
	defer w.Close()
	zw := zlib.NewWriter(w)
	defer zw.Close()
	_, err = zw.Write(b)
	if err != nil {
		return err
	}

	return nil
}

var count int64

func init() {
	var lastCount int64
	go func() {
		for {
			if count != lastCount {
				atomic.StoreInt64(&lastCount, count)
			}
			time.Sleep(2 * time.Second)
		}
	}()
}

func sendStaticDataToWS(c *websocket.Conn, chunkify bool, write writeFunc, shds logic.ShapeDataset, sds logic.StopDataset, rds logic.RouteDataset) error {
	var errored int64
	wg := &sync.WaitGroup{}
	wg.Add(3)

	var stops []logic.StopWithDistance
	go func() {
		defer wg.Done()
		var err error
		stops, err = sds.FetchAllStops()
		if err != nil {
			atomic.AddInt64(&errored, 1)
			log.Println(err)
		}
	}()

	var shapes []*logic.RouteShape
	go func() {
		defer wg.Done()
		var err error
		shapes, err = shds.FetchRouteShapes()
		if err != nil {
			atomic.AddInt64(&errored, 1)
			log.Println(err)
		}
	}()

	var routes []trimet.Route
	go func() {
		defer wg.Done()
		var err error
		routes, err = rds.FetchRoutes()
		if err != nil {
			atomic.AddInt64(&errored, 1)
			log.Println(err)
		}
	}()

	wg.Wait()
	if errored > 0 {
		return fmt.Errorf("error fetching static data")
	}
	err := write(c, wsMessage{
		Type: "totals",
		Data: map[string]int{
			"stops":        len(stops),
			"routes":       len(routes),
			"route_shapes": len(shapes),
		},
	})
	if err != nil {
		return err
	}
	if chunkify {
		var stopsPacket [100]logic.StopWithDistance
		for i, s := range stops {
			if i%100 == 0 || i == len(stops)-1 {
				if err := write(c, wsMessage{Type: "stops", Data: stopsPacket}); err != nil {
					return err
				}
				time.Sleep(25 * time.Millisecond)
			} else {
				stopsPacket[i%100] = s
			}
		}
	} else {
		if err := write(c, wsMessage{Type: "stops", Data: stops}); err != nil {
			return err
		}
	}

	if err := write(c, wsMessage{Type: "routes", Data: routes}); err != nil {
		return err
	}
	time.Sleep(25 * time.Millisecond)

	if err := write(c, wsMessage{Type: "route_shapes", Data: shapes}); err != nil {
		return err
	}
	return nil
}

func parseBool(s string, defaultValue bool) (bool, error) {
	if s == "" {
		return defaultValue, nil
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		return defaultValue, err
	}
	return b, nil
}

func parseInt(s string, defaultValue int64) (int64, error) {
	if s == "" {
		return defaultValue, nil
	}
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return defaultValue, err
	}
	return n, nil

}

// HandleWS maintains a websocket connection and pushes content to
// connected web clients.
func HandleWS(vds logic.VehicleDataset, shds logic.ShapeDataset, sds logic.StopDataset, rds logic.RouteDataset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Client Connecting")
		var err error

		version, err := parseInt(r.URL.Query().Get("version"), 0)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		chunkify, err := parseBool(r.URL.Query().Get("chunkify"), false)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		sendStaticData, err := parseBool(r.URL.Query().Get("static"), true)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		compressData, err := parseBool(r.URL.Query().Get("compress"), true)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var write writeFunc
		if compressData {
			write = writeCompressedJSON
		} else {
			write = writeJSON
		}

		var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print(err)
			return
		}
		defer c.Close()
		log.Printf("Client v%d Connected", version)

		atomic.AddInt64(&count, 1)
		defer atomic.AddInt64(&count, -1)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go func() {
			for {
				_, _, err := c.ReadMessage() // just reading in case we need to kill conn
				if err != nil {
					cancel()
					log.Println(err)
					break
				}
			}
		}()

		var since uint64
		var vehicles []logic.VehiclePositionWithRouteType
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			var err error
			vehicles, err = vds.FetchVehiclePositions(0)
			if err != nil {
				log.Println(err)
				cancel()
			}
			for _, v := range vehicles {
				if v.Timestamp > since {
					since = v.Timestamp
				}
			}
			wg.Done()
		}()

		if sendStaticData {
			// Will block until static data is sent or it errors
			if err = sendStaticDataToWS(c, chunkify, write, shds, sds, rds); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		// Send initial vehicles
		wg.Wait()
		if err := write(c, wsMessage{Type: "vehicles", Data: vehicles}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ticker := time.NewTicker(5000 * time.Millisecond)
		defer ticker.Stop()
	LOOP:
		for {
			select {
			case <-ctx.Done():
				break LOOP

			case <-ticker.C:
			}
			// 0 sends all positions on each push rather than just what has changed.
			// This is a temp debug setting that should be fixed.
			vehicles, err := vds.FetchVehiclePositions(0) //int(since)
			if err != nil {
				log.Println(err)
				continue
			}
			for _, v := range vehicles {
				if v.Timestamp > since {
					since = v.Timestamp
				}
			}

			if len(vehicles) == 0 {
				continue
			}
			if err := write(c, wsMessage{Type: "vehicles", Data: vehicles}); err != nil {
				log.Println(err)
				continue
			}
		}
	}
}
