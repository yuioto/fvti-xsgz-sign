package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LeaveList struct {
	List struct {
		Items []LeaveItem `json:"Items"`
	} `json:"List"`
}

type LeaveItem struct {
	ID         string `json:"Id"`
	StatusName string `json:"StatusName"`
	LeaveTime  string `json:"LeaveTime"`
}

func (c *Client) GetLeaveList(ctx context.Context, token string) (*LeaveList, error) {
	// pageIndex=1 is sufficient because we only need the newest leave status.
	// This keeps the Items slice small and avoids unnecessary iteration later.
	resp, err := c.doRequestWithRetry(ctx, true, requestSpec{
		method:   http.MethodGet,
		path:     pathLeaveList,
		rawQuery: "pageIndex=1",
		apply: func(req *http.Request) {
			req.Header.Set("Authorization", token)
		},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%w: status %d, body: %s", ErrGetLeaveListFailed, resp.StatusCode, string(body))
	}

	var leaveList LeaveList
	if err := json.NewDecoder(resp.Body).Decode(&leaveList); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &leaveList, nil
}

func (l LeaveItem) onLeave() bool {
	return l.StatusName == LeaveStatus
}

func (leaveList LeaveList) IsOnLeave() bool {
	if len(leaveList.List.Items) == 0 {
		return false
	}
	return leaveList.List.Items[0].onLeave()
}
