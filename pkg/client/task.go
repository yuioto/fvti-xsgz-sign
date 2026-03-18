package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strings"
)

// GetTaskList retrieves the list of tasks.
func (c *Client) GetTaskList(ctx context.Context, token string) (*TaskList, error) {
	u := url.URL{Scheme: "http", Host: c.config.Host, Path: pathGetTaskList}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	c.setCommonHeaders(req)
	req.Header.Set("Authorization", token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%w: status %d, body: %s", ErrGetTaskListFailed, resp.StatusCode, string(body))
	}

	var taskList TaskList
	if err := json.NewDecoder(resp.Body).Decode(&taskList); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &taskList, nil
}

// Rule represents a criterion for filtering and scoring an Item.
// The ok return value indicates whether the item satisfies the rule.
// The score represents the weight or priority of the item; higher scores
// are given higher preference during selection.
type Rule func(i *Item) (ok bool, score int)

// Scoring weights.
const (
	HighestPriority = 1000
	Important       = 100
	Preference      = 10
)

// TODO: switch to SelectAll and support multiple results.

// Select returns the Item that satisfies all provided rules and has the
// highest cumulative score.
//
// For an item to be considered, it must return ok=true for every rule
// in the rules slice. If multiple items satisfy all rules, Select picks
// the one with the maximum total score.
//
// Select returns an error if the TaskList is empty or if no items in the
// list satisfy all the provided rules.
func (tl *TaskList) Select(rules ...Rule) (*Item, error) {
	if len(tl.List.Items) == 0 {
		return nil, errors.New("task list is empty")
	}

	var bestItem *Item
	bestScore := math.MinInt

	for i := range tl.List.Items {
		item := &tl.List.Items[i]

		score, ok := func(item *Item) (int, bool) {
			total := 0
			for _, r := range rules {
				ok, score := r(item)
				if !ok {
					return total, false
				}
				total += score
			}
			return total, true
		}(item)

		if !ok {
			continue
		}

		if bestItem == nil || score > bestScore {
			bestItem = item
			bestScore = score
		}
	}

	if bestItem == nil {
		return nil, errors.New("no matching item")
	}

	return bestItem, nil
}

// TODO: remove Select* rule helpers and rename to non-prefixed versions
// (e.g., SelectNonMakeup -> NonMakeup, SelectUnsigned -> Unsigned)
// after client-side migration is complete.

// Deprecated: use NonMakeup instead.
// SelectNonMakeup prefers tasks that are not makeup sign-in tasks.
// It always returns ok=true and assigns a higher score to tasks whose name
// does not contain the makeup keyword.
func SelectNonMakeup(i *Item) (ok bool, score int) {
	ok = true
	if !strings.Contains(i.Name, MakeupSign) {
		return ok, Preference
	}
	return ok, score
}

// Deprecated: use Unsigned instead.
// SelectUnsigned filters out signed tasks and prioritizes unsigned ones.
// It returns ok=false if the task has already been signed.
func SelectUnsigned(i *Item) (bool, int) {
	return i.QD != StatusSignSuccessfullyOk, Important
}

// Deprecated: use Select instead.
// FindTaskIDByName finds a task ID by its name.
func (tl *TaskList) FindTaskIDByName(name string) (string, error) {
	for _, item := range tl.List.Items {
		if item.Name == name {
			return item.ID, nil
		}
	}
	return "", fmt.Errorf("%w: name %q", ErrTaskNotFound, name)
}

// Deprecated: use Select instead.
// FindSignIDByTaskID finds a sign ID by its task ID.
func (tl *TaskList) FindSignIDByTaskID(taskID string) (string, error) {
	for _, item := range tl.List.Items {
		if item.ID == taskID {
			return item.SignID, nil
		}
	}
	return "", fmt.Errorf("%w: id %q", ErrTaskNotFound, taskID)
}

// IsTaskSigned checks if a task is signed.
func (tl *TaskList) IsTaskSigned(taskID string) (bool, error) {
	for _, item := range tl.List.Items {
		if item.ID == taskID {
			return item.QD == StatusSignSuccessfullyOk, nil
		}
	}
	return false, fmt.Errorf("%w: id %q", ErrTaskNotFound, taskID)
}

// TaskList represents the response from the task list API.
type TaskList struct {
	List List `json:"List"`
}

// List contains the items.
type List struct {
	Items []Item `json:"Items"`
}

// Item represents a single task item.
type Item struct {
	ID         string `json:"Id"`
	Name       string `json:"Name"`
	QD         string `json:"QD"`
	SignID     string `json:"SignID"`
	QDTimeText string `json:"QDTimeText"`
}
