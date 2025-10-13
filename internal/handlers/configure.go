package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	p "sentinel/internal/pools"
	"sentinel/internal/scan"

	"github.com/gin-gonic/gin"
)

func Configure(ctx *gin.Context, scanner *scan.Scanner) {

	request := ctx.Request
	var pools []p.Pool

	body, err := io.ReadAll(request.Body)
	if err != nil {
		log.Print(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(body, &pools); err != nil {
		log.Printf("Error parsing request %v", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	scanner.StopScanning()

	ctx.JSON(200, pools)
}
