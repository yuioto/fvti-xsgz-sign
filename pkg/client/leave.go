package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	u := url.URL{Scheme: "http", Host: c.config.Host, Path: pathLeaveList, RawQuery: "pageIndex=1"}
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
