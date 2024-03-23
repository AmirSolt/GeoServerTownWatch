// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package models

import (
	"database/sql/driver"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
)

type CrimeType string

const (
	CrimeTypeAssault               CrimeType = "Assault"
	CrimeTypeAutoTheft             CrimeType = "Auto Theft"
	CrimeTypeTheftfromMotorVehicle CrimeType = "Theft from Motor Vehicle"
	CrimeTypeBreakandEnter         CrimeType = "Break and Enter"
	CrimeTypeSexualViolation       CrimeType = "Sexual Violation"
	CrimeTypeRobbery               CrimeType = "Robbery"
	CrimeTypeTheftOver             CrimeType = "Theft Over"
	CrimeTypeBikeTheft             CrimeType = "Bike Theft"
	CrimeTypeShooting              CrimeType = "Shooting"
	CrimeTypeHomicide              CrimeType = "Homicide"
)

func (e *CrimeType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = CrimeType(s)
	case string:
		*e = CrimeType(s)
	default:
		return fmt.Errorf("unsupported scan type for CrimeType: %T", src)
	}
	return nil
}

type NullCrimeType struct {
	CrimeType CrimeType `json:"crime_type"`
	Valid     bool      `json:"valid"` // Valid is true if CrimeType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullCrimeType) Scan(value interface{}) error {
	if value == nil {
		ns.CrimeType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.CrimeType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullCrimeType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.CrimeType), nil
}

type Area struct {
	ID        pgtype.UUID        `json:"id"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UserID    string             `json:"user_id"`
	IsActive  bool               `json:"is_active"`
	Address   string             `json:"address"`
	Region    string             `json:"region"`
	Radius    float64            `json:"radius"`
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
	CrimeType    CrimeType          `json:"crime_type"`
	Region       string             `json:"region"`
	Point        *string            `json:"point"`
	Lat          float64            `json:"lat"`
	Long         float64            `json:"long"`
}

type Report struct {
	ID         pgtype.UUID        `json:"id"`
	CreatedAt  pgtype.Timestamptz `json:"created_at"`
	UserID     string             `json:"user_id"`
	IsReported bool               `json:"is_reported"`
	AreaID     pgtype.UUID        `json:"area_id"`
}

type ReportEvent struct {
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	ReportID  pgtype.UUID        `json:"report_id"`
	EventID   int32              `json:"event_id"`
}

type Scan struct {
	ID          int32              `json:"id"`
	CreatedAt   pgtype.Timestamptz `json:"created_at"`
	Radius      float64            `json:"radius"`
	FromDate    pgtype.Timestamptz `json:"from_date"`
	ToDate      pgtype.Timestamptz `json:"to_date"`
	EventsCount int32              `json:"events_count"`
	Region      string             `json:"region"`
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
	CrimeType    CrimeType          `json:"crime_type"`
	Region       string             `json:"region"`
	Point        *string            `json:"point"`
	Lat          float64            `json:"lat"`
	Long         float64            `json:"long"`
}
