

-- =========================================
-- events


-- name: CreateTempEventsTable :exec
CREATE TEMPORARY TABLE _temp_events (LIKE events INCLUDING ALL) ON COMMIT DROP;


-- name: CreateTempEvents :copyfrom
INSERT INTO _temp_events (
    occur_at,
    external_id,
    details,
    lat,
    long
) VALUES ($1,$2,$3,$4,$5);

-- name: MoveFromTempEventsToEvents :exec
INSERT INTO events (
    occur_at,
    external_id,
    details,
    lat,
    long
)
SELECT
    occur_at,
    external_id,
    details,
    lat,
    long
FROM _temp_events
ON CONFLICT (external_id) DO NOTHING;

-- name: CountEvents :one
SELECT count(*) FROM events;

-- name: CountTempEvents :one
SELECT count(*) FROM _temp_events;

-- name: CreateScan :one
INSERT INTO scans (radius, from_date, to_date, events_count, address, lat, long) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING *;

