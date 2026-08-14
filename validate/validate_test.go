package validate

import "testing"

// TestCheck_Valid verifies that Check accepts a struct that satisfies
// all of its validation tags.
func TestCheck_Valid(t *testing.T) {
	type person struct {
		Name  string `json:"name" validate:"required"`
		Email string `json:"email" validate:"required,email"`
	}

	p := person{Name: "Bill", Email: "bill@example.com"}
	if err := Check(p); err != nil {
		t.Fatalf("Check(%+v) = %v, want nil", p, err)
	}
}

// TestCheck_Invalid verifies that Check returns FieldErrors describing
// each invalid field, keyed by the JSON tag name.
func TestCheck_Invalid(t *testing.T) {
	type person struct {
		Name  string `json:"name" validate:"required"`
		Email string `json:"email" validate:"required,email"`
	}

	err := Check(person{Name: "", Email: "not-an-email"})
	if err == nil {
		t.Fatal("Check() = nil, want a FieldErrors value")
	}

	if !IsFieldErrors(err) {
		t.Fatalf("Check() error is %T, want FieldErrors", err)
	}

	fields := GetFieldErrors(err).Fields()
	if _, ok := fields["name"]; !ok {
		t.Errorf("FieldErrors missing entry for json field %q: %v", "name", fields)
	}
	if _, ok := fields["email"]; !ok {
		t.Errorf("FieldErrors missing entry for json field %q: %v", "email", fields)
	}
}

// TestCheck_Nested verifies that Check validates nested structs and
// reports the nested field using its JSON tag name.
func TestCheck_Nested(t *testing.T) {
	type address struct {
		City string `json:"city" validate:"required"`
	}
	type person struct {
		Name    string  `json:"name" validate:"required"`
		Address address `json:"address"`
	}

	err := Check(person{Name: "Bill", Address: address{City: ""}})
	if err == nil {
		t.Fatal("Check() = nil, want a FieldErrors value")
	}

	if !IsFieldErrors(err) {
		t.Fatalf("Check() error is %T, want FieldErrors", err)
	}

	fields := GetFieldErrors(err).Fields()
	if _, ok := fields["city"]; !ok {
		t.Errorf("FieldErrors missing entry for nested json field %q: %v", "city", fields)
	}
}

// TestCheckID_ValidUUID verifies that CheckID accepts a well-formed
// uuid value.
func TestCheckID_ValidUUID(t *testing.T) {
	if err := CheckID("54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"); err != nil {
		t.Fatalf("CheckID() = %v, want nil", err)
	}
}

// TestGenerateID_AcceptedByCheckID verifies that an id produced by
// GenerateID is accepted by CheckID.
func TestGenerateID_AcceptedByCheckID(t *testing.T) {
	id := GenerateID()
	if err := CheckID(id); err != nil {
		t.Fatalf("GenerateID() produced %q, which CheckID rejected: %v", id, err)
	}
}

// TestGenerateID_Unique verifies that two calls to GenerateID do not
// produce the same value.
func TestGenerateID_Unique(t *testing.T) {
	a, b := GenerateID(), GenerateID()
	if a == b {
		t.Fatalf("two calls to GenerateID() returned the same value %q", a)
	}
}

// TestCheckID_Invalid verifies that CheckID rejects malformed ids.
func TestCheckID_Invalid(t *testing.T) {
	tt := []struct {
		name string
		id   string
	}{
		{"empty string", ""},
		{"not a uuid", "not-a-uuid"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if err := CheckID(tc.id); err == nil {
				t.Errorf("CheckID(%q) = nil, want an error", tc.id)
			}
		})
	}
}
