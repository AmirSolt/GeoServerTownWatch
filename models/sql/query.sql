

-- =========================================
-- events


-- name: CreateTempEventsTable :exec
CREATE TEMPORARY TABLE _temp_events (LIKE events INCLUDING ALL) ON COMMIT DROP;


-- name: CreateTempEvents :copyfrom
INSERT INTO _temp_events (
    occur_at,
    external_id,
    neighborhood,
    location_type,
    crime_type,
    lat,
    long
) VALUES ($1,$2,$3,$4,$5,$6,$7);

-- name: MoveFromTempEventsToEvents :exec
INSERT INTO events (
    occur_at,
    external_id,
    neighborhood,
    location_type,
    crime_type,
    lat,
    long
)
SELECT
    occur_at,
    external_id,
    neighborhood,
    location_type,
    crime_type,
    lat,
    long
FROM _temp_events
ON CONFLICT (external_id) DO NOTHING;

-- name: CountEvents :one
SELECT count(*) FROM events;

-- name: CountTempEvents :one
SELECT count(*) FROM _temp_events;

-- name: CreateScan :one
INSERT INTO scans (radius, from_date, to_date, events_count, address, user_id, lat, long) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING *;


-- =========================================
--  areas

-- name: CountAreasByUser :one
SELECT count(*) FROM areas
WHERE user_id = $1;

-- name: GetArea :one
SELECT * FROM areas
WHERE id = $1 AND user_id=$2;

-- name: GetAreasByUser :many
SELECT * FROM areas
WHERE user_id = $1 ORDER BY created_at;


-- name: CreateArea :one
INSERT INTO areas (user_id, address, radius, lat, long) VALUES ($1,$2,$3,$4,$5) RETURNING *;

-- name: UpdateArea :one
UPDATE areas SET 
address = $3,
radius = $4,
lat = $5,
long = $6
WHERE id = $1 AND user_id = $2
RETURNING *;

-- name: DeleteArea :one
DELETE FROM areas WHERE id = $1 AND user_id=$2
RETURNING *;

-- =========================================
-- reports

-- name: GetReportsByUser :many
SELECT * FROM reports
WHERE user_id = $1 ORDER BY created_at;


-- name: GetReportDetails :one
SELECT sqlc.embed(r), sqlc.embed(a)
FROM reports r
INNER JOIN areas a ON r.area_id = a.id
WHERE r.id = $1;

-- name: GetEventsByReport :many
SELECT e.*
FROM events e
INNER JOIN report_events re ON e.id = re.event_id
INNER JOIN reports r ON re.report_id = r.id
WHERE r.id = $1;


-- =========================================
-- custom functions




