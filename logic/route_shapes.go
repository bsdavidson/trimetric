package logic

import (
	"log"
	"sort"

	postgis "github.com/cridenour/go-postgis"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

// RouteShape represents a complete shape line for a given RouteID
type RouteShape struct {
	DirectionID int          `json:"direction_id"`
	RouteID     string       `json:"route_id"`
	Color       string       `json:"color"`
	Points      []RoutePoint `json:"points"`
}

// RoutePoint represents a single point along a route
type RoutePoint struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// TripShape adds a TripID to a RouteShape in order for it to be associated later
// in the client.
type TripShape struct {
	RouteShape
	TripID string `json:"trip_id"`
}

// FetchTripShapes takes a slice of trip ID's and returns a map that associates
// trip id's with all of their trip shapes. This is used client side to
// render a route line for a specific trip a vehicle is on.
func (sd *ShapeSQLDataset) FetchTripShapes(tripIDs []string) (map[string]*TripShape, error) {
	log.Println("Fetchiung Shapes for tripID:", tripIDs)
	q := `
		SELECT
			shapes.id,
			shapes.pt_lon_lat,
			trips.id AS trip_id,
			trips.route_id,
			trips.direction_id,
			routes.color AS route_color
		FROM trips
		JOIN shapes ON shapes.id = trips.shape_id
		JOIN routes ON routes.id = trips.route_id
		WHERE trips.id = ANY($1)
		ORDER BY trips.id, shapes.pt_sequence ASC
	`

	rows, err := sd.DB.Query(q, pq.Array(tripIDs))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	shapes := map[string]*TripShape{}
	var lastShapeID int
	var lastShape *TripShape
	for rows.Next() {
		var tripID string
		var routeID string
		var direction int
		var id int
		var color string
		var p RoutePoint
		var lngLat postgis.PointS
		err := rows.Scan(&id, &lngLat, &tripID, &routeID, &direction, &color)
		if err != nil {
			return nil, err
		}
		p.Lat = lngLat.Y
		p.Lng = lngLat.X

		if lastShape == nil || lastShapeID != id {

			lastShape = &TripShape{
				TripID: tripID,
				RouteShape: RouteShape{
					RouteID:     routeID,
					DirectionID: direction,
					Color:       color,
				},
			}

			shapes[tripID] = lastShape

			lastShapeID = id
		}
		lastShape.Points = append(lastShape.Points, p)
	}
	return shapes, nil
}

// FetchRouteShapes returns all shapes for train routes and flattens them to
// reduce the amount of data.
func (sd *ShapeSQLDataset) FetchRouteShapes() ([]*RouteShape, error) {
	q := `
		SELECT id, pt_lon_lat, route_shapes.route_id, route_shapes.direction_id, route_shapes.route_color
		FROM shapes
		JOIN (
			SELECT DISTINCT trips.shape_id, routes.id as route_id, trips.direction_id as direction_id, routes.color as route_color
			FROM routes
			JOIN trips ON trips.route_id = routes.id
			WHERE routes.type = 0
		) AS route_shapes ON route_shapes.shape_id = shapes.id
		ORDER BY id, pt_sequence ASC
	`

	rows, err := sd.DB.Query(q)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()

	shapes := []*RouteShape{}
	var lastShapeID int
	var lastShape *RouteShape
	for rows.Next() {
		var routeID string
		var direction int
		var id int
		var color string
		var p RoutePoint
		var lngLat postgis.PointS
		err := rows.Scan(&id, &lngLat, &routeID, &direction, &color)
		if err != nil {
			return nil, err
		}
		p.Lat = lngLat.Y
		p.Lng = lngLat.X

		if lastShape == nil || lastShapeID != id {
			shapes = append(shapes, &RouteShape{RouteID: routeID, DirectionID: direction, Color: color})
			lastShape = shapes[len(shapes)-1]
			lastShapeID = id
		}
		lastShape.Points = append(lastShape.Points, p)
	}

	return removeDuplicateSegments(shapes), nil
}

type segment struct {
	start RoutePoint
	end   RoutePoint
}

type segmentedShape struct {
	shape    *RouteShape
	segments []segment
}

func removeDuplicateSegments(shapes []*RouteShape) []*RouteShape {
	segmentedShapes := make([]segmentedShape, len(shapes))
	for i, shp := range shapes {
		sshp := &segmentedShapes[i]
		sshp.shape = shp
		sshp.segments = make([]segment, len(shp.Points)-1)
		for j := 0; j < len(shp.Points)-1; j++ {
			var p, np RoutePoint

			switch shp.DirectionID {
			case 0:
				p = shp.Points[j]
				np = shp.Points[j+1]
			case 1:
				p = shp.Points[len(shp.Points)-j-1]
				np = shp.Points[len(shp.Points)-j-2]
			}
			sshp.segments[j] = segment{start: p, end: np}
		}
	}

	sort.Slice(segmentedShapes, func(i, j int) bool {
		return len(segmentedShapes[i].segments) > len(segmentedShapes[j].segments)
	})

	uniqueSegments := map[segment]struct{}{}
	var shape *RouteShape
	var uniqueShapes []*RouteShape
	for _, sshp := range segmentedShapes {
		for _, seg := range sshp.segments {
			if _, ok := uniqueSegments[seg]; !ok {
				uniqueSegments[seg] = struct{}{}
				if shape == nil {
					shape = &RouteShape{
						Color:   sshp.shape.Color,
						RouteID: sshp.shape.RouteID,
						Points:  []RoutePoint{seg.start},
					}
				}
				shape.Points = append(shape.Points, seg.end)
				continue
			}
			if shape != nil {
				uniqueShapes = append(uniqueShapes, shape)
				shape = nil
			}
		}
		if shape != nil {
			uniqueShapes = append(uniqueShapes, shape)
			shape = nil
		}
	}
	return uniqueShapes
}
