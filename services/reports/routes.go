package reports

import (
	"fmt"
	"net/http"
	"townwatch/base"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

func LoadRoutes(b *base.Base) {

	b.Engine.GET("/api/reports/:id", func(ctx *gin.Context) {
		reportID, exists := ctx.Params.Get("id")
		if !exists {
			err := fmt.Errorf("report id does not exist")
			eventID := sentry.CaptureException(err)
			cerr := &base.CError{
				EventID: eventID,
				Message: "Report id does not exist",
				Error:   err,
			}
			ctx.JSON(http.StatusInternalServerError, cerr)
			return
		}
		reportDetails, err := GetReportDetails(b, ctx, reportID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, reportDetails)
	})
	b.Engine.GET("/api/reports/:id/events", func(ctx *gin.Context) {
		censorEventsStr := ctx.DefaultQuery("censor_events", "true")
		censorEvents := true
		if censorEventsStr == "false" {
			censorEvents = false
		}

		reportID, exists := ctx.Params.Get("id")
		if !exists {
			err := fmt.Errorf("report id does not exist")
			eventID := sentry.CaptureException(err)
			cerr := &base.CError{
				EventID: eventID,
				Message: "Report id does not exist",
				Error:   err,
			}
			ctx.JSON(http.StatusInternalServerError, cerr)
			return
		}
		events, err := GetEventsByReport(b, ctx, reportID, censorEvents)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, events)
	})
}
