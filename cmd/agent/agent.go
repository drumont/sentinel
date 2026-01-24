package main

import (
	"net/http"
	"sentinel/internal/commons"
	p "sentinel/internal/pools"
	r "sentinel/internal/router"
	"sentinel/internal/scan"
	"sentinel/internal/services"
	"log/slog"
	"os"
)


func main() {

	var logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	var conf = commons.LoadConfig()

	scanner := scan.NewScanner(make([]p.Pool, 0), conf.OutputFilePath)
	sentinelServices := services.NewSentinelServices(scanner, conf)

	if conf.PoolsFilePath != "" {
		pools, err := p.ReadPools(conf.PoolsFilePath)
		if err != nil {
			slog.Error("Error when reading pools: "+ err.Error())
		}
		sentinelServices.CurrentScanner.Pools = pools
		sentinelServices.CurrentScanner.InitScanning()
	}

	router := r.SetupRouter(sentinelServices)
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		slog.Error("Agent start failed with error: "+ err.Error())
		os.Exit(1)
	}

}
