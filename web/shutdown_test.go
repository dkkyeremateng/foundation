package web

import (
	"errors"
	"fmt"
	"testing"
)

// TestNewShutdownError verifies that NewShutdownError produces an
// error whose message matches the provided message.
func TestNewShutdownError(t *testing.T) {
	err := NewShutdownError("integrity issue")
	if err == nil {
		t.Fatal("NewShutdownError() = nil, want an error")
	}
	if err.Error() != "integrity issue" {
		t.Errorf("Error() = %q, want %q", err.Error(), "integrity issue")
	}
}

// TestIsShutdown verifies that IsShutdown detects shutdown errors and
// rejects other errors.
func TestIsShutdown(t *testing.T) {
	tt := []struct {
		name string
		err  error
		want bool
	}{
		{"shutdown error", NewShutdownError("boom"), true},
		{"wrapped shutdown error", fmt.Errorf("handler: %w", NewShutdownError("boom")), true},
		{"plain error", errors.New("boom"), false},
		{"nil error", nil, false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsShutdown(tc.err); got != tc.want {
				t.Errorf("IsShutdown(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}
