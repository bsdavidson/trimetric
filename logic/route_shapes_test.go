package logic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveDuplicateSegments(t *testing.T) {

	routeLines := []*RouteShape{
		&RouteShape{
			RouteID: "100",
			Points: []RoutePoint{
				RoutePoint{Lat: 45, Lng: -122},
				RoutePoint{Lat: 46, Lng: -123},
				RoutePoint{Lat: 47, Lng: -124},
				RoutePoint{Lat: 48, Lng: -125},
			},
		},
		&RouteShape{
			RouteID:     "100",
			DirectionID: 1,
			Points: []RoutePoint{
				RoutePoint{Lat: 50, Lng: -127},
				RoutePoint{Lat: 49, Lng: -126},
				RoutePoint{Lat: 48, Lng: -125},
				RoutePoint{Lat: 47, Lng: -124},
				RoutePoint{Lat: 46, Lng: -123},
			},
		},
		&RouteShape{
			RouteID: "200",
			Points: []RoutePoint{
				RoutePoint{Lat: 48, Lng: -125},
				RoutePoint{Lat: 49, Lng: -126},
				RoutePoint{Lat: 50, Lng: -127},
				RoutePoint{Lat: 51, Lng: -128},
			},
		},

		&RouteShape{
			RouteID: "300",
			Points: []RoutePoint{
				RoutePoint{Lat: 46, Lng: -123},
				RoutePoint{Lat: 47, Lng: -124},
				RoutePoint{Lat: 48, Lng: -135},
				RoutePoint{Lat: 49, Lng: -126},
				RoutePoint{Lat: 50, Lng: -127},
			},
		},
	}

	expected := []*RouteShape{
		&RouteShape{
			RouteID: "100",
			Points: []RoutePoint{
				RoutePoint{
					Lat: 46,
					Lng: -123,
				},
				RoutePoint{
					Lat: 47,
					Lng: -124,
				},
				RoutePoint{
					Lat: 48,
					Lng: -125,
				},
				RoutePoint{
					Lat: 49,
					Lng: -126,
				},
				RoutePoint{
					Lat: 50,
					Lng: -127,
				},
			},
		},
		&RouteShape{
			RouteID: "300",
			Points: []RoutePoint{
				RoutePoint{
					Lat: 47,
					Lng: -124,
				},
				RoutePoint{
					Lat: 48,
					Lng: -135,
				},
				RoutePoint{
					Lat: 49,
					Lng: -126,
				},
			},
		},
		&RouteShape{
			RouteID: "100",
			Points: []RoutePoint{
				RoutePoint{
					Lat: 45,
					Lng: -122,
				},
				RoutePoint{
					Lat: 46,
					Lng: -123,
				},
			},
		},
		&RouteShape{
			RouteID: "200",
			Points: []RoutePoint{
				RoutePoint{
					Lat: 50,
					Lng: -127,
				},
				RoutePoint{
					Lat: 51,
					Lng: -128,
				},
			},
		},
	}

	newRouteLines := removeDuplicateSegments(routeLines)

	assert.Equal(t, expected, newRouteLines)
}
