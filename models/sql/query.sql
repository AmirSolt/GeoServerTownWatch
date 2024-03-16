


-- name: CreateEvents :copyfrom
INSERT INTO events (
    occur_at,
    external_src_id,
    neighborhood,
    location_type,
    crime_type,
    region,
    lat,
    long
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8);


-- name: ScanCustomArea :many
SELECT scan_custom_area($1, $2, $3, $4, $5, $6, $7);

-- name: ScanAreas :many
SELECT scan_areas($1, $2, $3);


