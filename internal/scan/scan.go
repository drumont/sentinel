package scan

import (
	"context"
	"errors"
	"log"
	"os/exec"
	"sentinel/internal/output"
	out "sentinel/internal/output"
	p "sentinel/internal/pools"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	poolsScanState = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "sentinel_pools_scan_state",
			Help: "State of pool scan",
		},
		[]string{"pool", "hostnames", "ports"},
	)
)

type Scan struct {
	Pool *p.Pool
}

type ScanResult struct {
	PoolName string       `json:"pool-name"`
	Output   *out.NmapRun `json:"output"`
}

func NewScan(pool *p.Pool) *Scan {
	return &Scan{Pool: pool}
}

func (s *Scan) Run(channel chan<- ScanResult, ctx context.Context) {
	if s.Pool.ExecuteOnce() {
		s.runOnce(channel)
	} else {
		s.runMany(channel, ctx)
	}
}

func (s *Scan) runOnce(channel chan<- ScanResult) {
	s.executeAndSendResult(channel)
}

func (s *Scan) runMany(channel chan<- ScanResult, ctx context.Context) {
	duration := time.Duration(s.Pool.Interval) * time.Second
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Pool %v scan cancel", s.Pool.Name)
			return
		case <-ticker.C:
			s.executeAndSendResult(channel)
		}
	}
}

func (s *Scan) executeAndSendResult(channel chan<- ScanResult) {
	r, err := execute(s.Pool)
	if err != nil {
		log.Printf("Error during scan of pool %v. err: %v", s.Pool.Name, err)
	}
	if r != nil {
		log.Printf("Pool %v - %v", s.Pool.Name, r.FormatNmapRun())
		scanResult := ScanResult{
			PoolName: s.Pool.Name,
			Output:   r,
		}
		channel <- scanResult
	}
}

func execute(pool *p.Pool) (*output.NmapRun, error) {
	data, err := executeCommand(pool)
	if err != nil {
		return nil, err
	}
	runOutput, err := out.ParseRunOutput(data)
	if err != nil {
		return nil, err
	}
	if !runOutput.IsSuccessfulScan() {
		poolsScanState.WithLabelValues(pool.Name, pool.FormatHosts(), pool.FormatPorts()).Set(0)
		return nil, errors.New("something went wrong during scanning")
	}
	poolsScanState.WithLabelValues(pool.Name, pool.FormatHosts(), pool.FormatPorts()).Set(1)
	return runOutput, nil
}

func executeCommand(pool *p.Pool) ([]byte, error) {
	ports := pool.FormatPorts()

	args := []string{"-Pn", "-p", ports, "-oX", "-"}
	args = append(args, pool.Hosts...)

	cmd := exec.Command("nmap", args...)
	return cmd.CombinedOutput()
}
