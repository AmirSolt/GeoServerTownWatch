-- Amirali Soltani
-- 2024-01-29
CREATE EXTENSION IF NOT EXISTS "postgis";

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE crime_type AS ENUM (
    'Assault',
    'Auto Theft',
    'Theft from Motor Vehicle',
    'Break and Enter',
    'Sexual Violation',
    'Robbery',
    'Theft Over',
    'Bike Theft',
    'Shooting',
    'Homicide'
);




CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    occur_at TIMESTAMPTZ NOT NULL,
    external_id TEXT NOT NULL UNIQUE,
    neighborhood TEXT,
    location_type TEXT,
    crime_type crime_type NOT NULL,
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
    user_id TEXT,
    events_count INT NOT NULL,
    address TEXT NOT NULL,
    point geography(Point, 4326),
    lat DOUBLE PRECISION NOT NULL,
    long DOUBLE PRECISION NOT NULL
);


-- ======

CREATE TABLE areas (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    address TEXT NOT NULL,
    radius INT NOT NULL,
    point geography(Point, 4326),
    lat DOUBLE PRECISION NOT NULL,
    long DOUBLE PRECISION NOT NULL
);
CREATE INDEX area_user_idx ON areas ("user_id");
CREATE FUNCTION area_insert() RETURNS trigger AS $$
    BEGIN
        NEW.point := ST_POINT(NEW.long, NEW.lat, 4326)::geography;
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER on_area_insert BEFORE INSERT OR UPDATE ON areas
    FOR EACH ROW EXECUTE FUNCTION area_insert();


-- ======

CREATE TABLE reports (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id TEXT NOT NULL UNIQUE,
    area_id uuid NOT NULL REFERENCES areas(id),
    CONSTRAINT fk_area FOREIGN KEY (area_id) REFERENCES areas(id) ON DELETE CASCADE ON UPDATE CASCADE
);



-- ======

CREATE TABLE report_events (
    PRIMARY KEY (report_id, event_id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    report_id uuid NOT NULL REFERENCES reports(id),
    event_id INTEGER NOT NULL REFERENCES events(id),
    CONSTRAINT fk_report FOREIGN KEY(report_id) REFERENCES reports(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_event FOREIGN KEY(event_id) REFERENCES events(id) ON DELETE CASCADE ON UPDATE CASCADE
);
CREATE INDEX report_idx ON report_events("report_id");
CREATE INDEX event_idx ON report_events("event_id");

-- ======

