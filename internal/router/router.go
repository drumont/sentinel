package router

import (
	"net/http"
	"sentinel/internal/handlers"
	s "sentinel/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter(services *s.SentinelServices) http.Handler {
	r := gin.Default()

	r.Handle(http.MethodGet, "/health", func(ctx *gin.Context) {
		handlers.HealthCheck(ctx)
	})

	r.Handle(http.MethodPost, "/configure", func(ctx *gin.Context) {
		handlers.Configure(ctx, services)
	})

	r.Handle(http.MethodGet, "/stop", func(ctx *gin.Context) {
		handlers.Stop(ctx, services)
	})

	return r
}
