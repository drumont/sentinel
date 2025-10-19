package commons

import (
	"errors"
	"log"
	"os"
	"strings"
)

type Config struct {
	PoolsFilePath  string
	OutputFilePath string
}

func (c *Config) valid() error {
	if c.PoolsFilePath == "" {
		return errors.New("pools file path is not set")
	}
	if c.OutputFilePath == "" || !strings.Contains(c.OutputFilePath, ".jsonl") {
		log.Printf("Only .jsonl is supported for output file. Fallback to default output location")
		c.OutputFilePath = "scan.jsonl"
	}
	return nil
}

func LoadConfig() *Config {
	filePath := os.Getenv("POOLS_FILEPATH")
	outputFilePath := os.Getenv("OUTPUT_FILEPATH")
	conf := &Config{PoolsFilePath: filePath, OutputFilePath: outputFilePath}
	if err := conf.valid(); err != nil {
		log.Printf("Error loading agent configuration %v", err)
	}
	return conf
}
