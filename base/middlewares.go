package base

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const HeaderSecretKeyName string = "Api-Key"

func SecretRouteMiddleware(b *Base) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := ctx.GetHeader(HeaderSecretKeyName)
		if key != b.SECRET_API_KEY {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Next()
	}
}
