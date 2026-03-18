package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/sblinch/kdl-go"
	"github.com/yuioto/fvti-xsgz-sign/pkg/client"
)

// DefaultNotifyTopic is the default topic for notifications.
const DefaultNotifyTopic = "fvti-xsgz-sign-task-default-status"

// NewDefaultConfig returns a Config with default values.
func NewDefaultConfig() Config {
	return Config{
		Login: Login{
			StudentID: "fvti_student_id",
			Password:  "fvti_xsgz_password",
		},
		Notify: Notify{
			Ntfy: struct {
				Topic string `kdl:"topic"`
			}{
				Topic: DefaultNotifyTopic,
			},
		},
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

	encoder := kdl.NewEncoder(file)
	if err := encoder.Encode(defaultConfig); err != nil {
		return fmt.Errorf("failed to write default config: %w", err)
	}
	log.Println("Create default config file successfully")
	return nil
}
