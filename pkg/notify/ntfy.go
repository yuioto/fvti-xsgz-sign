// Package notify provides notification services.
package notify

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var (
	// ErrNtfySendFailed is returned when ntfy notification send fails.
	ErrNtfySendFailed = fmt.Errorf("ntfy send failed")
)

const (
	// DefaultBaseURL is the default URL for the ntfy service.
	DefaultBaseURL = "https://ntfy.sh/"
	// DefaultAttachURL is the default attachment image URL.
	DefaultAttachURL = "https://i0.hdslb.com/bfs/article/b11ad7419cd98dfb661f23505a996288ef694932.jpg"
)

// Client is a client for sending notifications.
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// New creates a new notification client.
func New(httpClient *http.Client) *Client {
	if httpClient == nil {
		const defaultTimeout = 10 * time.Second
		httpClient = &http.Client{Timeout: defaultTimeout}
	}
	return &Client{
		httpClient: httpClient,
		baseURL:    DefaultBaseURL,
	}
}

// Send sends a notification to the specified topic.
func (c *Client) Send(ctx context.Context, topic, level, title, message string) error {
	url := c.baseURL + topic

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(message))
	if err != nil {
		return fmt.Errorf("%w: create request failed: %v", ErrNtfySendFailed, err)
	}

	req.Header.Set("Markdown", "yes")
	req.Header.Set("Title", title)
	req.Header.Set("Priority", level)
	req.Header.Set("Tags", "fvti,xsgz,sign")
	req.Header.Set("Attach", DefaultAttachURL)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w: do request failed: %v", ErrNtfySendFailed, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: unexpected status code: %d", ErrNtfySendFailed, resp.StatusCode)
	}

	return nil
}

// ErrorKey maps notify errors to i18n translation keys.
func ErrorKey(err error) string {
	if err == nil {
		return ""
	}
	if errors.Is(err, ErrNtfySendFailed) {
		return "error.ntfy_send_failed"
	}
	return "error.notify_failed"
}
