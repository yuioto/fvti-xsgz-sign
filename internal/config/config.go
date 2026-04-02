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

	const configFile = "config.kdl"
	const dirPerm = 0750

	if err := os.MkdirAll(configDir, dirPerm); err != nil {
		log.Printf("Failed to create config directory: %v", err)
		return configFile
	}

	return filepath.Join(configDir, configFile)
}

// Config represents the application configuration.
type Config struct {
	Login  Login         `kdl:"login"`
	Task   Task          `kdl:"task"`
	Notify Notify        `kdl:"notify"`
	Client client.Config `kdl:"client"`
	Log    Log           `kdl:"log"`
}

// Task represents the task configuration.
type Task struct {
	Name   string `kdl:"name"`
	ID     string `kdl:"id"`
	SignID string `kdl:"sign_id"`
}

// Login represents the login configuration.
type Login struct {
	StudentID     string `kdl:"student_id"`
	Password      string `kdl:"password"`
	Authorization string `kdl:"authorization"`
}

// Notify represents the notify configuration.
type Notify struct {
	Ntfy struct {
		Topic string `kdl:"topic"`
	} `kdl:"ntfy"`
	Email struct {
		Host     string `kdl:"host"`
		Port     string `kdl:"port"`
		Username string `kdl:"username"`
		Password string `kdl:"password"`
		From     string `kdl:"from"`
		To       string `kdl:"to"`
	} `kdl:"email"`
}

// Log represents the log configuration.
type Log struct {
	Console bool `kdl:"console"`
}
