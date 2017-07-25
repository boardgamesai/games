package game

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	TmpDir     string
	PlayerDir  string
	UseSandbox bool
}

const ConfigPath = "config.json"

func Config() (*Configuration, error) {
	file, err := os.Open(ConfigPath)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(file)
	config := Configuration{}
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
