// Package main is the entry point for the sign-in application.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/yuioto/fvti-xsgz-sign/internal/app"
	"github.com/yuioto/fvti-xsgz-sign/internal/config"
	"github.com/yuioto/fvti-xsgz-sign/internal/emailhtml"
	"github.com/yuioto/fvti-xsgz-sign/internal/i18n"
	"github.com/yuioto/fvti-xsgz-sign/pkg/notify"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf(i18n.T("", "error.load_config"), err)
	}

	ctx := context.Background()
	result, runErr := app.Run(cfg)

	locale := cfg.Locale
	runAt := time.Now().In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:04:05")

	var title, status, action string
	if runErr != nil {
		title = i18n.T(locale, "notify.failure_title")
		status = i18n.T(locale, "notify.status_failure")
		action = runErr.Error()
	} else {
		title = i18n.T(locale, "notify.success_title")
		status = i18n.T(locale, "notify.status_success")
		action = fmt.Sprintf(i18n.T(locale, "notify.email_action_template"), result.TaskName, result.TaskID)
	}

	plainMessage := fmt.Sprintf("%s\n%s\n%s\n",
		fmt.Sprintf("%s: %s", i18n.T(locale, "notify.email_status_label"), status),
		fmt.Sprintf("%s: %s", i18n.T(locale, "notify.email_action_label"), action),
		fmt.Sprintf("%s: %s", i18n.T(locale, "notify.email_run_at_label"), runAt),
	)
	htmlMessage := emailhtml.FormatSignEmailHTML(locale, status, action, runAt, result.NotifyTasks)

	sendNotification(ctx, cfg, title, plainMessage, htmlMessage)

	if runErr != nil {
		log.Fatalf(i18n.T(locale, "error.run_failed"), runErr)
	}
}

func sendNotification(ctx context.Context, cfg config.Config, title, plainMessage, htmlMessage string) {
	locale := cfg.Locale

	if cfg.Notify.Ntfy.Topic != "" {
		notifier := notify.New(nil)
		if err := notifier.Send(ctx, cfg.Notify.Ntfy.Topic, "high", title, plainMessage); err != nil {
			log.Printf(i18n.T(locale, "error.ntfy_send_failed"), err)
		}
	}

	emailCfg := cfg.Notify.Email
	if emailCfg.Host != "" && emailCfg.To != "" {
		emailClient := notify.NewEmail(notify.EmailConfig{
			Host:     emailCfg.Host,
			Port:     emailCfg.Port,
			Username: emailCfg.Username,
			Password: emailCfg.Password,
			From:     emailCfg.From,
			FromName: emailCfg.FromName,
			To:       emailCfg.To,
			Cc:       emailCfg.Cc,
		})
		if err := emailClient.Send(ctx, title, htmlMessage); err != nil {
			log.Printf(i18n.T(locale, "error.email_send_failed"), err)
		}
	}
}

func loadConfig() (config.Config, error) {
	configFile := config.GetConfigFilePath("fvti-xsgz-sign")
	if envConfigFile := os.Getenv("FvtiSign"); envConfigFile != "" {
		configFile = envConfigFile
	}
	return config.LoadConfig(configFile)
}
