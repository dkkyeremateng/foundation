package logger_test

import (
	"testing"

	"github.com/dkkyeremateng/foundation/logger"
)

// TestNew verifies that New constructs a working sugared logger.
func TestNew(t *testing.T) {
	log, err := logger.New("TEST")
	if err != nil {
		t.Fatalf("New() = %v, want nil", err)
	}
	if log == nil {
		t.Fatal("New() returned a nil logger")
	}

	// The logger should accept log calls without panicking.
	log.Info("test message")
	log.Errorw("test error message", "key", "value")

	// Sync may return an error on some platforms when flushing stdout;
	// it is called only to exercise the logger's flush path.
	_ = log.Sync()
}
