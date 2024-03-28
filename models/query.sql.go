// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package models

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countAreasByUser = `-- name: CountAreasByUser :one

SELECT count(*) FROM areas
WHERE user_id = $1
`

// =========================================
//
//	areas
func (q *Queries) CountAreasByUser(ctx context.Context, userID string) (int64, error) {
	row := q.db.QueryRow(ctx, countAreasByUser, userID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countEvents = `-- name: CountEvents :one
SELECT count(*) FROM events
`

func (q *Queries) CountEvents(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countEvents)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countTempEvents = `-- name: CountTempEvents :one
SELECT count(*) FROM _temp_events
`

func (q *Queries) CountTempEvents(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countTempEvents)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createArea = `-- name: CreateArea :one
INSERT INTO areas (user_id, address, radius, lat, long) VALUES ($1,$2,$3,$4,$5) RETURNING id, created_at, user_id, is_active, address, radius, point, lat, long
`

type CreateAreaParams struct {
	UserID  string  `json:"user_id"`
	Address string  `json:"address"`
	Radius  int32   `json:"radius"`
	Lat     float64 `json:"lat"`
	Long    float64 `json:"long"`
}

func (q *Queries) CreateArea(ctx context.Context, arg CreateAreaParams) (Area, error) {
	row := q.db.QueryRow(ctx, createArea,
		arg.UserID,
		arg.Address,
		arg.Radius,
		arg.Lat,
		arg.Long,
	)
	var i Area
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UserID,
		&i.IsActive,
		&i.Address,
		&i.Radius,
		&i.Point,
		&i.Lat,
		&i.Long,
	)
	return i, err
}

const createScan = `-- name: CreateScan :one
INSERT INTO scans (radius, from_date, to_date, events_count, address, user_id, lat, long) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id, created_at, radius, from_date, to_date, user_id, events_count, address, point, lat, long
`

type CreateScanParams struct {
	Radius      int32              `json:"radius"`
	FromDate    pgtype.Timestamptz `json:"from_date"`
	ToDate      pgtype.Timestamptz `json:"to_date"`
	EventsCount int32              `json:"events_count"`
	Address     string             `json:"address"`
	UserID      pgtype.Text        `json:"user_id"`
	Lat         float64            `json:"lat"`
	Long        float64            `json:"long"`
}

func (q *Queries) CreateScan(ctx context.Context, arg CreateScanParams) (Scan, error) {
	row := q.db.QueryRow(ctx, createScan,
		arg.Radius,
		arg.FromDate,
		arg.ToDate,
		arg.EventsCount,
		arg.Address,
		arg.UserID,
		arg.Lat,
		arg.Long,
	)
	var i Scan
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Radius,
		&i.FromDate,
		&i.ToDate,
		&i.UserID,
		&i.EventsCount,
		&i.Address,
		&i.Point,
		&i.Lat,
		&i.Long,
	)
	return i, err
}

type CreateTempEventsParams struct {
	OccurAt      pgtype.Timestamptz `json:"occur_at"`
	ExternalID   string             `json:"external_id"`
	Neighborhood pgtype.Text        `json:"neighborhood"`
	LocationType pgtype.Text        `json:"location_type"`
	CrimeType    CrimeType          `json:"crime_type"`
	Lat          float64            `json:"lat"`
	Long         float64            `json:"long"`
}

const createTempEventsTable = `-- name: CreateTempEventsTable :exec


CREATE TEMPORARY TABLE _temp_events (LIKE events INCLUDING ALL) ON COMMIT DROP
`

// =========================================
// events
func (q *Queries) CreateTempEventsTable(ctx context.Context) error {
	_, err := q.db.Exec(ctx, createTempEventsTable)
	return err
}

const deleteArea = `-- name: DeleteArea :one
DELETE FROM areas WHERE id = $1 AND user_id=$2
RETURNING id, created_at, user_id, is_active, address, radius, point, lat, long
`

type DeleteAreaParams struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

func (q *Queries) DeleteArea(ctx context.Context, arg DeleteAreaParams) (Area, error) {
	row := q.db.QueryRow(ctx, deleteArea, arg.ID, arg.UserID)
	var i Area
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UserID,
		&i.IsActive,
		&i.Address,
		&i.Radius,
		&i.Point,
		&i.Lat,
		&i.Long,
	)
	return i, err
}

const getArea = `-- name: GetArea :one
SELECT id, created_at, user_id, is_active, address, radius, point, lat, long FROM areas
WHERE id = $1 AND user_id=$2
`

type GetAreaParams struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

func (q *Queries) GetArea(ctx context.Context, arg GetAreaParams) (Area, error) {
	row := q.db.QueryRow(ctx, getArea, arg.ID, arg.UserID)
	var i Area
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UserID,
		&i.IsActive,
		&i.Address,
		&i.Radius,
		&i.Point,
		&i.Lat,
		&i.Long,
	)
	return i, err
}

