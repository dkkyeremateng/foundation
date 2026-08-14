package web

import (
	"context"
	"testing"
	"time"
)

// TestGetValues verifies that GetValues returns the stored Values and
// errors when no value is present in the context.
func TestGetValues(t *testing.T) {
	v := Values{TraceID: "trace-id", Now: time.Now()}
	ctx := context.WithValue(context.Background(), key, &v)

	got, err := GetValues(ctx)
	if err != nil {
		t.Fatalf("GetValues() = %v, want nil", err)
	}
	if got != &v {
		t.Errorf("GetValues() = %p, want %p", got, &v)
	}

	if _, err := GetValues(context.Background()); err == nil {
		t.Error("GetValues() with empty context = nil, want an error")
	}
}

// TestGetTraceID verifies that GetTraceID returns the stored trace id
// and falls back to the zero value when no value is present.
func TestGetTraceID(t *testing.T) {
	v := Values{TraceID: "trace-id", Now: time.Now()}
	ctx := context.WithValue(context.Background(), key, &v)

	if got := GetTraceID(ctx); got != "trace-id" {
		t.Errorf("GetTraceID() = %q, want %q", got, "trace-id")
	}

	want := "00000000-0000-0000-0000-000000000000"
	if got := GetTraceID(context.Background()); got != want {
		t.Errorf("GetTraceID() with empty context = %q, want %q", got, want)
	}
}

// TestSetStatusCode verifies that SetStatusCode stores the status code
// in the context values and errors when no value is present.
func TestSetStatusCode(t *testing.T) {
	v := Values{TraceID: "trace-id", Now: time.Now()}
	ctx := context.WithValue(context.Background(), key, &v)

	if err := SetStatusCode(ctx, 200); err != nil {
		t.Fatalf("SetStatusCode() = %v, want nil", err)
	}
	if v.StatusCode != 200 {
		t.Errorf("StatusCode = %d, want 200", v.StatusCode)
	}

	if err := SetStatusCode(context.Background(), 200); err == nil {
		t.Error("SetStatusCode() with empty context = nil, want an error")
	}
}
