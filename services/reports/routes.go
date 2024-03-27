package reports

import (
	"fmt"
	"net/http"
	"time"
	"townwatch/base"
	"townwatch/models"
	"townwatch/utils"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func LoadRoutes(b *base.Base) {

	b.Engine.GET("/api/reports/:id", func(ctx *gin.Context) {
		reportID, exists := ctx.Params.Get("id")
		if !exists {
			err := fmt.Errorf("report id does not exist")
			eventID := sentry.CaptureException(err)
			cerr := &utils.CError{
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
			cerr := &utils.CError{
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

	if !b.IS_PROD {
		testReportCron(b)
	}
}

func testReportCron(b *base.Base) {
	b.Engine.GET("/api/reports/test", func(ctx *gin.Context) {
		reports, err := b.DB.Queries.CreateGlobalReports(ctx, models.CreateGlobalReportsParams{
			FromDate: pgtype.Timestamptz{
				Time:  time.Now().Add(-time.Duration(24) * time.Hour).UTC(),
				Valid: true,
			},
			ToDate: pgtype.Timestamptz{
				Time:  time.Now().UTC(),
				Valid: true,
			},
			EventsLimit: int32(b.ScanEventCountLimit),
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		fmt.Println("=======")
		fmt.Println("reports", reports)

		ctx.JSON(http.StatusOK, reports)

	})
}
