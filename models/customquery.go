package models

import (
	"context"
	"townwatch/utils"

	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgx/v5/pgtype"
)

const scanPoint = `-- name: scanPoint :many
SELECT id, created_at, occur_at, external_id, details, point, lat, long
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
ORDER BY occur_at DESC
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
			&i.Details,
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
