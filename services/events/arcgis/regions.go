package arcgis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
	"townwatch/base"
	"townwatch/models"
	"townwatch/utils"

	"github.com/getsentry/sentry-go"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
)

// york
// "https://services8.arcgis.com/lYI034SQcOoxRCR7/arcgis/rest/services/PublicCrimeDataFME/FeatureServer/0/query?where=0%3D0&objectIds=&time=&geometry=&geometryType=esriGeometryEnvelope&inSR=&spatialRel=esriSpatialRelIntersects&resultType=none&distance=0.0&units=esriSRUnit_Meter&relationParam=&returnGeodetic=false&outFields=*&returnGeometry=true&featureEncoding=esriDefault&multipatchOption=xyFootprint&maxAllowableOffset=&geometryPrecision=&outSR=&defaultSR=&datumTransformation=&applyVCSProjection=false&returnIdsOnly=false&returnUniqueIdsOnly=false&returnCountOnly=false&returnExtentOnly=false&returnQueryGeometry=false&returnDistinctValues=false&cacheHint=false&orderByFields=OBJECTID&groupByFieldsForStatistics=&outStatistics=&having=&resultOffset=&resultRecordCount=&returnZ=false&returnM=false&returnExceededLimitFeatures=true&quantizationParameters=&sqlFormat=none&f=pjson&token="

// peel
// "https://services.arcgis.com/w0dAT1ctgtKwxvde/arcgis/rest/services/ECrimes_gdb/FeatureServer/0/query?where=0%3D0&objectIds=&time=&geometry=&geometryType=esriGeometryEnvelope&inSR=&spatialRel=esriSpatialRelIntersects&resultType=none&distance=0.0&units=esriSRUnit_Meter&relationParam=&returnGeodetic=false&outFields=*&returnGeometry=true&featureEncoding=esriDefault&multipatchOption=xyFootprint&maxAllowableOffset=&geometryPrecision=&outSR=&defaultSR=&datumTransformation=&applyVCSProjection=false&returnIdsOnly=false&returnUniqueIdsOnly=false&returnCountOnly=false&returnExtentOnly=false&returnQueryGeometry=false&returnDistinctValues=false&cacheHint=false&orderByFields=OBJECTID&groupByFieldsForStatistics=&outStatistics=&having=&resultOffset=&resultRecordCount=&returnZ=false&returnM=false&returnExceededLimitFeatures=true&quantizationParameters=&sqlFormat=none&f=pjson&token="

// halton
// "https://services2.arcgis.com/o1LYr96CpFkfsDJS/arcgis/rest/services/Crime_Map/FeatureServer/0/query?where=0%3D0&objectIds=&time=&geometry=&geometryType=esriGeometryEnvelope&inSR=&spatialRel=esriSpatialRelIntersects&resultType=none&distance=0.0&units=esriSRUnit_Meter&relationParam=&returnGeodetic=false&outFields=*&returnGeometry=true&featureEncoding=esriDefault&multipatchOption=xyFootprint&maxAllowableOffset=&geometryPrecision=&outSR=&defaultSR=&datumTransformation=&applyVCSProjection=false&returnIdsOnly=false&returnUniqueIdsOnly=false&returnCountOnly=false&returnExtentOnly=false&returnQueryGeometry=false&returnDistinctValues=false&cacheHint=false&orderByFields=OBJECTID&groupByFieldsForStatistics=&outStatistics=&having=&resultOffset=&resultRecordCount=&returnZ=false&returnM=false&returnExceededLimitFeatures=true&quantizationParameters=&sqlFormat=none&f=pjson&token="

type TorontoArcgisAttributes struct {
	EventUniqueId    string  `json:"EVENT_UNIQUE_ID" validate:"required"`
	OccDateEst       int64   `json:"OCC_DATE_EST" validate:"required"`
	OccDateAgol      int64   `json:"OCC_DATE_AGOL" validate:"required"`
	ReportDateEst    int64   `json:"REPORT_DATE_EST" validate:"required"`
	ReportDateAgol   int64   `json:"REPORT_DATE_AGOL" validate:"required"`
	Division         string  `json:"DIVISION"`
	PremisesType     string  `json:"PREMISES_TYPE"`
	Hour             int16   `json:"HOUR" validate:"required"`
	CrimeType        string  `json:"CRIME_TYPE" validate:"required"`
	Hood158          string  `json:"HOOD_158"`
	Neighbourhood158 string  `json:"NEIGHBOURHOOD_158"`
	Hood140          string  `json:"HOOD_140"`
	Neighbourhood140 string  `json:"NEIGHBOURHOOD_140"`
	Count            int16   `json:"COUNT_"`
	LongWgs84        float64 `json:"LONG_WGS84"`
	LatWgs84         float64 `json:"LAT_WGS84"`
	LocationCategory string  `json:"LOCATION_CATEGORY"`
}

func fetchArcgisToronto(b *base.Base, fromDate time.Time, toDate time.Time) (*ArcgisResponse[TorontoArcgisAttributes], *utils.CError) {
	toDateStr := fmt.Sprintf("AND OCC_DATE_AGOL <= date '%s'", convertToArcgisQueryTime(toDate))
	where := fmt.Sprintf("OCC_DATE_AGOL >= date '%s' %s", convertToArcgisQueryTime(fromDate), toDateStr)
	endpoint := NewArcgisQuery().DefaultQueries().QWhere(url.QueryEscape(where)).BuildWithURL(b.Env.ARCGIS_TORONTO_URL)
	resp, err := http.Get(endpoint)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return nil, &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}
	defer resp.Body.Close()

	var response ArcgisResponse[TorontoArcgisAttributes]
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		eventID := sentry.CaptureException(err)
		return nil, &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}

	if len(response.Features) == 0 {
		sentry.CaptureMessage(fmt.Sprintf("Toronto Arcgis Response: Feature Len is 0 | URL: %s", endpoint))
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	vErr := validate.Struct(response)
	if vErr != nil {
		eventID := sentry.CaptureException(err)
		return nil, &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}
	return &response, nil
}

func convertArcgisTorontoResponseToEventParams(arcgisResponse *ArcgisResponse[TorontoArcgisAttributes]) *[]models.CreateTempEventsParams {
	reportsParams := []models.CreateTempEventsParams{}

	for _, arcReport := range arcgisResponse.Features {
		if arcReport.Attributes.LatWgs84 == 0 || arcReport.Attributes.LongWgs84 == 0 {
			continue
		}

		secs := int64(arcReport.Attributes.OccDateAgol/1000.0) + int64(arcReport.Attributes.Hour*60*60)
		reportsParams = append(reportsParams, models.CreateTempEventsParams{
			OccurAt:      pgtype.Timestamptz{Time: time.Unix(secs, 0).UTC(), Valid: true},
			ExternalID:   arcReport.Attributes.EventUniqueId,
			Neighborhood: pgtype.Text{String: removeNeighExtraChars(arcReport.Attributes.Neighbourhood158), Valid: true},
			LocationType: pgtype.Text{String: arcReport.Attributes.LocationCategory, Valid: true},
			CrimeType:    models.CrimeType(arcReport.Attributes.CrimeType),
			Lat:          arcReport.Attributes.LatWgs84,
			Long:         arcReport.Attributes.LongWgs84,
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
