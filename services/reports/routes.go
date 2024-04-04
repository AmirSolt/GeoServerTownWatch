package reports

import (
	"fmt"
	"net/http"
	"townwatch/base"
	"townwatch/models"
	"townwatch/utils"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

type GetReportsByUserParams struct {
	UserID string `json:"user_id"`
}

func LoadRoutes(b *base.Base) {

	b.Engine.POST("/api/reports/user", func(ctx *gin.Context) {

		var params *GetReportsByUserParams
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

		reports, err := GetReportsByUser(b, ctx, params.UserID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, reports)
	})

	b.Engine.POST("/api/reports/area", func(ctx *gin.Context) {
		var params *models.GetReportsByAreaParams
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

		reports, err := GetReportsByArea(b, ctx, params)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, reports)
	})

	b.Engine.GET("/api/reports/:id", func(ctx *gin.Context) {
		reportID, exists := ctx.Params.Get("id")

		if !exists {
			err := fmt.Errorf("report id does not exist in URL")
			eventID := sentry.CaptureException(err)
			cerr := &utils.CError{
				EventID: eventID,
				Message: "This url is broken",
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
		cerr := CreateGlobalReports(b)
		if cerr != nil {
			ctx.JSON(http.StatusInternalServerError, cerr)
		}
		ctx.JSON(http.StatusOK, reports)
	})
}
