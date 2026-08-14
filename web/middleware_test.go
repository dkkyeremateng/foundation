package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestWrapMiddleware verifies that middleware executes in the order it
// is provided (first middleware runs first) and that nil middleware is
// skipped.
func TestWrapMiddleware(t *testing.T) {
	var order []string

	mw := func(name string) Middleware {
		return func(next Handler) Handler {
			return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
				order = append(order, name)
				return next(ctx, w, r)
			}
		}
	}

	final := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		order = append(order, "handler")
		return nil
	}

	h := wrapMiddleware([]Middleware{mw("first"), nil, mw("second")}, final)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	if err := h(context.Background(), rec, req); err != nil {
		t.Fatalf("wrapped handler = %v, want nil", err)
	}

	want := []string{"first", "second", "handler"}
	if len(order) != len(want) {
		t.Fatalf("execution order = %v, want %v", order, want)
	}
	for i := range want {
		if order[i] != want[i] {
			t.Fatalf("execution order = %v, want %v", order, want)
		}
	}
}
