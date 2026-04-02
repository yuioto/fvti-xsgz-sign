package client

import "testing"

func TestErrorKey(t *testing.T) {
	cases := []struct {
		err      error
		expected string
	}{
		{ErrLoginFailed, "error.login_failed"},
		{ErrGetLeaveListFailed, "error.get_leave_failed"},
		{ErrGetTaskListFailed, "error.get_task_list_failed"},
		{ErrSignFailed, "error.sign_failed"},
		{ErrTaskNotFound, "error.no_matching_task"},
		{nil, ""},
	}

	for _, c := range cases {
		if got := ErrorKey(c.err); got != c.expected {
			t.Fatalf("ErrorKey(%v) expected %q, got %q", c.err, c.expected, got)
		}
	}
}
