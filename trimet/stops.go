package trimet

import (
	"archive/zip"
	"encoding/csv"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Stop struct {
	ID                 string  `json:"id"`
	Code               string  `json:"code"`
	Name               string  `json:"name"`
	Desc               string  `json:"desc"`
	Lat                float64 `json:"lat"`
	Lon                float64 `json:"lon"`
	ZoneID             string  `json:"zone_id"`
	StopURL            string  `json:"stop_url"`
	LocationType       int     `json:"location_type"`
	ParentStation      string  `json:"parent_station"`
	Direction          string  `json:"direction"`
	Position           string  `json:"position"`
	WheelchairBoarding int     `json:"wheelchair_boarding"`
}

// RequestStops ...
func RequestStops() ([][]string, error) {
	log.Println("Fetching Stops")
	f, err := ioutil.TempFile("", "tmp")
	if err != nil {
		return nil, err
	}
	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	resp, err := http.Get("https://developer.trimet.org/schedule/gtfs.zip")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return nil, err
	}
	f.Close()

	z, err := zip.OpenReader(f.Name())
	if err != nil {
		return nil, err
	}

	var stopIdx int
	for i, zf := range z.File {
		if zf.Name == "stops.txt" {
			stopIdx = i
			break
		}
	}

	rc, err := z.File[stopIdx].Open()
	if err != nil {
		return nil, err
	}
	stops, err := csv.NewReader(rc).ReadAll()
	if err != nil {
		return nil, err
	}
	log.Println("Done getting stops")
	return stops, nil
}
