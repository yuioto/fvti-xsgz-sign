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
		Locale: "zh-CN",
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
			Email: struct {
				Host     string `kdl:"host"`
				Port     string `kdl:"port"`
				Username string `kdl:"username"`
				Password string `kdl:"password"`
				From     string `kdl:"from"`
				FromName string `kdl:"from_name"`
				To       string `kdl:"to"`
			}{
				Host:     "smtp.example.com",
				Port:     "587",
				Username: "user@example.com",
				Password: "email_password",
				From:     "user@example.com",
				FromName: "签到状态",
				To:       "recipient@example.com",
			},
		},
		Log: Log{
			Console: true,
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
