package helpers

import (
	"encoding/json"
	"os"
)

type Config struct {
	TraktorCollectionPath    string   `json:"traktorCollectionPath"`
	TmpDir                   string   `json:"tmpDir"`
	BaseDir                  string   `json:"baseDir"`
	BaseOutputDir            string   `json:"baseOutputDir"`
	ExtensionsToConvertToMp3 []string `json:"extensionsToConvertToMp3"`
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile("../config.json")
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

func (c *Config) CheckTraktorCollectionPath() (bool, string) {
	if _, err := os.Stat(c.TraktorCollectionPath); os.IsNotExist(err) {
		return false, "Traktor collection path does not exist"
	}
	return true, ""
}

func (c *Config) CheckTmpDir() (bool, string) {
	fi, err := os.Stat(c.TmpDir)
	if err != nil {
		return false, "Temporary directory does not exist"
	}
	if !fi.IsDir() {
		return false, "Temporary directory is not a directory"
	}
	return true, ""
}

func (c *Config) CheckBaseDir() (bool, string) {
	fi, err := os.Stat(c.BaseDir)
	if err != nil {
		return false, "Base directory does not exist"
	}
	if !fi.IsDir() {
		return false, "Base directory is not a directory"
	}
	return true, ""
}

func (c *Config) CheckBaseOutputDir() (bool, string) {
	fi, err := os.Stat(c.BaseOutputDir)
	if err != nil {
		return false, "Base output directory does not exist"
	}
	if !fi.IsDir() {
		return false, "Base output directory is not a directory"
	}
	return true, ""
}
