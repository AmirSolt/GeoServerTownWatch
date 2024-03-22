package events

import (
	"net/http"
	"reflect"
	"townwatch/base"
	"townwatch/models"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

func LoadRoutes(b *base.Base) {

	b.Engine.POST("/api/events/scan", func(ctx *gin.Context) {
		censorEventsStr := ctx.DefaultQuery("censor_events", "true")
		censorEvents := true
		if censorEventsStr == "false" {
			censorEvents = false
		}

		var params *models.ScanPointParams
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

		if params.ScanEventsCountLimit > int32(b.ScanEventCountLimit) {
			params.ScanEventsCountLimit = int32(b.ScanEventCountLimit)
		}

		eventsRaw, err := b.Queries.ScanPoint(ctx, *params)
		if err != nil {
			eventID := sentry.CaptureException(err)
			cerr := &base.CError{
				EventID: eventID,
				Message: "Internal Server Error",
				Error:   err,
			}
			ctx.JSON(http.StatusInternalServerError, cerr)
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
