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
    external_id TEXT NOT NULL,
    neighborhood TEXT,
    location_type TEXT,
    crime_type crime_type NOT NULL,
    region TEXT NOT NULL,
    point geometry(Point, 3857) NOT NULL,
    lat DOUBLE PRECISION NOT NULL,
    long DOUBLE PRECISION NOT NULL
);
CREATE INDEX event_occ_at_idx ON events ("occur_at");
CREATE INDEX event_point_idx ON events USING GIST ("point");
CREATE FUNCTION event_insert() RETURNS trigger AS $$
    BEGIN
        NEW.point := ST_POINT(NEW.lat, NEW.long, 3857)
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER on_event_insert BEFORE INSERT OR UPDATE ON events
    FOR EACH ROW EXECUTE FUNCTION event_insert();

CREATE FUNCTION scan_custom_area(
    lat DOUBLE PRECISION,
    long DOUBLE PRECISION,
    radius DOUBLE PRECISION,
    region TEXT,
    from_date TIMESTAMPTZ,
    to_date TIMESTAMPTZ,
    count_limit INT
    ) RETURNS SETOF reports AS $$
        RETURN QUERY
        SELECT *
        FROM events
        WHERE 
        ST_DWithin(
            point,
            ST_Point(lat, long, 3857),
           radius
        )
        AND region = region
        AND occur_at >= from_date
        AND occur_at <= to_date
        ORDER BY occur_at
        LIMIT count_limit;
$$ LANGUAGE plpgsql;



-- ======

CREATE TABLE areas (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user TEXT UNIQUE NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    address TEXT NOT NULL,
    region region NOT NULL,
    radius DOUBLE PRECISION NOT NULL,
    point geometry(Point, 3857) NOT NULL,
    lat DOUBLE PRECISION NOT NULL,
    long DOUBLE PRECISION NOT NULL
);
CREATE INDEX area_user_idx ON areas ("user");
CREATE FUNCTION area_insert() RETURNS trigger AS $$
    BEGIN
        NEW.point := ST_POINT(NEW.lat, NEW.long, 3857)
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER on_area_insert BEFORE INSERT OR UPDATE ON areas
    FOR EACH ROW EXECUTE FUNCTION area_insert();


-- ======

CREATE TABLE reports (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    is_reported BOOLEAN NOT NULL DEFAULT false,
    area_id INT NOT NULL REFERENCES areas(id),
    CONSTRAINT fk_area FOREIGN KEY (area_id) REFERENCES areas(id) ON DELETE CASCADE ON UPDATE CASCADE
);
CREATE OR REPLACE FUNCTION scan_areas(from_date TIMESTAMPTZ, to_date TIMESTAMPTZ, scan_events_limit INT)
RETURNS SETOF reports AS $$
DECLARE
    area_record RECORD;
    event_ids INT[];
    new_report RECORD;
BEGIN
    FOR area_record IN SELECT * FROM areas WHERE is_active = true LOOP
        event_ids := ARRAY(SELECT id FROM events
            WHERE ST_DWithin(point, area_record.point, area_record.radius)
            AND region = area_record.region
            AND occur_at >= from_date
            AND occur_at <= to_date
            ORDER BY occur_at
            LIMIT scan_events_limit
        );
        
        IF array_length(event_ids, 1) > 0 THEN
            INSERT INTO reports area_id VALUES area_record.id
            RETURNING * INTO new_report;
            
            FOREACH event_id IN ARRAY event_ids LOOP
                INSERT INTO scans (report_id, event_id) VALUES (new_report.id, event_id);
            END LOOP;
            
            RETURN NEXT new_report;
        END IF;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- ======

CREATE TABLE scans (
    PRIMARY KEY (report_id, event_id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    report_id uuid NOT NULL REFERENCES reports(id),
    event_id INTEGER NOT NULL REFERENCES events(id),
    CONSTRAINT fk_report FOREIGN KEY(report_id) REFERENCES reports(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_event FOREIGN KEY(event_id) REFERENCES events(id) ON DELETE CASCADE ON UPDATE CASCADE
);
CREATE UNIQUE INDEX report_idx ON scans("report_id");
