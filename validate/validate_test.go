package validate

import "testing"

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
