package events

import (
	"net/http"
	"townwatch/base"
	"townwatch/models"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

func LoadRoutes(b *base.Base) {

	b.Engine.GET("/api/events/scan", func(ctx *gin.Context) {

		var params *models.ScanPointParams
		if err := ctx.BindJSON(&params); err != nil {
			eventID := sentry.CaptureException(err)
			cerr := &base.CError{
				EventID: eventID,
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
			sentry.CaptureException(err)
			return
		}

		// convert to Report
		events, errconv := base.ConvertArrayInterface[models.Event](eventsRaw)
		if errconv != nil {
			return
		}

		ctx.JSON(http.StatusOK, events)
	})
}
