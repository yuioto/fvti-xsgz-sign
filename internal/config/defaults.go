package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"github.com/yuioto/fvti-xsgz-sign/pkg/client"
)

// DefaultNotifyTopic is the default topic for notifications.
const DefaultNotifyTopic = "fvti-xsgz-sign-task-default-status"

// NewDefaultConfig returns a Config with default values.
func NewDefaultConfig() Config {
	return Config{
		StudentID: "{fvti_student_id}",
		Login: Login{
			Password:      "{fvti_xsgz_password}",
			Authorization: "",
		},
		Task: Task{
			Name:   "",
			ID:     "",
			SignID: "",
		},
		Nofy: DefaultNotifyTopic,
		Client: client.Config{
			Host:      client.DefaultHost,
			UserAgent: client.DefaultUserAgent,
			Latitude:  client.DefaultLatitude,
			Longitude: client.DefaultLongitude,
			SignSite:  client.DefaultSignSite,
		},
	}
}

// CreateDefaultConfig creates a default configuration file.
func CreateDefaultConfig(filename string) error {
	defaultConfig := NewDefaultConfig()

	file, err := os.Create(filepath.Clean(filename))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(defaultConfig); err != nil {
		return fmt.Errorf("failed to write default config: %w", err)
	}
	log.Println("Create default config file successfully")
	return nil
}
