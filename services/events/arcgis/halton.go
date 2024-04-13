package arcgis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"townwatch/base"
	"townwatch/models"
	"townwatch/utils"

	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgx/v5/pgtype"
)

type HaltonArcgisAttributes struct {
	ObjectID    int    `json:"OBJECTID"`
	CaseNo      string `json:"CASE_NO"`
	Date        int64  `json:"DATE"`
	Description string `json:"DESCRIPTION"`
	Location    string `json:"LOCATION"`
	City        string `json:"CITY"`
	Latitude    string `json:"Latitude"`
	Longitude   string `json:"Longitude"`
	GlobalID    string `json:"GlobalID"`
}

func FetchAndConverHaltonEvents(b *base.Base, ctx context.Context, fromDate time.Time, toDate time.Time) (*[]models.CreateTempEventsParams, *utils.CError) {

	rawURL := "https://services2.arcgis.com/o1LYr96CpFkfsDJS/arcgis/rest/services/Crime_Map/FeatureServer/0/query"
	toDateStr := fmt.Sprintf("AND DATE <= date '%s'", convertToArcgisQueryTime(toDate))
	where := fmt.Sprintf("DATE >= date '%s' %s", convertToArcgisQueryTime(fromDate), toDateStr)
	endpoint := NewArcgisQuery().DefaultQueries().QWhere(where).BuildWithURL(rawURL)
	arcgisResponse, cerr := fetchArcgis[HaltonArcgisAttributes](endpoint)
	if cerr != nil {
		return nil, cerr
	}

	return convertArcgisHaltonResponseToEventParams(arcgisResponse), nil
}

func convertArcgisHaltonResponseToEventParams(arcgisResponse *ArcgisResponse[HaltonArcgisAttributes]) *[]models.CreateTempEventsParams {
	reportsParams := []models.CreateTempEventsParams{}

	for _, arcReport := range arcgisResponse.Features {
		if arcReport.Geometry == nil {
			continue
		}
		x := arcReport.Geometry.X
		y := arcReport.Geometry.Y

		detailParams := EventDetailsParams{
			"Case Number":  arcReport.Attributes.CaseNo,
			"Description":  arcReport.Attributes.Description,
			"Neighborhood": arcReport.Attributes.Location,
			"City":         arcReport.Attributes.City,
		}
		jsonString, err := json.Marshal(utils.EventDetailsStringCleaner(detailParams))
		if err != nil {
			sentry.CaptureException(err)
			continue
		}

		secs := int64(arcReport.Attributes.Date / 1000)
		reportsParams = append(reportsParams, models.CreateTempEventsParams{
			ExternalID: arcReport.Attributes.GlobalID,
			OccurAt:    pgtype.Timestamptz{Time: time.Unix(secs, 0).UTC(), Valid: true},
			Lat:        y,
			Long:       x,
			Details:    []byte(jsonString),
		})
	}

	return &reportsParams
}
