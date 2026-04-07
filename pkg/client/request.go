package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

const (
	requestRetryAttempts = 15
	requestRetryDelay    = 200 * time.Millisecond
)

type requestSpec struct {
	method    string
	path      string
	rawQuery  string
	bodyBytes func() []byte
	apply     func(*http.Request)
}

// doRequestWithRetry executes a request with optional retry logic.
// doRequestWithRetry 执行请求，并在需要时应用重试逻辑。
func (c *Client) doRequestWithRetry(ctx context.Context, retry bool, spec requestSpec) (*http.Response, error) {
	attempts := 1
	if retry {
		attempts = requestRetryAttempts
	}

	var lastErr error
	for attempt := 1; attempt <= attempts; attempt++ {
		req, err := c.newRequest(ctx, spec)
		if err != nil {
			return nil, err
		}

		resp, err := c.httpClient.Do(req)
		if err == nil {
			return resp, nil
		}

		lastErr = err
		if !retry || !isRetryableRequestError(err) || attempt == attempts {
			break
		}

		if err := waitForRetry(ctx, requestRetryDelay*time.Duration(attempt)); err != nil {
			return nil, err
		}
	}

	return nil, fmt.Errorf("do request: %w", lastErr)
}

// newRequest builds an HTTP request from the given spec.
// newRequest 根据给定的规格构造 HTTP 请求。
func (c *Client) newRequest(ctx context.Context, spec requestSpec) (*http.Request, error) {
	u := url.URL{Scheme: "http", Host: c.config.Host, Path: spec.path, RawQuery: spec.rawQuery}

	var body io.Reader
	if spec.bodyBytes != nil {
		body = bytes.NewReader(spec.bodyBytes())
	}

	req, err := http.NewRequestWithContext(ctx, spec.method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	c.setCommonHeaders(req)
	if spec.apply != nil {
		spec.apply(req)
	}

	return req, nil
}

// isRetryableRequestError reports whether the request error is retryable.
// isRetryableRequestError 判断该请求错误是否可以重试。
func isRetryableRequestError(err error) bool {
	if errors.Is(err, context.Canceled) {
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	var netErr net.Error
	if errors.As(err, &netErr) {
		return netErr.Timeout() || netErr.Temporary()
	}

	return false
}

// waitForRetry pauses before the next retry attempt.
// waitForRetry 在下一次重试前等待指定时长。
func waitForRetry(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
