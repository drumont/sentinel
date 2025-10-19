package main

import (
	"log"
	"net/http"
	"sentinel/internal/commons"
	p "sentinel/internal/pools"
	r "sentinel/internal/router"
	"sentinel/internal/scan"
	"sentinel/internal/services"
)

var conf = commons.LoadConfig()

func main() {

	scanner := scan.NewScanner(make([]p.Pool, 0), conf.OutputFilePath)
	sentinelServices := services.NewSentinelServices(scanner, conf)

	if conf.PoolsFilePath != "" {
		pools, err := p.ReadPools(conf.PoolsFilePath)
		if err != nil {
			log.Printf("Error when reading pools %v", err)
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
		log.Fatal("Agent start failed")
	}

}
