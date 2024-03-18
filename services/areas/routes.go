package areas

import (
	"net/http"
	"townwatch/base"
	"townwatch/models"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

func LoadRoutes(b *base.Base) {

	b.Engine.POST("/api/areas/create", base.SecretRouteMiddleware(b), func(ctx *gin.Context) {
		var params *models.CreateAreaParams
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
		err := CreateArea(b, ctx, params)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.Status(http.StatusOK)
	})

	b.Engine.GET("/api/areas/read", base.SecretRouteMiddleware(b), func(ctx *gin.Context) {
		var params *models.GetAreaParams
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
		area, err := ReadArea(b, ctx, params)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, area)
	})

	b.Engine.PATCH("/api/areas/update", base.SecretRouteMiddleware(b), func(ctx *gin.Context) {
		var params *models.UpdateAreaParams
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
		err := UpdateArea(b, ctx, params)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.Status(http.StatusOK)
	})

	b.Engine.DELETE("/api/areas/delete", base.SecretRouteMiddleware(b), func(ctx *gin.Context) {
		var params *models.DeleteAreaParams
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
		err := DeleteArea(b, ctx, params)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.Status(http.StatusOK)
	})

	b.Engine.GET("/api/areas/read-by-user", base.SecretRouteMiddleware(b), func(ctx *gin.Context) {
		var params *GetAreasByUserParams
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
		area, err := ReadAreasByUser(b, ctx, params)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, area)
	})
}
