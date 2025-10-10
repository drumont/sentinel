package pools

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Pool struct {
	Name     string   `json:"name"`
	Hosts    []string `json:"hosts"`
	Ports    []string `json:"ports"`
	Interval int      `json:"interval"`
}

var SUPPORTED_EXTENSION = ".json"

func (p *Pool) OneTineScan() bool {
	if p.Interval == 0 {
		return true
	}
	return false
}

func (p *Pool) FormatHosts() string {
	return strings.Join(p.Hosts, " ")
}

func (p *Pool) FormatPorts() string {
	return strings.Join(p.Ports, ",")
}

func ReadPools(path string) ([]Pool, error) {
	pools, err := ReadPoolsFromFile(path)
	if err != nil {
		return nil, err
	}
	return pools, err
}

func ReadPoolsFromFile(path string) ([]Pool, error) {
	err := verifyExtension(path)
	if err != nil {
		return nil, err
	}

	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	data, err := io.ReadAll(jsonFile)
	var pools []Pool

	err = json.Unmarshal(data, &pools)
	if err != nil {
		return nil, err
	}

	return pools, nil
}

func verifyExtension(path string) error {
	extension := filepath.Ext(path)

	if extension != SUPPORTED_EXTENSION {
		return errors.New("file extension not supported. Only json are supported now")
	}

	return nil
}
