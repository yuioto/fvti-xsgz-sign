package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Sign performs the sign-in action for a specific task.
func (c *Client) Sign(ctx context.Context, token, studentID, taskID string) (string, error) {
	u := url.URL{Scheme: "http", Host: c.config.Host, Path: pathSign}

	data := url.Values{
		"ApplyInfo[Id]":            {"00000000-0000-0000-0000-000000000000"},
		"ApplyInfo[OrderId]":       {taskID},
		"ApplyInfo[StudentId]":     {studentID},
		"ApplyInfo[SignWayText]":   {"定位"},
		"ApplyInfo[IsPhoto]":       {"false"},
		"ApplyInfo[IsLocal]":       {"true"},
		"ApplyInfo[IsQrCode]":      {"false"},
		"ApplyInfo[QrCodeContent]": {""},
		"ApplyInfo[IsDWQDW]":       {"0"},
		"ApplyInfo[SingnScope]":    {""},
		"ApplyInfo[Latitude]":      {c.config.Latitude},
		"ApplyInfo[Longitude]":     {c.config.Longitude},
		"ApplyInfo[SingnSite]":     {c.config.SignSite},
		// Empty fields
		"ApplyInfo[InputUser]":      {""},
		"ApplyInfo[InputDate]":      {""},
		"ApplyInfo[collegeNo]":      {""},
		"ApplyInfo[classNo]":        {""},
		"ApplyInfo[qdType]":         {""},
		"ApplyInfo[qdTime]":         {""},
		"ApplyInfo[InsertUserId]":   {""},
		"ApplyInfo[InsertUserName]": {""},
		"ApplyInfo[InsertDate]":     {""},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	c.setCommonHeaders(req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", token)
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response body: %w", err)
	}
	body := string(bodyBytes)

	if resp.StatusCode != StatusSignOkStatusCode {
		return body, fmt.Errorf("%w: status %d, body: %s", ErrSignFailed, resp.StatusCode, body)
	}

	return body, nil
}
