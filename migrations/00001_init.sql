-- +goose Up
CREATE EXTENSION postgis;

-- This doesn't map to any file in the GTFS dump, but is here to provide
-- a unique key for service_id to tie foreign keys to.
CREATE TABLE services (
  id text NOT NULL PRIMARY KEY
);

CREATE TABLE calendar_dates (
  service_id text NOT NULL,
  date date NOT NULL,
  exception_type SMALLINT NOT NULL,
  PRIMARY KEY (service_id, date),
  FOREIGN KEY (service_id) REFERENCES services (id) ON DELETE CASCADE
);
CREATE INDEX ON calendar_dates (date, service_id);

CREATE TABLE routes (
  id text NOT NULL PRIMARY KEY,
  agency_id text NOT NULL DEFAULT '',
  short_name text NOT NULL DEFAULT '',
  long_name text NOT NULL DEFAULT '',
  type integer NOT NULL DEFAULT 0,
  url text NOT NULL DEFAULT '',
  color text NOT NULL DEFAULT '',
  text_color text NOT NULL,
  sort_order integer NOT NULL
);

CREATE TABLE trips (
  id text NOT NULL PRIMARY KEY,
  route_id text NOT NULL,
  service_id text NOT NULL,
  direction_id SMALLINT,
  block_id text,
  shape_id text,
  headsign text,
  short_name text,
  bikes_allowed SMALLINT NOT NULL,
  wheelchair_accessible SMALLINT NOT NULL,
  FOREIGN KEY (route_id) REFERENCES routes (id) ON DELETE CASCADE,
  FOREIGN KEY (service_id) REFERENCES services (id) ON DELETE CASCADE
);

CREATE TABLE stops (
  id text NOT NULL PRIMARY KEY,
  code text NOT NULL DEFAULT '',
  name text NOT NULL,
  "desc" text NOT NULL DEFAULT '',
  lat_lon geography(POINT) NOT NULL,
  zone_id text NOT NULL DEFAULT '',
  url text NOT NULL DEFAULT '',
  location_type integer NOT NULL DEFAULT 0,
  parent_station text NOT NULL DEFAULT 0,
  direction text NOT NULL DEFAULT '',
  position text NOT NULL DEFAULT '',
  wheelchair_boarding integer NOT NULL DEFAULT 0
);

CREATE TABLE stop_times (
  trip_id text NOT NULL,
  arrival_time INTERVAL,
  departure_time INTERVAL,
  stop_id text NOT NULL,
  stop_sequence smallint NOT NULL,
  stop_headsign text,
  pickup_type SMALLINT NOT NULL DEFAULT 0,
  drop_off_type SMALLINT NOT NULL DEFAULT 0,
  shape_dist_traveled DOUBLE PRECISION,
  timepoint SMALLINT,
  continuous_drop_off SMALLINT NOT NULL DEFAULT 0,
  continuous_pickup SMALLINT NOT NULL DEFAULT 0,
  PRIMARY KEY (trip_id, stop_sequence),
  FOREIGN KEY (trip_id) REFERENCES trips (id) ON DELETE CASCADE,
  FOREIGN KEY (stop_id) REFERENCES stops (id) ON DELETE CASCADE
);
CREATE INDEX ON stop_times (arrival_time);
CREATE INDEX ON stop_times (stop_id, arrival_time);

CREATE TABLE vehicle_positions (
  vehicle_id text NOT NULL PRIMARY KEY,
  vehicle_label text,
  trip_id text,
  route_id text,
  stop_id text NOT NULL,
  current_stop_sequence BIGINT NOT NULL,
  current_status BIGINT NOT NULL,
  position_lon_lat geography(POINT) NOT NULL,
  position_bearing SMALLINT,
  timestamp BIGINT NOT NULL
);
CREATE INDEX ON vehicle_positions (trip_id, vehicle_id);

CREATE TABLE trip_updates (
  id serial NOT NULL PRIMARY KEY,
  trip_id text NOT NULL,
  route_id text,
  vehicle_id text,
  vehicle_label text,
  timestamp timestamp,
  delay integer
);

CREATE TABLE stop_time_updates (
  trip_update_id bigint NOT NULL,
  index smallint NOT NULL,
  stop_sequence integer,
  stop_id text,
  arrival_delay integer,
	arrival_time timestamp,
  arrival_uncertainty smallint,
  departure_delay integer,
  departure_time timestamp,
  departure_uncertainty smallint,
  schedule_relationship integer,
  PRIMARY KEY (trip_update_id, index),
  FOREIGN KEY (trip_update_id) REFERENCES trip_updates (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE stop_time_updates;
DROP TABLE trip_updates;
DROP TABLE vehicle_positions;
DROP TABLE stop_times;
DROP TABLE stops;
DROP TABLE trips;
DROP TABLE routes;
DROP TABLE calendar_dates;
DROP EXTENSION postgis;