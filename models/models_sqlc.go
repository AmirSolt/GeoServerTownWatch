// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Event struct {
	ID         int32              `json:"id"`
	CreatedAt  pgtype.Timestamptz `json:"created_at"`
	OccurAt    pgtype.Timestamptz `json:"occur_at"`
	ExternalID string             `json:"external_id"`
	Details    []byte             `json:"details"`
	Point      *string            `json:"point"`
	Lat        float64            `json:"lat"`
	Long       float64            `json:"long"`
}

type Scan struct {
	ID          int32              `json:"id"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	Radius      int32              `json:"radius"`
	FromDate    pgtype.Timestamptz `json:"from_date"`
	ToDate      pgtype.Timestamptz `json:"to_date"`
	EventsCount int32              `json:"events_count"`
	Address     string             `json:"address"`
	Point       *string            `json:"point"`
	Lat         float64            `json:"lat"`
	Long        float64            `json:"long"`
}

type TempEvent struct {
	ID         int32              `json:"id"`
	CreatedAt  pgtype.Timestamptz `json:"created_at"`
	OccurAt    pgtype.Timestamptz `json:"occur_at"`
	ExternalID string             `json:"external_id"`
	Details    []byte             `json:"details"`
	Point      *string            `json:"point"`
	Lat        float64            `json:"lat"`
	Long       float64            `json:"long"`
}
