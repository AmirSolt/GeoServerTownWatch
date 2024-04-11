package arcgis

import (
	"context"
	"fmt"
	"time"
	"townwatch/base"
	"townwatch/models"
	"townwatch/utils"

	"github.com/jackc/pgx/v5/pgtype"
)

type PeelArcgisAttributes struct {
	OBJECTID          int64  `json:"OBJECTID"`
	OccurrenceNumber  string `json:"OccurrenceNumber"`
	OccurrenceDate    string `json:"OccurrenceDate"`
	OccurrenceTime    string `json:"OccurrenceTime"`
	OccDate           int64  `json:"OccDate"`
	OccDateUTC        int64  `json:"OccDateUTC"`
	OccurrenceWeekday string `json:"OccurrenceWeekday"`
	OccurrenceHour    string `json:"OccurrenceHour"`
	OccMonth          string `json:"OccMonth"`
	OccYear           int64  `json:"OccYear"`
	Description       string `json:"Description"`
	ClearanceStatus   string `json:"ClearanceStatus"`
	StreetName        string `json:"StreetName"`
	StreetType        string `json:"StreetType"`
	Municipality      string `json:"Municipality"`
	PatrolZone        string `json:"PatrolZone"`
	Division          string `json:"Division"`
	OccType           string `json:"OccType"`
	Ward              string `json:"Ward"`
}

func FetchAndConverPeelEvents(b *base.Base, ctx context.Context, fromDate time.Time, toDate time.Time) (*[]models.CreateTempEventsParams, *utils.CError) {

	rawURL := "https://services.arcgis.com/w0dAT1ctgtKwxvde/arcgis/rest/services/Experience_gdb/FeatureServer/0/query"
	toDateStr := fmt.Sprintf("AND OccDateUTC <= date '%s'", convertToArcgisQueryTime(toDate))
	where := fmt.Sprintf("OccDateUTC >= date '%s' %s", convertToArcgisQueryTime(fromDate), toDateStr)
	endpoint := NewArcgisQuery().DefaultQueries().QWhere(where).BuildWithURL(rawURL)
	arcgisResponse, cerr := fetchArcgis[PeelArcgisAttributes](endpoint)
	if cerr != nil {
		return nil, cerr
	}

	return convertArcgisPeelResponseToEventParams(arcgisResponse), nil
}

func convertArcgisPeelResponseToEventParams(arcgisResponse *ArcgisResponse[PeelArcgisAttributes]) *[]models.CreateTempEventsParams {
	reportsParams := []models.CreateTempEventsParams{}

	for _, arcReport := range arcgisResponse.Features {
		if arcReport.Geometry == nil {
			continue
		}

		x := arcReport.Geometry.X
		y := arcReport.Geometry.Y

		secs := int64(arcReport.Attributes.OccDateUTC / 1000)
		reportsParams = append(reportsParams, models.CreateTempEventsParams{
			OccurAt:      pgtype.Timestamptz{Time: time.Unix(secs, 0).UTC(), Valid: true},
			ExternalID:   arcReport.Attributes.OccurrenceNumber,
			Neighborhood: pgtype.Text{String: removeNeighExtraChars(arcReport.Attributes.StreetName), Valid: true},
			LocationType: pgtype.Text{String: "", Valid: true},
			CrimeType:    arcReport.Attributes.Description,
			Lat:          x,
			Long:         y,
		})
	}

	return &reportsParams
}
