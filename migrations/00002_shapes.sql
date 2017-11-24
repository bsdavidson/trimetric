-- +goose Up
CREATE TABLE shapes (
  id TEXT NOT NULL,
  pt_lon_lat geography(POINT) NOT NULL,
  pt_sequence integer NOT NULL,
  dist_traveled DOUBLE PRECISION,
  PRIMARY KEY (id, pt_sequence)
);

-- +goose Down
DROP TABLE shapes;