package events

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
	"townwatch/base"
	"townwatch/models"

	"github.com/getsentry/sentry-go"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgtype"
)

type ArcgisResponse struct {
	Features []ArcgisReport `json:"features"`
}

type ArcgisReport struct {
	Attributes ArcgisAttributes `json:"attributes" validate:"required"`
	Geometry   ArcgisGeometry   `json:"geometry" validate:"required"`
}

type ArcgisGeometry struct {
	X float64 `json:"x" validate:"required"`
	Y float64 `json:"y" validate:"required"`
}

type ArcgisAttributes struct {
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

func FetchAndStoreTorontoEvents(b *base.Base, ctx context.Context, fromDate time.Time, toDate time.Time) (int, *base.CError) {
	response, cerr := fetchArcgisToronto(b, fromDate, toDate)
	if cerr != nil {
		return 0, cerr
	}
	eventParams := convertArcgisTorontoResponseToEventParams(response)
	fetchCount := len(*eventParams)

	eventParams = removeEventParamsDuplicates(eventParams)
	err := storeEvents(b, ctx, eventParams)
	if err != nil {
		return 0, err
	}
	return fetchCount, nil
}

func storeEvents(b *base.Base, ctx context.Context, eventParams *[]models.CreateTempEventsParams) *base.CError {
	tx, err := b.DB.Pool.Begin(ctx)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return &base.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}
	defer tx.Rollback(ctx)
	qtx := b.DB.Queries.WithTx(tx)

	if err := qtx.CreateTempEventsTable(ctx); err != nil {
		eventID := sentry.CaptureException(err)
		return &base.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}
	_, errInsert := qtx.CreateTempEvents(ctx, *eventParams)
	if errInsert != nil {
		eventID := sentry.CaptureException(errInsert)
		return &base.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   errInsert,
		}
	}

	if err := qtx.MoveFromTempEventsToEvents(ctx); err != nil {
		eventID := sentry.CaptureException(err)
		return &base.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}

	if err := tx.Commit(ctx); err != nil {
		eventID := sentry.CaptureException(err)
		return &base.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}

	return nil
}

func fetchArcgisToronto(b *base.Base, fromDate time.Time, toDate time.Time) (*ArcgisResponse, *base.CError) {
	toDateStr := fmt.Sprintf("AND OCC_DATE_AGOL <= date '%s'", convertToArcgisQueryTime(toDate))
	where := fmt.Sprintf("OCC_DATE_AGOL >= date '%s' %s", convertToArcgisQueryTime(fromDate), toDateStr)
	endpoint := fmt.Sprintf("%s?where=%s&objectIds=&time=&geometry=&geometryType=esriGeometryEnvelope&inSR=&spatialRel=esriSpatialRelIntersects&resultType=none&distance=0.0&units=esriSRUnit_Meter&relationParam=&returnGeodetic=false&outFields=*&returnGeometry=true&featureEncoding=esriDefault&multipatchOption=xyFootprint&maxAllowableOffset=&geometryPrecision=&outSR=&defaultSR=&datumTransformation=&applyVCSProjection=false&returnIdsOnly=false&returnUniqueIdsOnly=false&returnCountOnly=false&returnExtentOnly=false&returnQueryGeometry=false&returnDistinctValues=false&cacheHint=false&orderByFields=&groupByFieldsForStatistics=&outStatistics=&having=&resultOffset=&resultRecordCount=&returnZ=false&returnM=false&returnExceededLimitFeatures=true&quantizationParameters=&sqlFormat=none&f=pjson&token=", b.Env.ARCGIS_TORONTO_URL, url.QueryEscape(where))
	resp, err := http.Get(endpoint)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return nil, &base.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}
	defer resp.Body.Close()

	var response ArcgisResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		eventID := sentry.CaptureException(err)
		return nil, &base.CError{
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
		return nil, &base.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}
	return &response, nil
}

