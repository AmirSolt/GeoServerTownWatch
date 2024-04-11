package events

import (
	"net/http"
	"reflect"
	"time"
	"townwatch/base"
	"townwatch/models"
	"townwatch/utils"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type ScanPointAPIParams struct {
	Lat      float64 `json:"lat"`
	Long     float64 `json:"long"`
	Radius   int32   `json:"radius"`
	Region   string  `json:"region"`
	FromDate string  `json:"from_date"`
	ToDate   string  `json:"to_date"`
	Limit    int32   `json:"limit"`
	Address  string  `json:"address"`
	UserID   string  `json:"user_id"`
}

func LoadRoutes(b *base.Base) {

	b.Engine.POST("/api/events/scan", func(ctx *gin.Context) {
		censorEventsStr := ctx.DefaultQuery("censor_events", "true")
		censorEvents := true
		if censorEventsStr == "false" {
			censorEvents = false
		}

		var params *ScanPointAPIParams
		if err := ctx.BindJSON(&params); err != nil {
			eventID := sentry.CaptureException(err)
			cerr := &utils.CError{
				EventID: eventID,
				Message: "Internal Server Error",
				Error:   err,
			}
			ctx.JSON(http.StatusInternalServerError, cerr)
			return
		}

		// =========================
		// validate and convert

		if params.Limit > int32(b.ScanEventCountLimit) {
			params.Limit = int32(b.ScanEventCountLimit)
		}

		dbParams, errCon := ConvertParams(*params)
		if errCon != nil {
			ctx.JSON(http.StatusInternalServerError, errCon)
			return
		}

		// =========================

		events, err := b.Queries.ScanPoint(ctx, *dbParams)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		// Create a scan for record
		go func() {
			CreateScan(b, ctx, &models.CreateScanParams{
				Radius:      dbParams.Radius,
				FromDate:    dbParams.FromDate,
				ToDate:      dbParams.ToDate,
				EventsCount: int32(len(events)),
				Address:     params.Address,
				UserID:      pgtype.Text{String: params.UserID, Valid: params.UserID != ""},
				Lat:         dbParams.Lat,
				Long:        dbParams.Long,
			})
		}()

		cenEvents := events
		if censorEvents {
			cenEvents = CensorEvents(events)
		}

		ctx.JSON(http.StatusOK, cenEvents)
	})

}

func CreateScan(b *base.Base, ctx *gin.Context, params *models.CreateScanParams) (*models.Scan, *utils.CError) {
	scan, err := b.DB.Queries.CreateScan(ctx, *params)
	if err != nil {
		eventID := sentry.CaptureException(err)
		return nil, &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
	}

	return &scan, nil
}

func CensorEvents(events []models.Event) []models.Event {
	cenEvents := []models.Event{}
	for _, event := range events {
		cenEvents = append(cenEvents, CensorEvent(event))
	}
	return cenEvents
}
func CensorEvent(event models.Event) models.Event {

	uncensoredFields := map[string]bool{"ID": true, "Lat": true, "Long": true}

	eventType := reflect.TypeOf(event)
	eventValue := reflect.ValueOf(&event).Elem()
	for i := 0; i < eventType.NumField(); i++ {
		fieldName := eventType.Field(i).Name
		if !uncensoredFields[fieldName] {
			field := eventValue.FieldByName(fieldName)
			field.Set(reflect.Zero(field.Type()))
		}
	}
	return event
}

func ConvertParams(apiParams ScanPointAPIParams) (*models.ScanPointParams, *utils.CError) {

	fromDate, cerr := convertStrToTime(apiParams.FromDate)
	if cerr != nil {
		return nil, cerr
	}
	toDate, ferr := convertStrToTime(apiParams.ToDate)
	if ferr != nil {
		return nil, ferr
	}

	scanParams := models.ScanPointParams{
		Lat:      apiParams.Lat,
		Long:     apiParams.Long,
		Radius:   apiParams.Radius,
		FromDate: *fromDate,
		ToDate:   *toDate,
		Limit:    apiParams.Limit,
	}

	return &scanParams, nil
}

func convertStrToTime(timeStr string) (*pgtype.Timestamptz, *utils.CError) {
	layout := "Mon, 02 Jan 2006 15:04:05 GMT"
	convDate, err := time.Parse(layout, timeStr)
	if err != nil {
		eventID := sentry.CaptureException(err)
		cerr := &utils.CError{
			EventID: eventID,
			Message: "Internal Server Error",
			Error:   err,
		}
		return nil, cerr
	}
	pgTimestamp := pgtype.Timestamptz{
		Time:  convDate,
		Valid: true,
	}
	return &pgTimestamp, nil
}
