package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// LoadConfig loads the configuration from the specified file.
func LoadConfig(configFile string) (Config, error) {
	// Initialize with defaults so that missing fields in TOML get default values
	config := NewDefaultConfig()

	if _, err := os.Stat(configFile); err != nil {
		if !os.IsNotExist(err) {
			return config, fmt.Errorf("check config file: %w", err)
		}

		log.Printf("Config file not found at %s, creating default...", configFile)
		if err := CreateDefaultConfig(configFile); err != nil {
			return config, fmt.Errorf("create default config: %w", err)
		}
		return config, fmt.Errorf("please configure your settings in %s", configFile)
	}

	cfgBytes, err := os.ReadFile(filepath.Clean(configFile))
	if err != nil {
		return config, fmt.Errorf("read config file: %w", err)
	}
	if err := toml.Unmarshal(cfgBytes, &config); err != nil {
		return config, fmt.Errorf("unmarshal config: %w", err)
	}

	return config, nil
}
