// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package models

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

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

const createArea = `-- name: CreateArea :exec
INSERT INTO areas (user_id, address, region, radius, lat, long) VALUES ($1,$2,$3,$4,$5,$6)
`

type CreateAreaParams struct {
	UserID  string  `json:"user_id"`
	Address string  `json:"address"`
	Region  string  `json:"region"`
	Radius  float64 `json:"radius"`
	Lat     float64 `json:"lat"`
	Long    float64 `json:"long"`
}

func (q *Queries) CreateArea(ctx context.Context, arg CreateAreaParams) error {
	_, err := q.db.Exec(ctx, createArea,
		arg.UserID,
		arg.Address,
		arg.Region,
		arg.Radius,
		arg.Lat,
		arg.Long,
	)
	return err
}

const createGlobalReports = `-- name: CreateGlobalReports :many
SELECT create_global_reports($1, $2, $3)
`

type CreateGlobalReportsParams struct {
	FromDate             pgtype.Timestamptz `json:"from_date"`
	ToDate               pgtype.Timestamptz `json:"to_date"`
	ScanEventsCountLimit int32              `json:"scan_events_count_limit"`
}

func (q *Queries) CreateGlobalReports(ctx context.Context, arg CreateGlobalReportsParams) ([]interface{}, error) {
	rows, err := q.db.Query(ctx, createGlobalReports, arg.FromDate, arg.ToDate, arg.ScanEventsCountLimit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []interface{}
	for rows.Next() {
		var create_global_reports interface{}
		if err := rows.Scan(&create_global_reports); err != nil {
			return nil, err
		}
		items = append(items, create_global_reports)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

type CreateTempEventsParams struct {
	OccurAt      pgtype.Timestamptz `json:"occur_at"`
	ExternalID   string             `json:"external_id"`
	Neighborhood pgtype.Text        `json:"neighborhood"`
	LocationType pgtype.Text        `json:"location_type"`
	CrimeType    CrimeType          `json:"crime_type"`
	Region       string             `json:"region"`
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

const deleteArea = `-- name: DeleteArea :exec
DELETE FROM areas WHERE id = $1 AND user_id=$2
`

type DeleteAreaParams struct {
	ID     int32  `json:"id"`
	UserID string `json:"user_id"`
}

func (q *Queries) DeleteArea(ctx context.Context, arg DeleteAreaParams) error {
	_, err := q.db.Exec(ctx, deleteArea, arg.ID, arg.UserID)
	return err
}

const getArea = `-- name: GetArea :one


SELECT id, created_at, user_id, is_active, address, region, radius, point, lat, long FROM areas
WHERE id = $1 AND user_id=$2
`

type GetAreaParams struct {
	ID     int32  `json:"id"`
	UserID string `json:"user_id"`
}

// =========================================
//
//	areas
func (q *Queries) GetArea(ctx context.Context, arg GetAreaParams) (Area, error) {
	row := q.db.QueryRow(ctx, getArea, arg.ID, arg.UserID)
	var i Area
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UserID,
		&i.IsActive,
		&i.Address,
		&i.Region,
		&i.Radius,
		&i.Point,
		&i.Lat,
		&i.Long,
	)
	return i, err
}

const getAreasByUser = `-- name: GetAreasByUser :many
SELECT id, created_at, user_id, is_active, address, region, radius, point, lat, long FROM areas
WHERE user_id = $1
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
			&i.Region,
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
SELECT e.id, e.created_at, e.occur_at, e.external_id, e.neighborhood, e.location_type, e.crime_type, e.region, e.point, e.lat, e.long
FROM events e
INNER JOIN report_events re ON e.id = re.event_id
INNER JOIN reports r ON re.report_id = r.id
WHERE r.id = $1
`

func (q *Queries) GetEventsByReport(ctx context.Context, id pgtype.UUID) ([]Event, error) {
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
			&i.Region,
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

const getPrivateReportDetails = `-- name: GetPrivateReportDetails :one
SELECT r.id, r.created_at, r.user_id, r.is_reported, r.area_id, a.id, a.created_at, a.user_id, a.is_active, a.address, a.region, a.radius, a.point, a.lat, a.long
FROM reports r
INNER JOIN areas a ON r.area_id = a.id
WHERE r.id = $1 AND r.user_id = $2
`

type GetPrivateReportDetailsParams struct {
	ID     pgtype.UUID `json:"id"`
	UserID string      `json:"user_id"`
}

type GetPrivateReportDetailsRow struct {
	Report Report `json:"report"`
	Area   Area   `json:"area"`
}

func (q *Queries) GetPrivateReportDetails(ctx context.Context, arg GetPrivateReportDetailsParams) (GetPrivateReportDetailsRow, error) {
	row := q.db.QueryRow(ctx, getPrivateReportDetails, arg.ID, arg.UserID)
	var i GetPrivateReportDetailsRow
	err := row.Scan(
		&i.Report.ID,
		&i.Report.CreatedAt,
		&i.Report.UserID,
		&i.Report.IsReported,
		&i.Report.AreaID,
		&i.Area.ID,
		&i.Area.CreatedAt,
		&i.Area.UserID,
		&i.Area.IsActive,
		&i.Area.Address,
		&i.Area.Region,
		&i.Area.Radius,
		&i.Area.Point,
		&i.Area.Lat,
		&i.Area.Long,
	)
	return i, err
}

const getPublicReportDetails = `-- name: GetPublicReportDetails :one


SELECT id, created_at, is_reported  FROM reports
WHERE id = $1
`

type GetPublicReportDetailsRow struct {
	ID         pgtype.UUID        `json:"id"`
	CreatedAt  pgtype.Timestamptz `json:"created_at"`
	IsReported bool               `json:"is_reported"`
}

// =========================================
// reports
func (q *Queries) GetPublicReportDetails(ctx context.Context, id pgtype.UUID) (GetPublicReportDetailsRow, error) {
	row := q.db.QueryRow(ctx, getPublicReportDetails, id)
	var i GetPublicReportDetailsRow
	err := row.Scan(&i.ID, &i.CreatedAt, &i.IsReported)
	return i, err
}

const moveFromTempEventsToEvents = `-- name: MoveFromTempEventsToEvents :exec
INSERT INTO events (
    occur_at,
    external_id,
    neighborhood,
    location_type,
    crime_type,
    region,
    lat,
    long
)
SELECT
    occur_at,
    external_id,
    neighborhood,
    location_type,
    crime_type,
    region,
    lat,
    long
FROM _temp_events
ON CONFLICT (external_id) DO UPDATE
    SET
        occur_at = EXCLUDED.occur_at,
        neighborhood = EXCLUDED.neighborhood,
        location_type = EXCLUDED.location_type,
        crime_type = EXCLUDED.crime_type,
        region = EXCLUDED.region,
        lat = EXCLUDED.lat,
        long = EXCLUDED.long
`

func (q *Queries) MoveFromTempEventsToEvents(ctx context.Context) error {
	_, err := q.db.Exec(ctx, moveFromTempEventsToEvents)
	return err
}

const scanPoint = `-- name: ScanPoint :many

SELECT scan_point($1, $2, $3, $4, $5, $6, $7)
`

type ScanPointParams struct {
	Lat                  float64            `json:"lat"`
	Long                 float64            `json:"long"`
	Radius               float64            `json:"radius"`
	Region               string             `json:"region"`
	FromDate             pgtype.Timestamptz `json:"from_date"`
	ToDate               pgtype.Timestamptz `json:"to_date"`
	ScanEventsCountLimit int32              `json:"scan_events_count_limit"`
}

// =========================================
// custom functions
func (q *Queries) ScanPoint(ctx context.Context, arg ScanPointParams) ([]interface{}, error) {
	rows, err := q.db.Query(ctx, scanPoint,
		arg.Lat,
		arg.Long,
		arg.Radius,
		arg.Region,
		arg.FromDate,
		arg.ToDate,
		arg.ScanEventsCountLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []interface{}
	for rows.Next() {
		var scan_point interface{}
		if err := rows.Scan(&scan_point); err != nil {
			return nil, err
		}
		items = append(items, scan_point)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateArea = `-- name: UpdateArea :exec
UPDATE areas SET 
address = $3,
radius = $4,
lat = $5,
long = $6
WHERE id = $1 AND user_id = $2
`

type UpdateAreaParams struct {
	ID      int32   `json:"id"`
	UserID  string  `json:"user_id"`
	Address string  `json:"address"`
	Radius  float64 `json:"radius"`
	Lat     float64 `json:"lat"`
	Long    float64 `json:"long"`
}

func (q *Queries) UpdateArea(ctx context.Context, arg UpdateAreaParams) error {
	_, err := q.db.Exec(ctx, updateArea,
		arg.ID,
		arg.UserID,
		arg.Address,
		arg.Radius,
		arg.Lat,
		arg.Long,
	)
	return err
}
