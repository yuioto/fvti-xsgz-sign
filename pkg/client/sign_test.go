package client

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestSignRetriesTemporaryNetworkErrors(t *testing.T) {
	var attempts int

	c := New(WithHTTPClient(&http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			attempts++

			bodyBytes, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatalf("read request body: %v", err)
			}
			values, err := url.ParseQuery(string(bodyBytes))
			if err != nil {
				t.Fatalf("parse request body: %v", err)
			}
			if got := values.Get("ApplyInfo[OrderId]"); got != "task-1" {
				t.Fatalf("expected task id in request body, got %q", got)
			}
			if got := values.Get("ApplyInfo[StudentId]"); got != "student-1" {
				t.Fatalf("expected student id in request body, got %q", got)
			}

			if attempts < 3 {
				return nil, retryableNetError{}
			}

			return &http.Response{
				StatusCode: StatusSignOkStatusCode,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader("ok")),
			}, nil
		}),
	}))

	got, err := c.Sign(context.Background(), "Bearer token", "student-1", "task-1")
	if err != nil {
		t.Fatalf("Sign failed: %v", err)
	}
	if got != "ok" {
		t.Fatalf("expected response body %q, got %q", "ok", got)
	}
	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
}
