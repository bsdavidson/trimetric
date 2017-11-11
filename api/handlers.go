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
// It searches for stops within a specified distance from a lat/lng.
func HandleStops(sd logic.StopDataset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lat := r.URL.Query().Get("lat")
		lng := r.URL.Query().Get("lng")
		dist := r.URL.Query().Get("distance")

		if lat == "" {
			httpError(w, "HandleStops:", fmt.Errorf("Latitude cannot be blank"), http.StatusBadRequest)
			return
		}
		if lng == "" {
			httpError(w, "HandleStops:", fmt.Errorf("Longitude cannot be blank"), http.StatusBadRequest)
			return
		}
		if dist == "" {
			httpError(w, "HandleStops:", fmt.Errorf("Distance cannot be blank"), http.StatusBadRequest)
			return
		}

		stops, err := sd.FetchWithinDistance(lat, lng, dist)
		if err != nil {
			httpError(w, "HandleStops:", err, http.StatusInternalServerError)
			return
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

// HandleArrivals provides responses for the /api/v1/arrivals endpoint.
// It proxies requests mostly untouched to the trimet API and returns a list of
// arrivals for the specified location IDs.
func HandleArrivals(apiKey string) http.HandlerFunc {
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

// HandleVehicles provides responses for the /api/v1/vehicles endpoint.
// It returns a list of vehicles pulled from a local DB, optionally filtered
// by a list of IDs.
func HandleVehicles(vd logic.VehicleDataset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ids, err := commaSplitInts(r.URL.Query().Get("ids"))
		if err != nil {
			http.Error(w, fmt.Sprintf("error parsing ids: %v", err), http.StatusBadRequest)
			return
		}

		vehicles, err := vd.FetchByIDs(ids)
		if err != nil {
			httpError(w, "HandleVehicles:", err, http.StatusInternalServerError)
			return
		}

		var rv []trimet.RawVehicleData
		for _, v := range vehicles {
			rv = append(rv, v.Data)
		}

		b, err := json.Marshal(rv)
		if err != nil {
			httpError(w, "HandleVehicles:", err, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(b); err != nil {
			httpError(w, "HandleVehicles:", err, http.StatusInternalServerError)
			return
		}
	}
}

// HandleGTFSVehicles provides responses for the /api/v2/vehicles endpoint.
// It returns a list of vehicles pulled from a local DB populated from a GTFS feed.
func HandleGTFSVehicles(vd logic.VehicleDataset) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		ids, err := commaSplitInts(r.URL.Query().Get("ids"))
		if err != nil {
			http.Error(w, fmt.Sprintf("error parsing ids: %v", err), http.StatusBadRequest)
			return
		}

		vehicles, err := vd.FetchGTFSByIDs(ids)
		if err != nil {
			httpError(w, "HandleGTFSVehicles:", err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		b, err := json.Marshal(vehicles)
		if err != nil {
			httpError(w, "HandleGTFSVehicles:", err, http.StatusInternalServerError)
			return
		}
		w.Write(b)

	}
}
