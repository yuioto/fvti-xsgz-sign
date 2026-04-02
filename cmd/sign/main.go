// Package main is the entry point for the sign-in application.
package main

import (
	"context"
	"log"
	"os"

	"github.com/yuioto/fvti-xsgz-sign/internal/app"
	"github.com/yuioto/fvti-xsgz-sign/internal/config"
	"github.com/yuioto/fvti-xsgz-sign/internal/i18n"
	"github.com/yuioto/fvti-xsgz-sign/pkg/notify"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf(i18n.T("", "error.load_config"), err)
	}

	if err := app.Run(cfg); err != nil {
		if cfg.Notify.Ntfy.Topic != "" {
			notifier := notify.New(nil)
			if nErr := notifier.Send(context.Background(), cfg.Notify.Ntfy.Topic, "max", i18n.T(cfg.Locale, "notify.failure_title"), err.Error()); nErr != nil {
				log.Printf(i18n.T(cfg.Locale, "notify.notify_send_failed"), nErr)
			}
		}
		log.Fatalf(i18n.T(cfg.Locale, "error.run_failed"), err)
	}
}

func loadConfig() (config.Config, error) {
	configFile := config.GetConfigFilePath("fvti-xsgz-sign")
	if envConfigFile := os.Getenv("FvtiSign"); envConfigFile != "" {
		configFile = envConfigFile
	}
	return config.LoadConfig(configFile)
}
