-- +goose Up
-- +goose StatementBegin

-- Enable PostGIS extension
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS locations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    geom GEOGRAPHY (POINT, 4326),
    created_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create an index on the name for faster lookups
CREATE INDEX IF NOT EXISTS idx_locations_name ON locations (name);

-- Create spatial index for efficient geographic queries
CREATE INDEX IF NOT EXISTS idx_locations_geom ON locations USING GIST (geom);

-- Function to automatically update the geometry column
CREATE OR REPLACE FUNCTION update_geom_column()
  RETURNS TRIGGER AS
$$
BEGIN
  NEW.geom = ST_SetSRID(ST_MakePoint(NEW.longitude, NEW.latitude), 4326)::geography;
  RETURN NEW;
END;
$$
LANGUAGE plpgsql;

-- Trigger to automatically populate geometry column on insert/update
CREATE TRIGGER update_geom BEFORE INSERT OR UPDATE ON locations
  FOR EACH ROW EXECUTE PROCEDURE update_geom_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS update_geom ON locations;

DROP FUNCTION IF EXISTS update_geom_column ();

DROP INDEX IF EXISTS idx_locations_geom;

DROP INDEX IF EXISTS idx_locations_name;

DROP TABLE IF EXISTS locations;

DROP EXTENSION IF EXISTS postgis;

-- +goose StatementEnd