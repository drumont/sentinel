package scan

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	p "sentinel/internal/pools"
	"sync"
)

type ScannerState int64

const (
	RUNNING ScannerState = iota
	STOPPED
)

type Scanner struct {
	Pools          []p.Pool
	ResultsChannel chan ScanResult
	OutputFilePath string
	State          ScannerState
	ctx            context.Context
	cancel         context.CancelFunc
	writerWg       sync.WaitGroup
}

func NewScanner(pools []p.Pool, fp string) *Scanner {
	channel := make(chan ScanResult, 100)
	ctx, cancel := context.WithCancel(context.Background())
	return &Scanner{Pools: pools, ResultsChannel: channel, writerWg: sync.WaitGroup{}, ctx: ctx, cancel: cancel, OutputFilePath: fp, State: STOPPED}
}

func (s *Scanner) InitScanning() {
	if len(s.Pools) == 0 {
		slog.Info("No pool configured. No scan to process")
	}
	for i := range s.Pools {
		scan := NewScan(&s.Pools[i])
		go scan.Run(s.ResultsChannel, s.ctx)
	}
	s.writerWg.Go(func() { s.writeResult(s.ctx) })
	s.State = RUNNING
}

func (s *Scanner) StopScanning() {
	s.cancel()
	s.writerWg.Wait()
	s.State = STOPPED
	slog.Info("All scans stopped")
}

func (s *Scanner) writeResult(ctx context.Context) {
	f, err := os.OpenFile(s.OutputFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		slog.Error("Output file creation failed with error: "+ err.Error())
	}
	defer f.Close()

	enc := json.NewEncoder(f)

	for {
		select {
		case <-ctx.Done():
			slog.Info("Scanner result writer done")
			return
		case result := <-s.ResultsChannel:
			if err := enc.Encode(result); err != nil {
				slog.Error("Writing scan result failed with error: "+ err.Error())
				continue
			}
		}
	}
}
