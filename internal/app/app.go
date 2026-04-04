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
)

// RunResult holds the outcome of a Run call.
type RunResult struct {
	TaskName    string
	TaskID      string
	SignID      string
	NotifyTasks []client.TaskSummary
}

// Run executes the main application logic.
func Run(cfg config.Config) (RunResult, error) {
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
			return RunResult{}, translateClientError(cfg, err, "error.login_failed")
		}
		cfg.Login.Authorization = token
	}

	// Check leave status
	leaveList, err := c.GetLeaveList(ctx, cfg.Login.Authorization)
	if err != nil {
		return RunResult{}, translateClientError(cfg, err, "error.get_leave_failed")
	}
	if leaveList.IsOnLeave() {
		return RunResult{}, errors.New(i18n.T(cfg.Locale, "error.on_leave"))
	}

	// Get TaskList
	taskList, err := c.GetTaskList(ctx, cfg.Login.Authorization)
	if err != nil {
		return RunResult{}, translateClientError(cfg, err, "error.get_task_list_failed")
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
			return RunResult{NotifyTasks: taskList.ToSummary()},
				fmt.Errorf(i18n.T(cfg.Locale, "error.task_id_not_found"), cfg.Task.ID, err)
		}
	case cfg.Task.Name != "":
		task, err = taskList.Select(
			func(i *client.Item) (ok bool, score int) {
				return i.Name == cfg.Task.Name, 0
			},
		)

		if err != nil {
			return RunResult{NotifyTasks: taskList.ToSummary()},
				fmt.Errorf(i18n.T(cfg.Locale, "error.task_name_not_found"), err)
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
			return RunResult{NotifyTasks: taskList.ToSummary()},
				fmt.Errorf(i18n.T(cfg.Locale, "error.task_auto_select_failed"), err)
		}
	}

	cfg.Task.ID = task.ID

	// Sign
	_, err = c.Sign(ctx, cfg.Login.Authorization, cfg.Login.StudentID, cfg.Task.ID)
	if err != nil {
		return RunResult{TaskName: task.Name, TaskID: task.ID},
			translateClientError(cfg, err, "error.sign_failed")
	}

	// Verify (reuse the fetched taskList for NotifyTasks, avoiding an extra API call)
	taskList, err = c.GetTaskList(ctx, cfg.Login.Authorization)
	if err != nil {
		return RunResult{TaskName: task.Name, TaskID: task.ID},
			translateClientError(cfg, err, "error.verify_get_task_list_failed")
	}

	signed, err := taskList.IsTaskSigned(cfg.Task.ID)
	if err != nil {
		return RunResult{TaskName: task.Name, TaskID: task.ID, NotifyTasks: taskList.ToSummary()},
			fmt.Errorf(i18n.T(cfg.Locale, "error.verify_status_failed"), err)
	}
	if !signed {
		return RunResult{TaskName: task.Name, TaskID: task.ID, NotifyTasks: taskList.ToSummary()},
			errors.New(i18n.T(cfg.Locale, "error.server_succeed_not_signed"))
	}

	// Get SignID for notification/logging
	signID, err := taskList.FindSignIDByTaskID(cfg.Task.ID)
	if err != nil {
		return RunResult{TaskName: task.Name, TaskID: task.ID, NotifyTasks: taskList.ToSummary()},
			fmt.Errorf(i18n.T(cfg.Locale, "error.get_signid_failed"), err)
	}

	return RunResult{
		TaskName:    task.Name,
		TaskID:      task.ID,
		SignID:      signID,
		NotifyTasks: taskList.ToSummary(),
	}, nil
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
