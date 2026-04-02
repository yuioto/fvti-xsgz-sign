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
			return fmt.Errorf("login failed: %w", err)
		}
		cfg.Login.Authorization = token
	}

	// Check leave status
	leaveList, err := c.GetLeaveList(ctx, cfg.Login.Authorization)
	if err != nil {
		return fmt.Errorf("failed to get leave list: %w", err)
	}
	if leaveList.IsOnLeave() {
		return errors.New("student is currently on leave")
	}

	// Get TaskList
	taskList, err := c.GetTaskList(ctx, cfg.Login.Authorization)
	if err != nil {
		return fmt.Errorf("get task list failed: %w", err)
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
			return fmt.Errorf("task id %q not found: %w", cfg.Task.ID, err)
		}
	case cfg.Task.Name != "":
		task, err = taskList.Select(
			func(i *client.Item) (ok bool, score int) {
				return i.Name == cfg.Task.Name, 0
			},
		)

		if err != nil {
			return fmt.Errorf("find task by name failed: %w", err)
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
			message := fmt.Sprintf("automatic task selection failed: %v", err)
			sendSignNotification(ctx, cfg, "Sign Failed", message)
			return fmt.Errorf("no matching task found: %w", err)
		}
	}

	cfg.Task.ID = task.ID

	// Sign
	_, err = c.Sign(ctx, cfg.Login.Authorization, cfg.Login.StudentID, cfg.Task.ID)
	if err != nil {
		return fmt.Errorf("sign failed: %w", err)
	}

	// Verify
	taskList, err = c.GetTaskList(ctx, cfg.Login.Authorization)
	if err != nil {
		return fmt.Errorf("verification failed (get list): %w", err)
	}

	signed, err := taskList.IsTaskSigned(cfg.Task.ID)
	if err != nil {
		return fmt.Errorf("verification failed (check status): %w", err)
	}
	if !signed {
		return errors.New("server returned success but task is not marked as signed")
	}

	// Get SignID for notification/logging
	signID, err := taskList.FindSignIDByTaskID(cfg.Task.ID)
	if err != nil {
		return fmt.Errorf("get SignId failed: %w", err)
	}
	cfg.Task.SignID = signID

	msg := fmt.Sprintf("StudentId: %s Task.Name: %s Task.Id: %s Task.SignId: %s",
		cfg.Login.StudentID, cfg.Task.Name, cfg.Task.ID, cfg.Task.SignID)
	// log.Println("Sign successful:", msg)

	if cfg.Notify.Ntfy.Topic != "" {
		notifier := notify.New(nil)
		if err := notifier.Send(ctx, cfg.Notify.Ntfy.Topic, "high", "Sign Done", msg); err != nil {
			log.Printf("Failed to send notification: %v", err)
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
		log.Printf("Sending email notification to %s via %s:%s", emailCfg.To, emailCfg.Host, emailCfg.Port)
		if err := emailClient.Send(ctx, "Sign Done", msg); err != nil {
			log.Printf("Failed to send email notification: %v", err)
		} else {
			log.Println("Email notification sent successfully")
		}
	}

	return nil
}

func sendSignNotification(ctx context.Context, cfg config.Config, title, message string) {
	if cfg.Notify.Ntfy.Topic != "" {
		notifier := notify.New(nil)
		if err := notifier.Send(ctx, cfg.Notify.Ntfy.Topic, "high", title, message); err != nil {
			log.Printf("Failed to send notification: %v", err)
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
			log.Printf("Failed to send email notification: %v", err)
		}
	}
}
