-- +goose Up

CREATE TABLE IF NOT EXISTS gtfs_vehicles (
  vehicle_id text NOT NULL PRIMARY KEY,
  vehicle_label text NOT NULL,
  trip_trip_id text NOT NULL,
  trip_route_id text NOT NULL,
  position_lat_lon geography(POINT) NOT NULL,
  position_bearing SMALLINT NOT NULL,
  current_stop_sequence BIGINT NOT NULL,
  stop_id text NOT NULL,
  current_status BIGINT NOT NULL,
  timestamp BIGINT NOT NULL
);


-- +goose Down
DROP TABLE gtfs_vehicles;
