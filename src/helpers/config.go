package helpers

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	TraktorCollectionPath string `json:"traktorCollectionPath"`
}

func LoadConfig() (*Config, error) {
	data, err := ioutil.ReadFile("../config.json")
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func (c *Config) SaveConfig() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile("../config.json", data, 0644)
	if err != nil {
		return err
	}

	return nil
}
