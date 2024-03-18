package events

import (
	"net/http"
	"time"
	"townwatch/base"
	"townwatch/models"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

func LoadRoutes(b *base.Base) {

	if !b.IS_PROD {
		LoadTestRoutes(b)
	}

	b.Engine.GET("/api/events/scan", func(ctx *gin.Context) {

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

		ctx.JSON(http.StatusOK, events)
	})

}

func LoadTestRoutes(b *base.Base) {
	b.Engine.GET("/api/test/events/fetch", func(ctx *gin.Context) {

		count, err := FetchAndStoreTorontoEvents(b, ctx, time.Now().Add(-time.Duration(10)*time.Hour).UTC(), time.Now().UTC())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, map[string]any{"count": count})
	})
}
