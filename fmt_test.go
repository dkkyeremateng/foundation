package foundation_test

import (
	"os/exec"
	"strings"
	"testing"
)

// TestGofmt fails when any Go file in the repository is not gofmt-formatted.
//
// It exists because nothing else here could see formatting. The factory's
// verification allowlist is `go build`, `go vet`, `go test`, and unformatted Go
// passes all three — so agent-written code with statements at column 0 was
// merged twice without anything objecting.
//
// The obvious fix, adding `gofmt -l .` to the allowlist, does not work:
// `gofmt -l` LISTS offending files and exits 0, so it can never fail a build.
// The allowlist also runs commands as argv rather than through a shell, so the
// usual `test -z "$(gofmt -l .)"` cannot be expressed either. A test is the one
// place that can turn that listing into a failure.
func TestGofmt(t *testing.T) {
	out, err := exec.Command("gofmt", "-l", ".").CombinedOutput()
	if err != nil {
		t.Fatalf("running gofmt: %v\n%s", err, out)
	}
	if files := strings.TrimSpace(string(out)); files != "" {
		t.Errorf("these files are not gofmt-formatted; run `gofmt -w .`:\n%s", files)
	}
}
