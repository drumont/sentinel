package handlers

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	p "sentinel/internal/pools"
	"sentinel/internal/scan"
	s "sentinel/internal/services"

	"github.com/gin-gonic/gin"
)

func Stop(ctx *gin.Context, services *s.SentinelServices) {
	if services.CurrentScanner.State == scan.RUNNING {
		services.CurrentScanner.StopScanning()
	}
	ctx.Status(http.StatusAccepted)
}

func Configure(ctx *gin.Context, services *s.SentinelServices) {
	request := ctx.Request
	var pools []p.Pool

	body, err := io.ReadAll(request.Body)
	if err != nil {
		slog.Error("Reading request body failed with error: "+ err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(body, &pools); err != nil {
		slog.Error("Parsing request body failed with error: "+ err.Error())
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if services.CurrentScanner.State == scan.RUNNING {
		services.CurrentScanner.StopScanning()
	}

	services.CurrentScanner = scan.NewScanner(pools, services.Config.OutputFilePath)
	services.CurrentScanner.InitScanning()

	ctx.Status(http.StatusAccepted)
}
