package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func runAuditCommand(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("go", append([]string{"run", "../main.go"}, args...)...)
	cmd.Env = append(os.Environ(), "ENVOY_PROJECT_DIR="+dir)
	out, _ := cmd.CombinedOutput()
	return string(out)
}

func TestAuditEmptyLog(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".envoy"), 0755); err != nil {
		t.Fatal(err)
	}
	out := runAuditCommand(t, dir, "audit")
	if !strings.Contains(out, "No audit log") {
		t.Errorf("expected empty message, got: %s", out)
	}
}

func TestAuditShowsEntries(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".envoy"), 0755); err != nil {
		t.Fatal(err)
	}
	writeAuditEntry(dir, "switch", "prod")
	writeAuditEntry(dir, "delete", "old")

	entries, err := readAuditLog(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if !strings.Contains(entries[1], "delete") {
		t.Errorf("expected delete action in second entry: %s", entries[1])
	}
}
