package commons

import (
	"errors"
	"log"
	"os"
)

type Config struct {
	PoolsFilePath  string
	OutputFilePath string
}

func (c *Config) valid() error {
	if c.PoolsFilePath == "" {
		return errors.New("pools file path must be set")
	}
	return nil
}

func LoadConfig() *Config {
	filePath := os.Getenv("POOLS_FILEPATH")
	outputFilePath := os.Getenv("OUTPUT_FILEPATH")
	if outputFilePath == "" {
		outputFilePath = "scan.json"
	}
	conf := &Config{PoolsFilePath: filePath, OutputFilePath: outputFilePath}
	if err := conf.valid(); err != nil {
		log.Fatalf("Error loading agent configuration %v", err)
	}
	return conf
}
