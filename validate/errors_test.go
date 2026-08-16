package validate

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

// TestFieldErrors_Error verifies that Error serializes the field
// errors to JSON.
func TestFieldErrors_Error(t *testing.T) {
	fe := FieldErrors{
		{Field: "name", Error: "name is a required field"},
		{Field: "email", Error: "email must be a valid email address"},
	}

	got := fe.Error()

	var decoded []FieldError
	if err := json.Unmarshal([]byte(got), &decoded); err != nil {
		t.Fatalf("Error() = %q, which is not valid JSON: %v", got, err)
	}

	if len(decoded) != len(fe) {
		t.Fatalf("decoded %d field errors, want %d", len(decoded), len(fe))
	}
	for i, fld := range decoded {
		if fld != fe[i] {
			t.Errorf("decoded[%d] = %+v, want %+v", i, fld, fe[i])
		}
	}
}

// TestFieldErrors_Fields verifies that Fields returns the field to
// message map.
func TestFieldErrors_Fields(t *testing.T) {
	fe := FieldErrors{
		{Field: "name", Error: "name is a required field"},
		{Field: "email", Error: "email must be a valid email address"},
	}

	m := fe.Fields()

	if len(m) != 2 {
		t.Fatalf("Fields() returned %d entries, want 2", len(m))
	}
	for _, fld := range fe {
		msg, ok := m[fld.Field]
		if !ok {
			t.Errorf("Fields() missing entry for field %q", fld.Field)
			continue
		}
		if msg != fld.Error {
			t.Errorf("Fields()[%q] = %q, want %q", fld.Field, msg, fld.Error)
		}
	}
}

// TestIsFieldErrors verifies that IsFieldErrors detects FieldErrors
// values and rejects other errors.
func TestIsFieldErrors(t *testing.T) {
	fe := FieldErrors{{Field: "name", Error: "name is a required field"}}

	tt := []struct {
		name string
		err  error
		want bool
	}{
		{"FieldErrors value", fe, true},
		{"wrapped FieldErrors", fmt.Errorf("validate: %w", fe), true},
		{"plain error", errors.New("boom"), false},
		{"nil error", nil, false},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsFieldErrors(tc.err); got != tc.want {
				t.Errorf("IsFieldErrors(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}

// TestGetFieldErrors verifies that GetFieldErrors returns the
// FieldErrors value, and nil for other errors.
func TestGetFieldErrors(t *testing.T) {
	fe := FieldErrors{{Field: "name", Error: "name is a required field"}}

	got := GetFieldErrors(fe)
	if got == nil {
		t.Fatal("GetFieldErrors(FieldErrors) = nil, want the FieldErrors value")
	}
	if len(got) != 1 || got[0] != fe[0] {
		t.Errorf("GetFieldErrors(FieldErrors) = %v, want %v", got, fe)
	}

	if got := GetFieldErrors(errors.New("boom")); got != nil {
		t.Errorf("GetFieldErrors(plain error) = %v, want nil", got)
	}

	if got := GetFieldErrors(nil); got != nil {
		t.Errorf("GetFieldErrors(nil) = %v, want nil", got)
	}
}

// TestGetFieldErrors_ReturnsCopy verifies that the slice returned by
// GetFieldErrors does not share a backing array with the original error,
// so writing through it does not mutate the error.
func TestGetFieldErrors_ReturnsCopy(t *testing.T) {
	fe := FieldErrors{{Field: "name", Error: "name is a required field"}}

	got := GetFieldErrors(fe)
	if got == nil || len(got) != len(fe) {
		t.Fatalf("GetFieldErrors(FieldErrors) = %v, want a copy of %v", got, fe)
	}

	got[0].Error = "redacted"

	if fe[0].Error != "name is a required field" {
		t.Errorf("GetFieldErrors returned an alias: mutating it changed the original "+
			"error to %q; want a copy", fe[0].Error)
	}
}
