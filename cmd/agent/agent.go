package main

import (
	"log"
	"net/http"
	"sentinel/internal/commons"
	p "sentinel/internal/pools"
	r "sentinel/internal/router"
	"sentinel/internal/scan"
)

var conf = commons.LoadConfig()

func main() {

	pools, err := p.ReadPools(conf.PoolsFilePath)
	if err != nil {
		log.Printf("Error when reading pools %v", err)
	}

	scanner := scan.NewScanner(pools, conf.OutputFilePath)
	scanner.InitScanning()

	router := r.SetupRouter(scanner)
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Agent start failed")
	}

}
