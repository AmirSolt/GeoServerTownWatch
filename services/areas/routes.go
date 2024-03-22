package areas

import (
	"errors"
	"net/http"
	"townwatch/base"
	"townwatch/models"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

func LoadRoutes(b *base.Base) {

	b.Engine.POST("/api/areas/create", func(ctx *gin.Context) {
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

		testErr := errors.New("test error")
		eventID := sentry.CaptureException(testErr)
		ctx.JSON(http.StatusInternalServerError, &base.CError{Message: testErr.Error(), EventID: eventID, Error: testErr})
		return

		area, err := CreateArea(b, ctx, params)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, area)
	})

	b.Engine.GET("/api/areas/read", func(ctx *gin.Context) {
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

	b.Engine.PATCH("/api/areas/update", func(ctx *gin.Context) {
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
		area, err := UpdateArea(b, ctx, params)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, area)
	})

	b.Engine.DELETE("/api/areas/delete", func(ctx *gin.Context) {
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
		area, err := DeleteArea(b, ctx, params)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, area)
	})

	b.Engine.GET("/api/areas/user", func(ctx *gin.Context) {
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
