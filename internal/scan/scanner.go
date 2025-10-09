package scan

import (
	"encoding/json"
	"log"
	"os"
	p "sentinel/internal/pools"
)

type Scanner struct {
	Pools          []p.Pool
	ResultsChannel chan ScanResult
}

func NewScanner(pools []p.Pool) *Scanner {
	channel := make(chan ScanResult, 100)
	return &Scanner{Pools: pools, ResultsChannel: channel}
}

func (s *Scanner) InitScanning() {
	for _, pool := range s.Pools {
		scan := NewScan(&pool)
		go scan.Run(s.ResultsChannel)
	}
	go s.writeResult()
}

func (s *Scanner) writeResult() {
	f, err := os.OpenFile("scan.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error occurs when creation file. err: %v", err)
	}
	defer f.Close()

	for result := range s.ResultsChannel {
		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			log.Printf("Error marshaling result: %v", err)
			continue
		}
		if _, err := f.Write(jsonData); err != nil {
			log.Printf("Error writing to file: %v", err)
			continue
		}
		if _, err := f.Write([]byte("\n")); err != nil {
			log.Printf("Error writing new line: %v", err)
			continue
		}
		if err := f.Sync(); err != nil {
			log.Printf("Error flushing writer: %v", err)
			continue
		}
		log.Printf("Written result to file for pool %v", result.PoolName)
	}
}
