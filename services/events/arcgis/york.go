package arcgis

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"townwatch/base"
	"townwatch/models"
	"townwatch/utils"

	"github.com/getsentry/sentry-go"
	"github.com/jackc/pgx/v5/pgtype"
)

type YorkArcgisAttributes struct {
	UniqueIdentifier string  `json:"UniqueIdentifier"`
	OccDate          int64   `json:"occ_date"`
	OccTime          int     `json:"occ_time"`
	CaseTypePubtrans string  `json:"case_type_pubtrans"`
	CaseCategory1    string  `json:"case_category1"`
	LocationCode     string  `json:"LocationCode"`
	District         string  `json:"district"`
	Municipality     string  `json:"municipality"`
	SpecialGrouping  string  `json:"Special_grouping"`
	XCoordinate      string  `json:"x_coordinate"`
	YCoordinate      string  `json:"y_coordinate"`
	WeekDay          int     `json:"week_day"`
	CrimePerPop      float64 `json:"crime_per_pop"`
	ObjectID         int     `json:"OBJECTID"`
	CrimePrevention  string  `json:"crimeprevention"`
	Shooting         string  `json:"Shooting"`
	MediaRelease     string  `json:"MediaRelease"`
	MediaReleaseFlag string  `json:"mediareleaseflag"`
	DemsPublicURL    string  `json:"dems_public_url"`
	DemsFlag         string  `json:"dems_flag"`
	OccID            string  `json:"occ_id"`
	HateCrime        string  `json:"hate_crime"`
	CaseStatus       string  `json:"case_status"`
	OccType          string  `json:"occ_type"`
	VehicleMake      string  `json:"Vehicle_make"`
	VehicleModel     string  `json:"Vehicle_model"`
	VehicleStyle     string  `json:"Vehicle_Style"`
	VehicleColour    string  `json:"Vehicle_colour"`
	ReportDate       int64   `json:"rep_date"`
}

func FetchAndConverYorkEvents(b *base.Base, ctx context.Context, fromDate time.Time, toDate time.Time) (*[]models.CreateTempEventsParams, *utils.CError) {

	rawURL := "https://services8.arcgis.com/lYI034SQcOoxRCR7/arcgis/rest/services/PublicCrimeDataFME/FeatureServer/0/query"
	toDateStr := fmt.Sprintf("AND occ_date <= date '%s'", convertToArcgisQueryTime(toDate))
	where := fmt.Sprintf("occ_date >= date '%s' %s", convertToArcgisQueryTime(fromDate), toDateStr)
	endpoint := NewArcgisQuery().DefaultQueries().QWhere(where).BuildWithURL(rawURL)
	arcgisResponse, cerr := fetchArcgis[YorkArcgisAttributes](endpoint)
	if cerr != nil {
		return nil, cerr
	}

	return convertArcgisYorkResponseToEventParams(arcgisResponse), nil
}

func convertArcgisYorkResponseToEventParams(arcgisResponse *ArcgisResponse[YorkArcgisAttributes]) *[]models.CreateTempEventsParams {
	reportsParams := []models.CreateTempEventsParams{}

	for _, arcReport := range arcgisResponse.Features {
		if arcReport.Geometry == nil {
			continue
		}

		x := arcReport.Geometry.X
		y := arcReport.Geometry.Y

		detailParams := EventDetailsParams{
			"Unique ID":         arcReport.Attributes.UniqueIdentifier,
			"Description":       arcReport.Attributes.CaseTypePubtrans,
			"Occurrence Type":   arcReport.Attributes.OccType,
			"Location Category": arcReport.Attributes.CaseCategory1,
			"Premises Type":     arcReport.Attributes.LocationCode,
			"City":              arcReport.Attributes.Municipality,
			"Week Day":          strconv.Itoa(arcReport.Attributes.WeekDay),
			"Status":            arcReport.Attributes.CaseStatus,
		}

		if arcReport.Attributes.VehicleMake != "" {
			detailParams["Vehicle Make"] = arcReport.Attributes.VehicleMake
		}
		if arcReport.Attributes.VehicleModel != "" {
			detailParams["Vehicle Model"] = arcReport.Attributes.VehicleModel
		}

		jsonString, err := json.Marshal(utils.EventDetailsStringCleaner(detailParams))
		if err != nil {
			sentry.CaptureException(err)
			continue
		}

		secs := int64(arcReport.Attributes.OccDate / 1000)
		reportsParams = append(reportsParams, models.CreateTempEventsParams{
			OccurAt:    pgtype.Timestamptz{Time: time.Unix(secs, 0).UTC(), Valid: true},
			ExternalID: arcReport.Attributes.UniqueIdentifier,
			Details:    []byte(jsonString),
			Lat:        y,
			Long:       x,
		})
	}

	return &reportsParams
}
