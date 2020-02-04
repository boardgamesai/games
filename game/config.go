package game

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	TmpDir string
}

const ConfigPath = "config.json"

func Config() (*Configuration, error) {
	file, err := os.Open(ConfigPath)
	if os.IsNotExist(err) {
		config := Configuration{
			TmpDir: "/tmp",
		}
		return &config, nil
	} else if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(file)
	config := Configuration{}
	err = decoder.Decode(&config)
	return &config, err
}
