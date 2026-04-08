package client

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

type retryableNetError struct{}

func (retryableNetError) Error() string   { return "temporary network error" }
func (retryableNetError) Timeout() bool   { return true }
func (retryableNetError) Temporary() bool { return true }

func TestDoRequestWithRetryRetriesTemporaryNetworkErrors(t *testing.T) {
	var attempts int
	var bodyCalls int

	c := New(WithHTTPClient(&http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			attempts++

			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("read request body: %v", err)
			}
			if string(bodyBytes) != "payload" {
				t.Fatalf("unexpected request body: %q", string(bodyBytes))
			}

			if attempts < 3 {
				return nil, retryableNetError{}
			}

			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader("ok")),
			}, nil
		}),
	}))

	resp, err := c.doRequestWithRetry(context.Background(), true, requestSpec{
		method: http.MethodPost,
		path:   "/retry",
		bodyBytes: func() []byte {
			bodyCalls++
			return []byte("payload")
		},
	})
	if err != nil {
		t.Fatalf("doRequestWithRetry failed: %v", err)
	}
	defer resp.Body.Close()

	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
	if bodyCalls != 3 {
		t.Fatalf("expected body factory to run 3 times, got %d", bodyCalls)
	}
}

func TestDoRequestWithRetryStopsOnPermanentError(t *testing.T) {
	var attempts int

	c := New(WithHTTPClient(&http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			attempts++
			return nil, errors.New("boom")
		}),
	}))

	_, err := c.doRequestWithRetry(context.Background(), true, requestSpec{
		method: http.MethodGet,
		path:   "/fail",
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if attempts != 1 {
		t.Fatalf("expected 1 attempt for permanent error, got %d", attempts)
	}
}
