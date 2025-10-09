package main

import (
	"log"
	"net/http"
	p "sentinel/internal/pools"
	"sentinel/internal/scan"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	pools, err := p.ReadPools("/Users/drumont/Developer/drumont/sentinel/scripts/data.json")
	if err != nil {
		log.Print("Error when reading pools")
	}

	scanner := scan.NewScanner(pools)
	scanner.InitScanning()

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Agent start failed")
	}

}
