package reports

import (
	"net/http"
	"townwatch/base"
	"townwatch/models"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

func LoadRoutes(b *base.Base) {

	b.Engine.GET("/api/private/reports/read", base.SecretRouteMiddleware(b), func(ctx *gin.Context) {
		var params *models.GetPrivateReportDetailsParams
		if err := ctx.BindJSON(&params); err != nil {
			eventID := sentry.CaptureException(err)
			cerr := &base.CError{
				EventID: eventID,
				Error:   err,
			}
			ctx.JSON(http.StatusInternalServerError, cerr)
			return
		}
		report, err := ReadPrivateReport(b, ctx, params)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, report)
	})
	b.Engine.GET("/api/public/reports/read", base.SecretRouteMiddleware(b), func(ctx *gin.Context) {
		var params *GetPublicReportDetailsParams
		if err := ctx.BindJSON(&params); err != nil {
			eventID := sentry.CaptureException(err)
			cerr := &base.CError{
				EventID: eventID,
				Error:   err,
			}
			ctx.JSON(http.StatusInternalServerError, cerr)
			return
		}
		report, err := ReadPublicReport(b, ctx, params)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, report)
	})
}
