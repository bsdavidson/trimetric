create database trimetric;
\c trimetric
CREATE EXTENSION Postgis;
create table "vehicles"
(
  "id" serial PRIMARY KEY,
  "vehicle_id" integer not null UNIQUE,
  "type" varchar(255) not null,
  "sign_message" varchar(255) not null,
  "time" timestamp not null
);
CREATE TABLE raw_vehicles
(
  vehicle_id integer NOT NULL PRIMARY KEY,
  data jsonb NOT NULL
);
create index idx_data_expires on raw_vehicles
(
((data->>'expires')::bigint));
CREATE TABLE stops
(
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


-- SELECT name, ST_Distance(ST_GeogFromText('SRID=4326;POINT(45.38472 -57.39681)'), lat_lon)
-- AS distance
-- FROM stops
-- WHERE ST_DWithin(ST_GeogFromText('SRID=4326;POINT(45.38472 -57.39681)'), lat_lon, 1000)
-- ORDER BY distance ASC;

