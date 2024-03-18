

-- =========================================
-- events

-- name: CreateEvents :copyfrom
INSERT INTO events (
    occur_at,
    external_id,
    neighborhood,
    location_type,
    crime_type,
    region,
    lat,
    long
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8);

-- =========================================
--  areas


-- name: GetArea :one
SELECT * FROM areas
WHERE id = $1 AND user_id=$2;

-- name: GetAreasByUser :many
SELECT * FROM areas
WHERE user_id = $1;


-- name: CreateArea :exec
INSERT INTO areas (user_id, address, region, radius, lat, long) VALUES ($1,$2,$3,$4,$5,$6);

-- name: UpdateArea :exec
UPDATE areas SET 
address = $3,
radius = $4,
lat = $5,
long = $6
WHERE id = $1 AND user_id = $2;

-- name: DeleteArea :exec
DELETE FROM areas WHERE id = $1 AND user_id=$2;

-- =========================================
-- reports


-- name: GetPublicReportDetails :one
SELECT id, created_at, is_reported  FROM reports
WHERE id = $1;

-- name: GetPrivateReportDetails :one
SELECT sqlc.embed(r), sqlc.embed(a)
FROM reports r
INNER JOIN areas a ON r.area_id = a.id
WHERE r.id = $1 AND r.user_id = $2;

-- name: GetEventsByReport :many
SELECT e.*
FROM events e
INNER JOIN report_events re ON e.id = re.event_id
INNER JOIN reports r ON re.report_id = r.id
WHERE r.id = $1;


-- =========================================
-- custom functions

-- name: ScanPoint :many
SELECT scan_point($1, $2, $3, $4, $5, $6, $7);

-- name: CreateGlobalReports :many
SELECT create_global_reports($1, $2, $3);

