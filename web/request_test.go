package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dimfeld/httptreemux/v5"
)

// TestDecode verifies that Decode unmarshals a valid JSON body into
// the provided value.
func TestDecode(t *testing.T) {
	type person struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	body := `{"name":"Bill","email":"bill@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	var p person
	if err := Decode(req, &p); err != nil {
		t.Fatalf("Decode() = %v, want nil", err)
	}
	if p.Name != "Bill" || p.Email != "bill@example.com" {
		t.Errorf("Decode() = %+v, want name and email populated", p)
	}
}

// TestDecode_UnknownField verifies that Decode rejects a body
// containing fields not present in the target struct.
func TestDecode_UnknownField(t *testing.T) {
	type person struct {
		Name string `json:"name"`
	}

	body := `{"name":"Bill","unknown":"field"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	var p person
	if err := Decode(req, &p); err == nil {
		t.Fatal("Decode() with unknown field = nil, want an error")
	}
}

// TestDecode_MalformedJSON verifies that Decode returns an error for
// a body that is not valid JSON.
func TestDecode_MalformedJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":`))

	var p struct {
		Name string `json:"name"`
	}
	if err := Decode(req, &p); err == nil {
		t.Fatal("Decode() with malformed JSON = nil, want an error")
	}
}

// TestParam verifies that Param returns route parameters stored in the
// request context by httptreemux.
func TestParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users/42", nil)

	ctx := httptreemux.AddParamsToContext(req.Context(), map[string]string{"id": "42"})
	req = req.WithContext(ctx)

	if got := Param(req, "id"); got != "42" {
		t.Errorf("Param(id) = %q, want %q", got, "42")
	}

	if got := Param(req, "missing"); got != "" {
		t.Errorf("Param(missing) = %q, want %q", got, "")
	}
}

// TestParam_NoParams verifies that Param returns the empty string and
// does not panic when the request context carries no params.
func TestParam_NoParams(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users/42", nil)

	if got := Param(req, "id"); got != "" {
		t.Errorf("Param(id) = %q, want %q", got, "")
	}
}
