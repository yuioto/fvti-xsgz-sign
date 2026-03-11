// Package app contains the main application logic.
package app

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/yuioto/fvti-xsgz-sign/internal/config"
	"github.com/yuioto/fvti-xsgz-sign/pkg/client"
	"github.com/yuioto/fvti-xsgz-sign/pkg/notify"
)

// Run executes the main application logic.
func Run(cfg config.Config) error {
	ctx := context.Background()
	c := client.New(client.WithConfig(cfg.Client))

	// Login if no authorization token
	if cfg.Login.Authorization == "" {
		token, err := c.Login(ctx, cfg.StudentID, cfg.Login.Password)
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

	// Get Task ID if missing
	if cfg.Task.ID == "" {
		taskList, err := c.GetTaskList(ctx, cfg.Login.Authorization)
		if err != nil {
			return fmt.Errorf("get task list failed: %w", err)
		}
		taskID, err := taskList.FindTaskIDByName(cfg.Task.Name)
		if err != nil {
			return fmt.Errorf("find task failed: %w", err)
		}
		cfg.Task.ID = taskID
	}

	// Sign
	_, err = c.Sign(ctx, cfg.Login.Authorization, cfg.StudentID, cfg.Task.ID)
	if err != nil {
		return fmt.Errorf("sign failed: %w", err)
	}

	// Verify
	taskList, err := c.GetTaskList(ctx, cfg.Login.Authorization)
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
		cfg.StudentID, cfg.Task.Name, cfg.Task.ID, cfg.Task.SignID)
	// log.Println("Sign successful:", msg)

	if cfg.Nofy != "" {
		notifier := notify.New(nil)
		if err := notifier.Send(ctx, cfg.Nofy, "high", "Sign Done", msg); err != nil {
			log.Printf("Failed to send notification: %v", err)
		}
	}

	return nil
}
