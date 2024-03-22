package base

// Webframework, handles mostly routes and requests

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (base *Base) loadEngine() {
	gin.SetMode(gin.DebugMode)
	// gin.DisableConsoleColor()
	engine := gin.Default()

	engine.GET("/ping", func(ctx *gin.Context) { ctx.String(http.StatusOK, "pong") })
	engine.Use(SecretRouteMiddleware(base))

	base.Engine = engine
}
