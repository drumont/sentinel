package router

import (
	"net/http"
	"sentinel/internal/handlers"
	"sentinel/internal/scan"

	"github.com/gin-gonic/gin"
)

func SetupRouter(scanner *scan.Scanner) http.Handler {
	r := gin.Default()

	r.Handle("GET", "/health", func(ctx *gin.Context) {
		handlers.HealthCheck(ctx)
	})

	r.Handle("POST", "/configure", func(ctx *gin.Context) {
		handlers.Configure(ctx, scanner)
	})

	return r
}
