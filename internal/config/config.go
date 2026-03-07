// Package config provides configuration management for the application.
package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/yuioto/fvti-xsgz-sign/pkg/client"
)

// GetConfigFilePath returns the path to the configuration file.
func GetConfigFilePath(appName string) string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	} else {
		configDir = filepath.Join(configDir, appName)
	}

	const configFile = "Config.toml"
	const dirPerm = 0750

	if err := os.MkdirAll(configDir, dirPerm); err != nil {
		log.Printf("Failed to create config directory: %v", err)
		return configFile
	}

	return filepath.Join(configDir, configFile)
}

// Config represents the application configuration.
type Config struct {
	StudentID string        `toml:"StudentId"`
	Login     Login         `toml:"Login"`
	Task      Task          `toml:"Task"`
	Nofy      string        `toml:"Nofy"`
	Client    client.Config `toml:"Client"`
}

// Task represents the task configuration.
type Task struct {
	Name   string `toml:"Name"`
	ID     string `toml:"Id"`
	SignID string `toml:"SignId"`
}

// Login represents the login configuration.
type Login struct {
	Password      string `toml:"Password"`
	Authorization string `toml:"Authorization"`
}
