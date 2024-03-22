package base

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const HeaderSecretKeyName string = "Api-Key"

func SecretRouteMiddleware(b *Base) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key := ctx.GetHeader(HeaderSecretKeyName)
		fmt.Println("=====")
		fmt.Println(ctx.Request.Header)
		fmt.Println(key)
		fmt.Println(b.SECRET_API_KEY)
		fmt.Println("=====")
		if key != b.SECRET_API_KEY {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Next()
	}
}
