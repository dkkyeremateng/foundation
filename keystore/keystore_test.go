package keystore_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/dkkyeremateng/foundation/keystore"
)

// generateKey produces an RSA private key for use in tests.
func generateKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generating RSA key: %v", err)
	}
	return key
}

// pemEncode marshals an RSA private key into PEM form.
func pemEncode(t *testing.T, key *rsa.PrivateKey) []byte {
	t.Helper()

	return pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
}

// TestNew verifies that an empty store fails lookups with the
// "kid lookup failed" error.
func TestNew(t *testing.T) {
	ks := keystore.New()
	if ks == nil {
		t.Fatal("New() returned nil")
	}

	if _, err := ks.PrivateKey("unknown"); err == nil ||
		!strings.Contains(err.Error(), "kid lookup failed") {
		t.Errorf("PrivateKey(unknown) = %v, want \"kid lookup failed\" error", err)
	}

	if _, err := ks.PublicKey("unknown"); err == nil ||
		!strings.Contains(err.Error(), "kid lookup failed") {
		t.Errorf("PublicKey(unknown) = %v, want \"kid lookup failed\" error", err)
	}
}

// TestNewMap verifies that a store initialized with NewMap resolves
// the provided keys.
func TestNewMap(t *testing.T) {
	key := generateKey(t)

	ks := keystore.NewMap(map[string]*rsa.PrivateKey{"kid1": key})

	got, err := ks.PrivateKey("kid1")
	if err != nil {
		t.Fatalf("PrivateKey(kid1) = %v, want nil", err)
	}
	if got != key {
		t.Errorf("PrivateKey(kid1) returned a different key than the one provided")
	}
}

// TestAddRemove verifies that keys can be added to and removed from
// the store.
func TestAddRemove(t *testing.T) {
	key := generateKey(t)
	ks := keystore.New()

	ks.Add(key, "kid1")
	if _, err := ks.PrivateKey("kid1"); err != nil {
		t.Fatalf("PrivateKey(kid1) after Add = %v, want nil", err)
	}

	ks.Remove("kid1")
	if _, err := ks.PrivateKey("kid1"); err == nil {
		t.Fatal("PrivateKey(kid1) after Remove = nil, want an error")
	}
}

// TestPrivatePublicKey verifies that the public key returned for a kid
// matches the stored private key's public key.
func TestPrivatePublicKey(t *testing.T) {
	key := generateKey(t)
	ks := keystore.NewMap(map[string]*rsa.PrivateKey{"kid1": key})

	priv, err := ks.PrivateKey("kid1")
	if err != nil {
		t.Fatalf("PrivateKey(kid1) = %v, want nil", err)
	}

	pub, err := ks.PublicKey("kid1")
	if err != nil {
		t.Fatalf("PublicKey(kid1) = %v, want nil", err)
	}

	if pub.N.Cmp(priv.PublicKey.N) != 0 || pub.E != priv.PublicKey.E {
		t.Errorf("PublicKey(kid1) does not match PrivateKey(kid1).PublicKey")
	}
}

// TestNewFS verifies that NewFS loads PEM files from a file system,
// keying them by filename minus the .pem extension, and ignores
// non-PEM files and directories.
func TestNewFS(t *testing.T) {
	key := generateKey(t)

	fsys := fstest.MapFS{
		"54bb2165-71e1-41a6-af3e-7da4a0e1e2c1.pem": &fstest.MapFile{Data: pemEncode(t, key)},
		"README.txt":                              &fstest.MapFile{Data: []byte("not a key")},
		"nested":                                  &fstest.MapFile{Mode: fs.ModeDir},
	}

	ks, err := keystore.NewFS(fsys)
	if err != nil {
		t.Fatalf("NewFS() = %v, want nil", err)
	}

	got, err := ks.PrivateKey("54bb2165-71e1-41a6-af3e-7da4a0e1e2c1")
	if err != nil {
		t.Fatalf("PrivateKey() = %v, want nil", err)
	}
	if got.N.Cmp(key.N) != 0 || got.E != key.E {
		t.Errorf("PrivateKey() does not match the key stored in the PEM file")
	}

	if _, err := ks.PrivateKey("README"); err == nil {
		t.Errorf("PrivateKey(README) = nil, want an error for a non-PEM file")
	}
}

// TestNewFS_MalformedPEM verifies that NewFS returns an error when a
// .pem file cannot be parsed.
func TestNewFS_MalformedPEM(t *testing.T) {
	fsys := fstest.MapFS{
		"bad.pem": &fstest.MapFile{Data: []byte("this is not a PEM file")},
	}

	if _, err := keystore.NewFS(fsys); err == nil {
		t.Fatal("NewFS() with a malformed PEM file = nil, want an error")
	}
}
