package areas

import (
	"net/http"
	"townwatch/base"
	"townwatch/models"
	"townwatch/utils"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

func LoadRoutes(b *base.Base) {

	b.Engine.POST("/api/areas/create", func(ctx *gin.Context) {

		var params *models.CreateAreaParams
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

		area, err := CreateArea(b, ctx, params)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, area)
	})

	b.Engine.POST("/api/areas/read", func(ctx *gin.Context) {
		var params *models.GetAreaParams
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
		area, err := ReadArea(b, ctx, params)
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
			cerr := &utils.CError{
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

	b.Engine.POST("/api/areas/user", func(ctx *gin.Context) {
		var params *GetAreasByUserParams
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
		areas, err := ReadAreasByUser(b, ctx, params)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}
		ctx.JSON(http.StatusOK, areas)
	})
}
