// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package models

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Area struct {
	ID        string             `json:"id"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UserID    string             `json:"user_id"`
	IsActive  bool               `json:"is_active"`
	Address   string             `json:"address"`
	Radius    int32              `json:"radius"`
	Point     *string            `json:"point"`
	Lat       float64            `json:"lat"`
	Long      float64            `json:"long"`
}

type Event struct {
	ID           int32              `json:"id"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	OccurAt      pgtype.Timestamptz `json:"occur_at"`
	ExternalID   string             `json:"external_id"`
	Neighborhood pgtype.Text        `json:"neighborhood"`
	LocationType pgtype.Text        `json:"location_type"`
	CrimeType    string             `json:"crime_type"`
	Point        *string            `json:"point"`
	Lat          float64            `json:"lat"`
	Long         float64            `json:"long"`
}

type Report struct {
	ID        string             `json:"id"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UserID    string             `json:"user_id"`
	AreaID    string             `json:"area_id"`
}

type ReportEvent struct {
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	ReportID  string             `json:"report_id"`
	EventID   int32              `json:"event_id"`
	AreaID    string             `json:"area_id"`
}

type Scan struct {
	ID          int32              `json:"id"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	Radius      int32              `json:"radius"`
	FromDate    pgtype.Timestamptz `json:"from_date"`
	ToDate      pgtype.Timestamptz `json:"to_date"`
	UserID      pgtype.Text        `json:"user_id"`
	EventsCount int32              `json:"events_count"`
	Address     string             `json:"address"`
	Point       *string            `json:"point"`
	Lat         float64            `json:"lat"`
	Long        float64            `json:"long"`
}

type TempEvent struct {
	ID           int32              `json:"id"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	OccurAt      pgtype.Timestamptz `json:"occur_at"`
	ExternalID   string             `json:"external_id"`
	Neighborhood pgtype.Text        `json:"neighborhood"`
	LocationType pgtype.Text        `json:"location_type"`
	CrimeType    string             `json:"crime_type"`
	Point        *string            `json:"point"`
	Lat          float64            `json:"lat"`
	Long         float64            `json:"long"`
}
