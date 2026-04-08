package notify

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorKey(t *testing.T) {
	if got := ErrorKey(nil); got != "" {
		t.Fatalf("expected empty for nil, got %q", got)
	}

	if got := ErrorKey(fmt.Errorf("%w: network error", ErrNtfySendFailed)); got != "error.ntfy_send_failed" {
		t.Fatalf("expected error.ntfy_send_failed, got %q", got)
	}

	if got := ErrorKey(fmt.Errorf("%w: smtp error", ErrEmailSendFailed)); got != "error.notify_failed" {
		t.Fatalf("expected error.notify_failed, got %q", got)
	}

	if got := ErrorKey(errors.New("other")); got != "error.notify_failed" {
		t.Fatalf("expected error.notify_failed for unknown, got %q", got)
	}
}
