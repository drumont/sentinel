package services

import (
	"sentinel/internal/commons"
	"sentinel/internal/scan"
)

type SentinelServices struct {
	CurrentScanner *scan.Scanner
	Config         *commons.Config
}

func NewSentinelServices(scanner *scan.Scanner, config *commons.Config) *SentinelServices {
	return &SentinelServices{CurrentScanner: scanner, Config: config}
}