const getAreasByUser = `-- name: GetAreasByUser :many
SELECT id, created_at, user_id, is_active, address, radius, point, lat, long FROM areas
WHERE user_id = $1 ORDER BY created_at
`

func (q *Queries) GetAreasByUser(ctx context.Context, userID string) ([]Area, error) {
	rows, err := q.db.Query(ctx, getAreasByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Area
	for rows.Next() {
		var i Area
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UserID,
			&i.IsActive,
			&i.Address,
			&i.Radius,
			&i.Point,
			&i.Lat,
			&i.Long,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEventsByReport = `-- name: GetEventsByReport :many
SELECT e.id, e.created_at, e.occur_at, e.external_id, e.neighborhood, e.location_type, e.crime_type, e.point, e.lat, e.long
FROM events e
INNER JOIN report_events re ON e.id = re.event_id
INNER JOIN reports r ON re.report_id = r.id
WHERE r.id = $1
`

func (q *Queries) GetEventsByReport(ctx context.Context, id string) ([]Event, error) {
	rows, err := q.db.Query(ctx, getEventsByReport, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.OccurAt,
			&i.ExternalID,
			&i.Neighborhood,
			&i.LocationType,
			&i.CrimeType,
			&i.Point,
			&i.Lat,
			&i.Long,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getReportDetails = `-- name: GetReportDetails :one
SELECT r.id, r.created_at, r.user_id, r.area_id, a.id, a.created_at, a.user_id, a.is_active, a.address, a.radius, a.point, a.lat, a.long
FROM reports r
INNER JOIN areas a ON r.area_id = a.id
WHERE r.id = $1
`

type GetReportDetailsRow struct {
	Report Report `json:"report"`
	Area   Area   `json:"area"`
}

func (q *Queries) GetReportDetails(ctx context.Context, id string) (GetReportDetailsRow, error) {
	row := q.db.QueryRow(ctx, getReportDetails, id)
	var i GetReportDetailsRow
	err := row.Scan(
		&i.Report.ID,
		&i.Report.CreatedAt,
		&i.Report.UserID,
		&i.Report.AreaID,
		&i.Area.ID,
		&i.Area.CreatedAt,
		&i.Area.UserID,
		&i.Area.IsActive,
		&i.Area.Address,
		&i.Area.Radius,
		&i.Area.Point,
		&i.Area.Lat,
		&i.Area.Long,
	)
	return i, err
}

const getReportsByArea = `-- name: GetReportsByArea :many
SELECT id, created_at, user_id, area_id FROM reports
WHERE user_id = $1 AND area_id = $2 ORDER BY created_at
`

type GetReportsByAreaParams struct {
	UserID string `json:"user_id"`
	AreaID string `json:"area_id"`
}

func (q *Queries) GetReportsByArea(ctx context.Context, arg GetReportsByAreaParams) ([]Report, error) {
	rows, err := q.db.Query(ctx, getReportsByArea, arg.UserID, arg.AreaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Report
	for rows.Next() {
		var i Report
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UserID,
			&i.AreaID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getReportsByUser = `-- name: GetReportsByUser :many

SELECT id, created_at, user_id, area_id FROM reports
WHERE user_id = $1 ORDER BY created_at
`

// =========================================
// reports
func (q *Queries) GetReportsByUser(ctx context.Context, userID string) ([]Report, error) {
	rows, err := q.db.Query(ctx, getReportsByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Report
	for rows.Next() {
		var i Report
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UserID,
			&i.AreaID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const moveFromTempEventsToEvents = `-- name: MoveFromTempEventsToEvents :exec
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
ON CONFLICT (external_id) DO NOTHING
`

func (q *Queries) MoveFromTempEventsToEvents(ctx context.Context) error {
	_, err := q.db.Exec(ctx, moveFromTempEventsToEvents)
	return err
}

const updateArea = `-- name: UpdateArea :one
UPDATE areas SET 
address = $3,
radius = $4,
lat = $5,
long = $6
WHERE id = $1 AND user_id = $2
RETURNING id, created_at, user_id, is_active, address, radius, point, lat, long
`

type UpdateAreaParams struct {
	ID      string  `json:"id"`
	UserID  string  `json:"user_id"`
	Address string  `json:"address"`
	Radius  int32   `json:"radius"`
	Lat     float64 `json:"lat"`
	Long    float64 `json:"long"`
}

func (q *Queries) UpdateArea(ctx context.Context, arg UpdateAreaParams) (Area, error) {
	row := q.db.QueryRow(ctx, updateArea,
		arg.ID,
		arg.UserID,
		arg.Address,
		arg.Radius,
		arg.Lat,
		arg.Long,
	)
	var i Area
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UserID,
		&i.IsActive,
		&i.Address,
		&i.Radius,
		&i.Point,
		&i.Lat,
		&i.Long,
	)
	return i, err
}
