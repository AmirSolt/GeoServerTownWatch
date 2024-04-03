package models

import (
	"context"
	"fmt"
	"townwatch/utils"

	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgx/v5/pgtype"
)

const scanPoint = `-- name: scanPoint :many
SELECT id, created_at, occur_at, external_id, neighborhood, location_type, crime_type, point, lat, long
FROM events
WHERE 
ST_DWithin(
    point::geography,
    ST_Point($1, $2, 4326)::geography,
    CAST ($3 AS DOUBLE PRECISION),
	false
)
AND occur_at >= $4
AND occur_at <= $5
ORDER BY occur_at
LIMIT $6
`

type ScanPointParams struct {
	Lat      float64            `json:"lat"`
	Long     float64            `json:"long"`
	Radius   int32              `json:"radius"`
	FromDate pgtype.Timestamptz `json:"to_date"`
	ToDate   pgtype.Timestamptz `json:"from_date"`
	Limit    int32              `json:"limit"`
}

func (q *Queries) ScanPoint(ctx context.Context, arg ScanPointParams) ([]Event, *utils.CError) {
	rows, err := q.db.Query(ctx, scanPoint,
		arg.Long,
		arg.Lat,
		arg.Radius,
		arg.FromDate,
		arg.ToDate,
		arg.Limit,
	)
	if err != nil {
		eventID := sentry.CaptureException(err)
		cerr := &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
		return nil, cerr
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
			eventID := sentry.CaptureException(err)
			cerr := &utils.CError{
				EventID: eventID,
				Message: "Internal Server Error",
				Error:   err,
			}
			return nil, cerr
		}

		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		eventID := sentry.CaptureException(err)
		cerr := &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
		return nil, cerr
	}
	return items, nil
}

// ==============================================

const createGlobalReports = `
-- name: createGlobalReports :many
WITH scanned AS(
	SELECT
	a.id as area_id,
	a.user_id,
	e.id as event_id
	FROM
	areas a
	JOIN
	events e 
	ON ST_DWithin(e.point, a.point, a.radius, false)
	AND e.created_at >= $1
	AND e.created_at <= $2 	
	AND NOT EXISTS (
	  SELECT 1
	  FROM report_events
	  WHERE report_events.event_id = e.id AND report_events.area_id = a.id
	)
	ORDER BY e.created_at
	LIMIT $3
  ),
  inserted_reports AS (
	  INSERT INTO reports (area_id, user_id)
	  SELECT DISTINCT ON (scanned.area_id) scanned.area_id, scanned.user_id
	  FROM
		  scanned
		RETURNING *
  ),
  inserted_events_report AS (
   INSERT INTO report_events (report_id, area_id, event_id)
	SELECT inserted_reports.id, inserted_reports.area_id, scanned.event_id
	FROM
		scanned
	JOIN
	  inserted_reports ON inserted_reports.area_id = scanned.area_id
  )
  select * from inserted_reports
`

type CreateGlobalReportsParams struct {
	FromDate    pgtype.Timestamptz `json:"from_date"`
	ToDate      pgtype.Timestamptz `json:"to_date"`
	EventsLimit int32              `json:"events_limit"`
}

// =========================================
// custom functions
func (q *Queries) CreateGlobalReports(ctx context.Context, arg CreateGlobalReportsParams) ([]Report, *utils.CError) {
	rows, err := q.db.Query(ctx, createGlobalReports, arg.FromDate, arg.ToDate, arg.EventsLimit)
	if err != nil {
		eventID := sentry.CaptureException(err)
		cerr := &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
		return nil, cerr
	}
	defer rows.Close()
	var items []Report
	for rows.Next() {

		fmt.Println("=====----=-=--=-")
		columnValues, _ := rows.Values()
		for i, v := range columnValues {
			fmt.Printf("Type of value at %v=%T, value=%v | \n", i, v, v)
		}
		fmt.Println("=====----=-=--=-")

		var i Report
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UserID,
			&i.AreaID,
		); err != nil {
			eventID := sentry.CaptureException(err)
			cerr := &utils.CError{
				EventID: eventID,
				Message: "Internal Server Error",
				Error:   err,
			}
			return nil, cerr
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		eventID := sentry.CaptureException(err)
		cerr := &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
		return nil, cerr
	}
	return items, nil
}