func convertArcgisTorontoResponseToEventParams(arcgisResponse *ArcgisResponse) *[]models.CreateTempEventsParams {
	reportsParams := []models.CreateTempEventsParams{}

	for _, arcReport := range arcgisResponse.Features {
		secs := int64(arcReport.Attributes.OccDateAgol/1000.0) + int64(arcReport.Attributes.Hour*60*60)
		reportsParams = append(reportsParams, models.CreateTempEventsParams{
			OccurAt:      pgtype.Timestamptz{Time: time.Unix(secs, 0).UTC(), Valid: true},
			ExternalID:   arcReport.Attributes.EventUniqueId,
			Neighborhood: pgtype.Text{String: removeNeighExtraChars(arcReport.Attributes.Neighbourhood158), Valid: true},
			LocationType: pgtype.Text{String: arcReport.Attributes.LocationCategory, Valid: true},
			CrimeType:    models.CrimeType(arcReport.Attributes.CrimeType),
			Region:       string(TorontoRegion),
			Lat:          arcReport.Geometry.X,
			Long:         arcReport.Geometry.Y,
		})
	}

	return &reportsParams
}

func removeEventParamsDuplicates(params *[]models.CreateTempEventsParams) *[]models.CreateTempEventsParams {
	uniqueMap := make(map[string]models.CreateTempEventsParams)
	for _, param := range *params {
		// Only add to the map if ExternalID doesn't exist already
		if _, ok := uniqueMap[param.ExternalID]; !ok {
			uniqueMap[param.ExternalID] = param
		}
	}
	uniqueParams := make([]models.CreateTempEventsParams, 0, len(uniqueMap))
	for _, v := range uniqueMap {
		uniqueParams = append(uniqueParams, v)
	}
	return &uniqueParams
}

// func CreateReports(base *base.Base, reportsParams *[]models.CreateEventsParams) {

// 	count, err := base.DB.Queries.CreateEvents(context.Background(), *reportsParams)
// 	if err != nil {
// 		log.Fatalln("ERROR: bulk insert reports failed:", err, " || count:", count)
// 	}
// }

func convertToArcgisQueryTime(time time.Time) string {
	return fmt.Sprintf("%d-%d-%d %d:%d:%d\n",
		time.Year(),
		time.Month(),
		time.Day(),
		time.Hour(),
		time.Hour(),
		time.Second())
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

// function getInsertValuesStr(rawReports: any[]) {
// 	return rawReports.map(raw => {
// 		const date_at = setDateToSpecificHour(raw["attributes"][DATE_TYPE], raw["attributes"]["HOUR"])

// 		return `
// 			(
// 				${date_at},
// 				${region}_${raw["attributes"]["EVENT_UNIQUE_ID"]},
// 				${removeNeighExtraChars(raw["attributes"]["NEIGHBOURHOOD_158"])},
// 				${raw["attributes"]["LOCATION_CATEGORY"]},
// 				${crimeTypeCleaning(raw["attributes"]["CRIME_TYPE"]) as CrimeType},
// 				${region},
// 				ST_Point(${raw["geometry"]["x"]}, ${raw["geometry"]["y"]}, 3857),
// 				${raw["geometry"]["x"]},
// 				${raw["geometry"]["y"]}
// 			)
// 	`}).join(",")
// }

// return await fastify.prisma.$executeRaw`
// 	INSERT INTO "Report" (
// 			reported_at,
// 			external_src_id,
// 			neighborhood,
// 			location_type,
// 			crime_type,
// 			region,
// 			point,
// 			lat,
// 			long
// 		)

// 		VALUES
// 		${getInsertValuesStr(rawReports)}
// 		ON CONFLICT DO NOTHING;
// `

// }

// function UTCToStr(date: Date): string {
// // Format the date and time
// let year = date.getUTCFullYear();
// let month = ("0" + (date.getUTCMonth() + 1)).slice(-2);
// let day = ("0" + date.getUTCDate()).slice(-2);
// let hours = ("0" + date.getUTCHours()).slice(-2);
// let minutes = ("0" + date.getUTCMinutes()).slice(-2);
// let seconds = ("0" + date.getUTCSeconds()).slice(-2);

// return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
// }

// function removeNeighExtraChars(inputString: string): string {
// let result: string = "";
// let flag: boolean = false;   // True when between "(" and ")"

// for (let char of inputString) {
// 	if (char === "(") {
// 		flag = true;
// 	} else if (char === ")") {
// 		flag = false;
// 	} else if (!flag) {
// 		result += char;
// 	}
// }

// return result;
// }

// function setDateToSpecificHour(epochTime: number, hour: number): Date {
// const date = new Date(epochTime);
// date.setUTCHours(hour, 0, 0, 0);
// return date;
// }

// function crimeTypeCleaning(rawCrimeType: string) {
// return rawCrimeType.toUpperCase().replace(/\s+/g, '_');
// }
