package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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

// HandleVehiclePositions provides responses for the /api/v2/vehicles endpoint.
// It returns a list of vehicles pulled from a local DB populated from a GTFS feed.
func HandleVehiclePositions(vd logic.VehicleDataset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		ids, err := commaSplitInts(r.URL.Query().Get("ids"))
		if err != nil {
			http.Error(w, fmt.Sprintf("error parsing ids: %v", err), http.StatusBadRequest)
			return
		}

		vehicles, err := vd.FetchVehiclePositionsByIDs(ids)
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

// HandleArrivals returns a list of upcoming arrivals for a list of stops.
func HandleArrivals(tds logic.StopDataset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ids := strings.Split(r.URL.Query().Get("locIDs"), ",")
		if len(ids) == 0 && ids[0] == "" {
			http.Error(w, "error: a list of stop IDs are required", http.StatusBadRequest)
			return
		}
		ar, err := tds.FetchArrivals(ids)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(ar)
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

// HandleWSVehicles ...
func HandleWSVehicles(vd logic.VehicleDataset, updateChan <-chan trimet.VehiclePosition) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("websocket open")
		var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print(err)
			return
		}
		defer func() {
			c.Close()
			log.Println("websocket closed")
		}()

		for {
			_, _, err := c.ReadMessage() // just reading in case we need to kill conn
			if err != nil {
				log.Println(err)
				break
			}
			select {
			case v := <-updateChan:
				if err := c.WriteJSON(v); err != nil {
					log.Println(err)
					break
				}
			}
		}

	}
}
