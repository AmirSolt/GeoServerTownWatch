package events

import (
	"net/http"
	"time"
	"townwatch/base"

	"github.com/gin-gonic/gin"
)

func LoadRoutes(b *base.Base) {

	if !b.IS_PROD {
		loadTestRoutes(b)
	}
}

func loadTestRoutes(b *base.Base) {
	b.Engine.POST("/join/signin", func(ctx *gin.Context) {

		err := FetchAndStoreTorontoEvents(b, ctx, time.Now().Add(-time.Duration(10)*time.Hour), time.Now())
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err)
			return
		}

		ctx.Status(http.StatusOK)
		return
	})
}
