package events

import (
	"net/http"
	"reflect"
	"time"
	"townwatch/base"
	"townwatch/models"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type ScanPointAPIParams struct {
	Lat                  float64 `json:"lat"`
	Long                 float64 `json:"long"`
	Radius               float64 `json:"radius"`
	Region               string  `json:"region"`
	FromDate             string  `json:"from_date"`
	ToDate               string  `json:"to_date"`
	ScanEventsCountLimit int32   `json:"scan_events_count_limit"`
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
			cerr := &base.CError{
				EventID: eventID,
				Message: "Internal Server Error",
				Error:   err,
			}
			ctx.JSON(http.StatusInternalServerError, cerr)
			return
		}

		// =========================
		// validate and convert

		if params.ScanEventsCountLimit > int32(b.ScanEventCountLimit) {
			params.ScanEventsCountLimit = int32(b.ScanEventCountLimit)
		}

		dbParams, errCon := ConvertParams(*params)
		if errCon != nil {
			ctx.JSON(http.StatusInternalServerError, errCon)
			return
		}

		// =========================

		eventsRaw, err := b.Queries.ScanPoint(ctx, *dbParams)
		if err != nil {
			eventID := sentry.CaptureException(err)
			cerr := &base.CError{
				EventID: eventID,
				Message: "Internal Server Error",
				Error:   err,
			}
			ctx.JSON(http.StatusInternalServerError, *cerr)
			return
		}

		// convert to Report
		events, errconv := base.ConvertArrayInterface[models.Event](eventsRaw)
		if errconv != nil {
			ctx.JSON(http.StatusInternalServerError, errconv)
			return
		}

		var cenEvents []models.Event
		if censorEvents {
			cenEvents = CensorEvents(events)
		}

		ctx.JSON(http.StatusOK, cenEvents)
	})

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
		if uncensoredFields[fieldName] {
			field := eventValue.FieldByName(fieldName)
			field.Set(reflect.Zero(field.Type()))
		}
	}
	return event
}

func ConvertParams(apiParams ScanPointAPIParams) (*models.ScanPointParams, *base.CError) {
	scanParams := models.ScanPointParams{}

	apiFields := reflect.ValueOf(&apiParams).Elem()
	scanFields := reflect.ValueOf(&scanParams).Elem()

	for i := 0; i < apiFields.NumField(); i++ {
		apiField := apiFields.Field(i)
		scanField := scanFields.Field(i)

		if scanField.Type() == reflect.TypeOf(pgtype.Timestamptz{}) {
			// ================
			// Convert time from string
			layout := "Mon, 02 Jan 2006 15:04:05 GMT"
			convDate, err := time.Parse(layout, apiField.Interface().(string))
			if err != nil {
				eventID := sentry.CaptureException(err)
				cerr := &base.CError{
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
			scanField.Set(reflect.ValueOf(&pgTimestamp).Elem())
			continue
			// =================
		} else {
			scanField.Set(apiField)
		}

	}

	return &scanParams, nil
}
