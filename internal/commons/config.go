package commons

import (
	"errors"
	"log"
	"os"
)

type Config struct {
	PoolsFilePath string
}

func (c *Config) valid() error {
	if c.PoolsFilePath == "" {
		return errors.New("pools file path must be set")
	}
	return nil
}

func LoadConfig() *Config {
	filePath := os.Getenv("POOLS_FILEPATH")
	conf := &Config{PoolsFilePath: filePath}
	if err := conf.valid(); err != nil {
		log.Fatalf("Error loading agent configuration %v", err)
	}
	return conf
}
