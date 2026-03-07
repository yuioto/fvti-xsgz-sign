// Package main is the entry point for the sign-in application.
package main

import (
	"context"
	"log"
	"os"

	"github.com/yuioto/fvti-xsgz-sign/internal/app"
	"github.com/yuioto/fvti-xsgz-sign/internal/config"
	"github.com/yuioto/fvti-xsgz-sign/pkg/notify"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := app.Run(cfg); err != nil {
		if cfg.Nofy != "" {
			notifier := notify.New(nil)
			if nErr := notifier.Send(context.Background(), cfg.Nofy, "max", "Sign Failed", err.Error()); nErr != nil {
				log.Printf("Failed to send failure notification: %v", nErr)
			}
		}
		log.Fatalf("Run failed: %v", err)
	}
}

func loadConfig() (config.Config, error) {
	configFile := config.GetConfigFilePath("fvti-xsgz-sign")
	if envConfigFile := os.Getenv("FvtiSign"); envConfigFile != "" {
		configFile = envConfigFile
	}
	return config.LoadConfig(configFile)
}
