package scan

import (
	"context"
	"encoding/json"
	"log"
	"os"
	p "sentinel/internal/pools"
	"sync"
)

type Scanner struct {
	Pools          []p.Pool
	ResultsChannel chan ScanResult
	wg             sync.WaitGroup
	Ctx            context.Context
	Cancel         context.CancelFunc
}

func NewScanner(pools []p.Pool) *Scanner {
	channel := make(chan ScanResult, 100)
	ctx, cancel := context.WithCancel(context.Background())
	return &Scanner{Pools: pools, ResultsChannel: channel, Ctx: ctx, Cancel: cancel}
}

func (s *Scanner) InitScanning() {
	for i := range s.Pools {
		s.wg.Add(1)
		s.wg.Done()
		scan := NewScan(&s.Pools[i])
		go scan.Run(s.ResultsChannel, s.Ctx)
	}
	go s.writeResult()
}

func (s *Scanner) StopScanning() {
	s.Cancel()
	s.wg.Wait()
	log.Printf("All scanning operations stopped")
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
