// Package app contains the main application logic.
package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/yuioto/fvti-xsgz-sign/internal/config"
	"github.com/yuioto/fvti-xsgz-sign/internal/i18n"
	"github.com/yuioto/fvti-xsgz-sign/pkg/client"
	"github.com/yuioto/fvti-xsgz-sign/pkg/notify"
)

// Run executes the main application logic.
func Run(cfg config.Config) error {
	if !cfg.Log.Console {
		log.SetOutput(io.Discard)
	} else {
		log.SetOutput(os.Stderr)
	}
	ctx := context.Background()
	c := client.New(client.WithConfig(cfg.Client))

	// Login if no authorization token
	if cfg.Login.Authorization == "" {
		token, err := c.Login(ctx, cfg.Login.StudentID, cfg.Login.Password)
		if err != nil {
			return translateClientError(cfg, err, "error.login_failed")
		}
		cfg.Login.Authorization = token
	}

	// Check leave status
	leaveList, err := c.GetLeaveList(ctx, cfg.Login.Authorization)
	if err != nil {
		return translateClientError(cfg, err, "error.get_leave_failed")
	}
	if leaveList.IsOnLeave() {
		return errors.New(i18n.T(cfg.Locale, "error.on_leave"))
	}

	// Get TaskList
	taskList, err := c.GetTaskList(ctx, cfg.Login.Authorization)
	if err != nil {
		return translateClientError(cfg, err, "error.get_task_list_failed")
	}

	// Get Task
	var task *client.Item

	switch {
	case cfg.Task.ID != "":
		task, err = taskList.Select(
			func(i *client.Item) (ok bool, score int) {
				return i.ID == cfg.Task.ID, 0
			},
		)

		if err != nil {
			return fmt.Errorf(i18n.T(cfg.Locale, "error.task_id_not_found"), cfg.Task.ID, err)
		}
	case cfg.Task.Name != "":
		task, err = taskList.Select(
			func(i *client.Item) (ok bool, score int) {
				return i.Name == cfg.Task.Name, 0
			},
		)

		if err != nil {
			return fmt.Errorf(i18n.T(cfg.Locale, "error.task_name_not_found"), err)
		}
	default:
		task, err = taskList.Select(
			client.SelectUnsigned,
			client.SelectNonMakeup,
			func(i *client.Item) (ok bool, score int) {
				return i.QD != client.SignOutOfTimeRange, client.Important

				/*
					// Check i.QDTimeText version

					layout := "15:04"
					parts := strings.Split(i.QDTimeText, "至")
					if len(parts) != 2 {
						return false, 0
					}

					start, err1 := time.Parse(layout, strings.TrimSpace(parts[0]))
					end, err2 := time.Parse(layout, strings.TrimSpace(parts[1]))
					if err1 != nil || err2 != nil {
						return false, 0
					}

					now := time.Now()
					if now.Before(start) || now.After(end) {
						return false, 0
					}

					return true, client.Important
				*/
			},
		)

		if err != nil {
			message := fmt.Sprintf(i18n.T(cfg.Locale, "error.task_auto_select_failed"), err)
			sendSignNotification(ctx, cfg, i18n.T(cfg.Locale, "notify.failure_title"), message)
			return errors.New(i18n.T(cfg.Locale, "error.no_matching_task"))
		}
	}

	cfg.Task.ID = task.ID

	// Sign
	_, err = c.Sign(ctx, cfg.Login.Authorization, cfg.Login.StudentID, cfg.Task.ID)
	if err != nil {
		return translateClientError(cfg, err, "error.sign_failed")
	}

	// Verify
	taskList, err = c.GetTaskList(ctx, cfg.Login.Authorization)
	if err != nil {
		return translateClientError(cfg, err, "error.verify_get_task_list_failed")
	}

	signed, err := taskList.IsTaskSigned(cfg.Task.ID)
	if err != nil {
		return fmt.Errorf(i18n.T(cfg.Locale, "error.verify_status_failed"), err)
	}
	if !signed {
		return errors.New(i18n.T(cfg.Locale, "error.server_succeed_not_signed"))
	}

	// Get SignID for notification/logging
	signID, err := taskList.FindSignIDByTaskID(cfg.Task.ID)
	if err != nil {
		return fmt.Errorf(i18n.T(cfg.Locale, "error.get_signid_failed"), err)
	}
	cfg.Task.SignID = signID

	msg := fmt.Sprintf("StudentId: %s Task.Name: %s Task.Id: %s Task.SignId: %s",
		cfg.Login.StudentID, cfg.Task.Name, cfg.Task.ID, cfg.Task.SignID)

	notifyTitle := i18n.T(cfg.Locale, "notify.success_title")
	if cfg.Notify.Ntfy.Topic != "" {
		notifier := notify.New(nil)
		if err := notifier.Send(ctx, cfg.Notify.Ntfy.Topic, "high", notifyTitle, msg); err != nil {
			key := notify.ErrorKey(err)
			if key == "" {
				key = "error.ntfy_send_failed"
			}
			log.Printf(i18n.T(cfg.Locale, key), err)
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
			To:       emailCfg.To,
		})
		log.Printf(i18n.T(cfg.Locale, "notify.email_sending"), emailCfg.To, emailCfg.Host, emailCfg.Port)
		if err := emailClient.Send(ctx, notifyTitle, msg); err != nil {
			key := notify.ErrorKey(err)
			if key == "" {
				key = "error.email_send_failed"
			}
			log.Printf(i18n.T(cfg.Locale, key), err)
		} else {
			log.Println(i18n.T(cfg.Locale, "notify.email_send_success"))
		}
	}

	return nil
}

func sendSignNotification(ctx context.Context, cfg config.Config, title, message string) {
	if cfg.Notify.Ntfy.Topic != "" {
		notifier := notify.New(nil)
		if err := notifier.Send(ctx, cfg.Notify.Ntfy.Topic, "high", title, message); err != nil {
			log.Printf(i18n.T(cfg.Locale, "error.ntfy_send_failed"), err)
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
			To:       emailCfg.To,
		})
		if err := emailClient.Send(ctx, title, message); err != nil {
			log.Printf(i18n.T(cfg.Locale, "error.email_send_failed"), err)
		}
	}
}

// translateClientError maps client-layer errors to i18n keys and returns formatted error.
//
// This helps keep the client package free of direct localization logic. Only the
// application layer performs translation with locale templates.
func translateClientError(cfg config.Config, err error, fallbackKey string) error {
	key := client.ErrorKey(err)
	if key == "" {
		key = fallbackKey
	}

	return fmt.Errorf(i18n.T(cfg.Locale, key), err)
}
