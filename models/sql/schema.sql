-- Amirali Soltani
-- 2024-01-29
CREATE EXTENSION IF NOT EXISTS "postgis";
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";



CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    occur_at TIMESTAMPTZ NOT NULL,
    external_id TEXT NOT NULL UNIQUE,
    details JSONB NOT NULL,
    point geography(Point, 4326),
    lat DOUBLE PRECISION NOT NULL,
    long DOUBLE PRECISION NOT NULL
);
CREATE INDEX event_occ_at_idx ON events ("occur_at");
CREATE INDEX event_point_idx ON events USING GIST ("point");
CREATE TEMPORARY TABLE _temp_events (LIKE events INCLUDING ALL) ON COMMIT DROP;
CREATE FUNCTION event_insert() RETURNS trigger AS $$
    BEGIN
        NEW.point := ST_POINT(NEW.long, NEW.lat, 4326)::geography;
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER on_event_insert BEFORE INSERT OR UPDATE ON events
    FOR EACH ROW EXECUTE FUNCTION event_insert();

CREATE TABLE scans (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    radius INT NOT NULL,
    from_date TIMESTAMPTZ NOT NULL,
    to_date TIMESTAMPTZ NOT NULL,
    events_count INT NOT NULL,
    address TEXT NOT NULL,
    point geography(Point, 4326),
    lat DOUBLE PRECISION NOT NULL,
    long DOUBLE PRECISION NOT NULL
);


