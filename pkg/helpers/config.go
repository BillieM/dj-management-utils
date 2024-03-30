package helpers

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/billiem/seren-management/pkg/projectpath"
)

/*
Config is the main config struct for the application

Built from config.json
*/
type Config struct {
	Development                 bool     `json:"development"`
	TraktorCollectionPath       string   `json:"traktorCollectionPath"`
	BaseDir                     string   `json:"baseDir"`
	DownloadDir                 string   `json:"downloadDir"`
	ExtensionsToConvertToMp3    []string `json:"extensionsToConvertToMp3"`
	ExtensionsToSeparateToStems []string `json:"extensionsToSeparateToStems"`
	CudaEnabled                 bool     `json:"cudaEnabled"`
	DemucsBatchSize             int      `json:"demucsBatchSize"`
	MergeWorkers                int      `json:"mergeWorkers"`
	CleanUpWorkers              int      `json:"cleanUpWorkers"`

	// these are not stored in config.json
	SoundCloudClientID    string `json:"-"`
	SoundCloudSecretToken string `json:"-"`
}

// buildDefaultConfig builds default config values and saves them to config.json
// this is called when the application is first run or when the config file is deleted
func buildDefaultConfig() (*Config, error) {
	cfg := &Config{
		TraktorCollectionPath:       "",
		BaseDir:                     "",
		DownloadDir:                 "",
		ExtensionsToConvertToMp3:    []string{"wav", "aiff", "flac", "ogg", "m4a"},
		ExtensionsToSeparateToStems: []string{"mp3", "wav"},
	}

	cfg.loadEnvConfig()

	err := cfg.SaveConfig()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

/*
LoadCLIConfig loads the config from the given path (if given) or from the default path

No reference to the config is required for the CLI
*/
func LoadCLIConfig(configPath string) (Config, error) {

	// load default config if no config path is given
	if configPath == "" {
		configPath := JoinFilepathToSlash(projectpath.Root, "config.json")
		cfg, err := loadConfig(configPath)
		return *cfg, err
	}

	configPath, err := GetAbsOrWdPath(configPath)

	if err != nil {
		return Config{}, err
	}

	cfg, err := loadConfig(configPath)
	return *cfg, err
}

func LoadGUIConfig() (*Config, error) {

	configPath := JoinFilepathToSlash(projectpath.Root, "config.json")

	_, err := os.Stat(configPath)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// build config from defaults if it doesn't exist
			cfg, err := buildDefaultConfig()
			if err != nil {
				return nil, err
			}
			return cfg, nil
		}
		return nil, err
	}

	return loadConfig(configPath)
}

func loadConfig(configPath string) (*Config, error) {

	// should never hit this but just in case
	_, err := os.Stat(configPath)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	config.loadEnvConfig()

	return &config, nil
}

func (c *Config) SaveConfig() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(JoinFilepathToSlash(projectpath.Root, "config.json"), data, 0644)
	if err != nil {
		return err
	}

	return nil
}

/*
loadEnvConfig loads config values stored in environment variables

Such as API keys/secrets
*/
func (c *Config) loadEnvConfig() {

	_, isDev := os.LookupEnv("DEVELOPMENT")

	c.Development = isDev
	c.SoundCloudClientID = os.Getenv("SOUNDCLOUD_CLIENT_ID")
	c.SoundCloudSecretToken = os.Getenv("SOUNDCLOUD_SECRET_TOKEN")
}

/*
Below config validation functions are called by views to ensure required config values are
present and valid before allowing the user to proceed
*/

func (c *Config) CheckTraktorCollectionPath() (bool, string) {
	if _, err := os.Stat(c.TraktorCollectionPath); os.IsNotExist(err) {
		return false, "Traktor collection path does not exist"
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
