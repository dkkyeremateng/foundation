package web

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// testContext returns a context carrying a Values value, as the
// framework would during a request.
func testContext() (context.Context, *Values) {
	v := &Values{TraceID: "trace-id", Now: time.Now()}
	return context.WithValue(context.Background(), key, v), v
}

// TestRespond verifies that Respond writes the JSON body, content
// type, and status code, and records the status code in the context.
func TestRespond(t *testing.T) {
	ctx, v := testContext()
	rec := httptest.NewRecorder()

	data := map[string]string{"name": "Bill"}
	if err := Respond(ctx, rec, data, http.StatusOK); err != nil {
		t.Fatalf("Respond() = %v, want nil", err)
	}

	if rec.Code != http.StatusOK {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusOK)
	}

	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want %q", ct, "application/json")
	}

	var got map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}
	if got["name"] != "Bill" {
		t.Errorf("response body = %v, want name=Bill", got)
	}

	if v.StatusCode != http.StatusOK {
		t.Errorf("context StatusCode = %d, want %d", v.StatusCode, http.StatusOK)
	}
}

// TestRespond_NoContent verifies that Respond with http.StatusNoContent
// writes the status code and no body.
func TestRespond_NoContent(t *testing.T) {
	ctx, _ := testContext()
	rec := httptest.NewRecorder()

	if err := Respond(ctx, rec, nil, http.StatusNoContent); err != nil {
		t.Fatalf("Respond() = %v, want nil", err)
	}

	if rec.Code != http.StatusNoContent {
		t.Errorf("status code = %d, want %d", rec.Code, http.StatusNoContent)
	}
	if rec.Body.Len() != 0 {
		t.Errorf("body = %q, want empty", rec.Body.String())
	}
}

// TestRespond_MarshalError verifies that Respond returns an error when
// the data cannot be marshaled to JSON.
func TestRespond_MarshalError(t *testing.T) {
	ctx, _ := testContext()
	rec := httptest.NewRecorder()

	if err := Respond(ctx, rec, make(chan int), http.StatusOK); err == nil {
		t.Fatal("Respond() with unmarshalable data = nil, want an error")
	}
}
