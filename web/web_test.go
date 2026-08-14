package web_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/dkkyeremateng/foundation/web"
)

// TestApp_Routing verifies that a handler registered on a route group
// is reached for the matching method and path, and that the framework
// populates the request context values.
func TestApp_Routing(t *testing.T) {
	shutdown := make(chan os.Signal, 1)
	app := web.NewApp(shutdown)

	type result struct {
		traceID string
		now     time.Time
	}
	resCh := make(chan result, 1)

	app.Handle(http.MethodGet, "v1", "/users", func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		v, err := web.GetValues(ctx)
		if err != nil {
			return err
		}
		resCh <- result{traceID: v.TraceID, now: v.Now}
		return web.Respond(ctx, w, map[string]string{"status": "ok"}, http.StatusOK)
	})

	before := time.Now()
	req := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)
	after := time.Now()

	if rec.Code != http.StatusOK {
		t.Fatalf("status code = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want %q", ct, "application/json")
	}

	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}
	if body["status"] != "ok" {
		t.Errorf("response body = %v, want status=ok", body)
	}

	select {
	case res := <-resCh:
		if res.traceID == "" {
			t.Error("TraceID is empty, want a non-empty value")
		}
		if res.now.Before(before) || res.now.After(after) {
			t.Errorf("Now = %v, want between %v and %v", res.now, before, after)
		}
	case <-time.After(time.Second):
		t.Fatal("handler was not called")
	}

	// The shutdown channel must not have been signaled.
	select {
	case sig := <-shutdown:
		t.Fatalf("unexpected shutdown signal: %v", sig)
	default:
	}
}

// TestApp_Middleware verifies that application and handler middleware
// execute in order around the handler.
func TestApp_Middleware(t *testing.T) {
	shutdown := make(chan os.Signal, 1)

	var order []string
	mw := func(name string) web.Middleware {
		return func(next web.Handler) web.Handler {
			return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
				order = append(order, name)
				return next(ctx, w, r)
			}
		}
	}

	app := web.NewApp(shutdown, mw("app"))

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		order = append(order, "handler")
		w.Header().Set("X-Test", "middleware-ran")
		return web.Respond(ctx, w, nil, http.StatusNoContent)
	}
	app.Handle(http.MethodGet, "", "/test", handler, mw("route"))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status code = %d, want %d", rec.Code, http.StatusNoContent)
	}
	if got := rec.Header().Get("X-Test"); got != "middleware-ran" {
		t.Errorf("X-Test header = %q, want %q", got, "middleware-ran")
	}

	want := []string{"app", "route", "handler"}
	if len(order) != len(want) {
		t.Fatalf("middleware order = %v, want %v", order, want)
	}
	for i := range want {
		if order[i] != want[i] {
			t.Fatalf("middleware order = %v, want %v", order, want)
		}
	}
}

// TestApp_SignalShutdown verifies that a handler returning an error
// causes the app to signal a shutdown.
func TestApp_SignalShutdown(t *testing.T) {
	shutdown := make(chan os.Signal, 1)
	app := web.NewApp(shutdown)

	app.Handle(http.MethodGet, "", "/fail", func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		return errors.New("handler failure")
	})

	req := httptest.NewRequest(http.MethodGet, "/fail", nil)
	rec := httptest.NewRecorder()
	app.ServeHTTP(rec, req)

	select {
	case sig := <-shutdown:
		if sig != syscall.SIGTERM {
			t.Errorf("shutdown signal = %v, want %v", sig, syscall.SIGTERM)
		}
	case <-time.After(time.Second):
		t.Fatal("no shutdown signal received after handler error")
	}
}
