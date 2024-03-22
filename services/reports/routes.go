package reports

import (
	"net/http"
	"townwatch/base"
	"townwatch/models"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

func LoadRoutes(b *base.Base) {

	b.Engine.GET("/api/reports/read", func(ctx *gin.Context) {
		censorEventsStr := ctx.DefaultQuery("censor_events", "true")
		censorEvents := true
		if censorEventsStr == "false" {
			censorEvents = false
		}

		var params *models.GetReportDetailsParams
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
		report, err := ReadPrivateReport(b, ctx, params, censorEvents)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		ctx.JSON(http.StatusOK, report)
	})
}
