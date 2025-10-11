package main

import (
	"log"
	"net/http"
	"sentinel/internal/commons"
	p "sentinel/internal/pools"
	"sentinel/internal/scan"

	"github.com/gin-gonic/gin"
)

var conf = commons.LoadConfig()

func main() {

	router := gin.Default()
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	pools, err := p.ReadPools(conf.PoolsFilePath)
	if err != nil {
		log.Printf("Error when reading pools %v", err)
	}

	scanner := scan.NewScanner(pools)
	scanner.InitScanning()

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Agent start failed")
	}

}
