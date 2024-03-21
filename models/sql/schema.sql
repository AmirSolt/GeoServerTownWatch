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
    region TEXT NOT NULL,
    point geometry(Point, 3857),
    lat DOUBLE PRECISION NOT NULL,
    long DOUBLE PRECISION NOT NULL
);
CREATE INDEX event_occ_at_idx ON events ("occur_at");
CREATE INDEX event_point_idx ON events USING GIST ("point");
CREATE TEMPORARY TABLE _temp_events (LIKE events INCLUDING ALL) ON COMMIT DROP;
CREATE FUNCTION event_insert() RETURNS trigger AS $$
    BEGIN
        NEW.point := ST_POINT(NEW.lat, NEW.long, 3857);
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER on_event_insert BEFORE INSERT OR UPDATE ON events
    FOR EACH ROW EXECUTE FUNCTION event_insert();

CREATE TABLE scans (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    radius DOUBLE PRECISION NOT NULL,
    from_date TIMESTAMPTZ NOT NULL,
    to_date TIMESTAMPTZ NOT NULL,

    events_count INT NOT NULL,

    region TEXT NOT NULL,
    point geometry(Point, 3857),
    lat DOUBLE PRECISION NOT NULL,
    long DOUBLE PRECISION NOT NULL
);

CREATE OR REPLACE FUNCTION scan_point(
  lat DOUBLE PRECISION,
  long DOUBLE PRECISION,
  radius DOUBLE PRECISION,
  region TEXT,
  from_date TIMESTAMPTZ,
  to_date TIMESTAMPTZ,
  scan_events_count_limit INT
) RETURNS SETOF events AS $$
DECLARE
    scan_results CURSOR FOR
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
        LIMIT scan_events_count_limit;

    event_row events%ROWTYPE;
    scan_record scans%ROWTYPE;
    events_found_count INT := 0;
BEGIN
    -- Perform the scan and store the results in a temporary table
    CREATE TEMP TABLE temp_scan_results AS
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
        LIMIT scan_events_count_limit;

    -- Iterate over the results to count the events
    FOR event_row IN scan_results LOOP
        events_found_count := events_found_count + 1;
    END LOOP;

    -- Record the completed scan
    scan_record.radius := radius;
    scan_record.from_date := from_date;
    scan_record.to_date := to_date;
    scan_record.events_found_count := events_found_count;
    scan_record.region := region;
    scan_record.point := ST_SetSRID(ST_MakePoint(long, lat), 3857);
    scan_record.lat := lat;
    scan_record.long := long;

    INSERT INTO scans (radius, from_date, to_date, events_found_count, region, point, lat, long)
    VALUES (scan_record.radius, scan_record.from_date, scan_record.to_date, scan_record.events_found_count, scan_record.region, scan_record.point, scan_record.lat, scan_record.long);

    -- Return the original scan results
    RETURN QUERY SELECT * FROM temp_scan_results;

    DROP TABLE IF EXISTS temp_scan_results;

    RETURN;
END;
$$ LANGUAGE plpgsql;
-- ======

CREATE TABLE areas (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id TEXT NOT NULL UNIQUE,
    is_active BOOLEAN NOT NULL DEFAULT true,
    address TEXT NOT NULL,
    region TEXT NOT NULL,
    radius DOUBLE PRECISION NOT NULL,
    point geometry(Point, 3857),
    lat DOUBLE PRECISION NOT NULL,
    long DOUBLE PRECISION NOT NULL
);
CREATE INDEX area_user_idx ON areas ("user_id");
CREATE FUNCTION area_insert() RETURNS trigger AS $$
    BEGIN
        NEW.point := ST_POINT(NEW.lat, NEW.long, 3857);
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
    is_reported BOOLEAN NOT NULL DEFAULT false,
    area_id uuid NOT NULL REFERENCES areas(id),
    CONSTRAINT fk_area FOREIGN KEY (area_id) REFERENCES areas(id) ON DELETE CASCADE ON UPDATE CASCADE
);
CREATE OR REPLACE FUNCTION create_global_reports(from_date TIMESTAMPTZ, to_date TIMESTAMPTZ, scan_events_count_limit INT)
RETURNS SETOF record AS $$
DECLARE
    area_record RECORD;
    event_ids INT[];
    event_id INT;
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
            INSERT INTO reports (area_id, user_id) VALUES (area_record.id, area_record.user_id)
            RETURNING * INTO new_report;
            
            FOREACH event_id IN ARRAY event_ids LOOP
                INSERT INTO report_events (report_id, event_id) VALUES (new_report.id, event_id);
            END LOOP;

            RETURN NEXT new_report;
        END IF;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- ======

CREATE TABLE report_events (
    PRIMARY KEY (report_id, event_id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    report_id uuid NOT NULL REFERENCES reports(id),
    event_id INTEGER NOT NULL REFERENCES events(id),
    CONSTRAINT fk_report FOREIGN KEY(report_id) REFERENCES reports(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_event FOREIGN KEY(event_id) REFERENCES events(id) ON DELETE CASCADE ON UPDATE CASCADE
);
CREATE UNIQUE INDEX report_idx ON report_events("report_id");

-- ======

