-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE EXTENSION IF NOT EXISTS Postgis;

CREATE TABLE IF NOT EXISTS raw_vehicles (
  vehicle_id integer NOT NULL PRIMARY KEY,
  data jsonb NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_raw_vehicles_data_expires ON raw_vehicles
  (((data->>'expires')::bigint));

CREATE TABLE IF NOT EXISTS stops (
  id text NOT NULL PRIMARY KEY,
  code text NOT NULL DEFAULT '',
  name text NOT NULL,
  "desc" text NOT NULL DEFAULT '',
  lat_lon geography(POINT) NOT NULL,
  zone_id text NOT NULL DEFAULT '',
  stop_url text NOT NULL DEFAULT '',
  location_type integer NOT NULL DEFAULT 0,
  parent_station text NOT NULL DEFAULT 0,
  direction text NOT NULL DEFAULT '',
  position text NOT NULL DEFAULT '',
  wheelchair_boarding integer NOT NULL DEFAULT 0
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE stops;
DROP INDEX IF EXISTS idx_data_expires;
DROP TABLE raw_vehicles;
DROP EXTENSION IF EXISTS Postgis;