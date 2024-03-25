package models

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const scanPoint = `-- name: scanPoint :many
SELECT id, created_at, occur_at, external_id, neighborhood, location_type, crime_type, region, point, lat, long
FROM events
WHERE 
ST_DWithin(
    point::geography,
    ST_Point($1, $2, 4326)::geography,
    $3,
	true
)
AND occur_at >= $4
AND occur_at <= $5
ORDER BY occur_at
LIMIT $6
`

type ScanPointParams struct {
	Lat      float64            `json:"lat"`
	Long     float64            `json:"long"`
	Radius   float64            `json:"radius"`
	Region   string             `json:"region"`
	FromDate pgtype.Timestamptz `json:"to_date"`
	ToDate   pgtype.Timestamptz `json:"from_date"`
	Limit    int32              `json:"limit"`
}

func (q *Queries) ScanPoint(ctx context.Context, arg ScanPointParams) ([]Event, error) {
	rows, err := q.db.Query(ctx, scanPoint,
		arg.Long,
		arg.Lat,
		arg.Radius,
		arg.FromDate,
		arg.ToDate,
		arg.Limit,
	)
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
