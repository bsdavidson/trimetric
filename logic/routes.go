package logic

import (
	"database/sql"

	"github.com/bsdavidson/trimetric/trimet"
	"github.com/pkg/errors"
)

// RouteDataset provides methods to query and update a database table of Shapes
type RouteDataset interface {
	FetchRoutes() ([]trimet.Route, error)
}

// RouteSQLDataset stores a DB instance and provides access to methods to
// retrieve and update shapes from the database
type RouteSQLDataset struct {
	DB *sql.DB
}

// FetchRoutes ...
func (sd *RouteSQLDataset) FetchRoutes() ([]trimet.Route, error) {
	q := `
		SELECT id,agency_id,short_name,long_name,type,url,color,text_color,sort_order
		FROM routes
		ORDER BY id ASC
	`
	var routes []trimet.Route
	rows, err := sd.DB.Query(q)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer rows.Close()
	for rows.Next() {
		var r trimet.Route
		err := rows.Scan(&r.RouteID, &r.AgencyID, &r.ShortName, &r.LongName, &r.Type, &r.URL, &r.Color, &r.TextColor, &r.SortOrder)
		if err != nil {
			return nil, err
		}
		routes = append(routes, r)

	}
	return routes, nil
}
