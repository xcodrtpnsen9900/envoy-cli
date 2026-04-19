package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestLintProfileViaCommand(t *testing.T) {
	dir := setupTempDir(t)
	writeProfile(t, dir, "prod", "HOST=localhost\nPORT=8080\n")

	origDir, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(origDir)

	issues, err := lintProfile(dir, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 0 {
		t.Fatalf("expected clean profile, got issues: %v", issues)
	}
}

func TestLintOutputContainsLineNumbers(t *testing.T) {
	dir := setupTempDir(t)
	writeProfile(t, dir, "bad", "GOOD=ok\nBADLINE\nANOTHER=fine\n")

	issues, err := lintProfile(dir, "bad")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) == 0 {
		t.Fatal("expected at least one issue")
	}
	var buf bytes.Buffer
	for _, iss := range issues {
		buf.WriteString(iss)
	}
	if !bytes.Contains(buf.Bytes(), []byte("line 2")) {
		t.Fatalf("expected line number in output, got: %s", buf.String())
	}
}

func TestLintProfileNotInEnvoyDir(t *testing.T) {
	dir := setupTempDir(t)
	// Write file outside profiles dir — should not be found
	os.WriteFile(filepath.Join(dir, "stray.env"), []byte("KEY=val\n"), 0644)
	_, err := lintProfile(dir, "stray")
	if err == nil {
		t.Fatal("expected error when profile not in .envoy/profiles")
	}
}
