package arcgis

import (
	"context"
	"fmt"
	"strings"
	"time"
	"townwatch/base"
	"townwatch/models"
	"townwatch/utils"

	"github.com/jackc/pgx/v5/pgtype"
)

type TorontoArcgisAttributes struct {
	EventUniqueId    string  `json:"EVENT_UNIQUE_ID"`
	OccDateEst       int64   `json:"OCC_DATE_EST"`
	OccDateAgol      int64   `json:"OCC_DATE_AGOL"`
	ReportDateEst    int64   `json:"REPORT_DATE_EST"`
	ReportDateAgol   int64   `json:"REPORT_DATE_AGOL"`
	Division         string  `json:"DIVISION"`
	PremisesType     string  `json:"PREMISES_TYPE"`
	Hour             int16   `json:"HOUR"`
	CrimeType        string  `json:"CRIME_TYPE"`
	Hood158          string  `json:"HOOD_158"`
	Neighbourhood158 string  `json:"NEIGHBOURHOOD_158"`
	Hood140          string  `json:"HOOD_140"`
	Neighbourhood140 string  `json:"NEIGHBOURHOOD_140"`
	Count            int16   `json:"COUNT_"`
	LongWgs84        float64 `json:"LONG_WGS84"`
	LatWgs84         float64 `json:"LAT_WGS84"`
	LocationCategory string  `json:"LOCATION_CATEGORY"`
}

func FetchAndConverTorontoEvents(b *base.Base, ctx context.Context, fromDate time.Time, toDate time.Time) (*[]models.CreateTempEventsParams, *utils.CError) {

	rawURL := "https://services.arcgis.com/S9th0jAJ7bqgIRjw/ArcGIS/rest/services/YTD_CRIME_WM/FeatureServer/0/query"
	toDateStr := fmt.Sprintf("AND OCC_DATE_AGOL <= date '%s'", convertToArcgisQueryTime(toDate))
	where := fmt.Sprintf("OCC_DATE_AGOL >= date '%s' %s", convertToArcgisQueryTime(fromDate), toDateStr)
	endpoint := NewArcgisQuery().DefaultQueries().QWhere(where).BuildWithURL(rawURL)
	arcgisResponse, cerr := fetchArcgis[TorontoArcgisAttributes](endpoint)
	if cerr != nil {
		return nil, cerr
	}

	return convertArcgisTorontoResponseToEventParams(arcgisResponse), nil
}

func convertArcgisTorontoResponseToEventParams(arcgisResponse *ArcgisResponse[TorontoArcgisAttributes]) *[]models.CreateTempEventsParams {
	reportsParams := []models.CreateTempEventsParams{}

	for _, arcReport := range arcgisResponse.Features {
		if arcReport.Geometry == nil {
			continue
		}

		x := arcReport.Geometry.X
		y := arcReport.Geometry.Y

		secs := int64(arcReport.Attributes.OccDateAgol / 1000)
		tempTime := time.Unix(secs, 0)
		tempTime = tempTime.Add(time.Hour * time.Duration(int(arcReport.Attributes.Hour)-tempTime.Hour()))
		reportsParams = append(reportsParams, models.CreateTempEventsParams{
			OccurAt:      pgtype.Timestamptz{Time: tempTime.UTC(), Valid: true},
			ExternalID:   arcReport.Attributes.EventUniqueId,
			Neighborhood: pgtype.Text{String: removeNeighExtraChars(arcReport.Attributes.Neighbourhood158), Valid: true},
			LocationType: pgtype.Text{String: arcReport.Attributes.LocationCategory, Valid: true},
			CrimeType:    arcReport.Attributes.CrimeType,
			Lat:          y,
			Long:         x,
		})
	}

	return &reportsParams
}

func removeNeighExtraChars(inputString string) string {
	var result strings.Builder
	flag := false // True when between "(" and ")"

	for _, char := range inputString {
		if char == '(' {
			flag = true
		} else if char == ')' {
			flag = false
		} else if !flag {
			// Write the character to the result if we're not between "(" and ")"
			result.WriteRune(char)
		}
	}

	return result.String()
}
