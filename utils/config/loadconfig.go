package config

import (
	"fmt"
	"log"
	"os"

	"github.com/pelletier/go-toml/v2"
)

func LoadConfig(configFile string) (Config, error) {
	var config Config

	if _, err := os.Stat(configFile); err != nil {
		if os.IsNotExist(err) {
			log.Println("Error:", err)
			log.Println("Creating default config file")
			CreateDefaultConfig(configFile)
			fmt.Println("Example: 'Bearer {token}' --> 'Bearer you_token'")
			log.Fatalln("Please write your config file at", configFile)
		} else {
			return config, fmt.Errorf("failed to check config file: %w", err)
		}
	} else {
		cfg, err := os.ReadFile(configFile)
		if err != nil {
			return config, fmt.Errorf("failed to read file: %w", err)
		}
		if err := toml.Unmarshal(cfg, &config); err != nil {
			return config, fmt.Errorf("failed to unmarshal TOML file: %w", err)
		}
		//log.Println("Config file read successfully")
	}

	return config, nil
}
