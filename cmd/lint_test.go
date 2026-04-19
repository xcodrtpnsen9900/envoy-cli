package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func writeLintProfile(t *testing.T, dir, name, content string) {
	t.Helper()
	profilesDir := filepath.Join(dir, ".envoy", "profiles")
	os.MkdirAll(profilesDir, 0755)
	os.WriteFile(filepath.Join(profilelesDir, name+".env"), []byte(content), 0644)
}

func TestLintValidProfile(t *testing.T) {
	dir := setupTempDir(t)
	writeProfile(t, dir, "dev", "KEY=value\nFOO=bar\n")
	issues, err := lintProfile(dir, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got: %v", issues)
	}
}

func TestLintMissingEquals(t *testing.T) {
	dir := setupTempDir(t)
	writeProfile(t, dir, "dev", "KEYONLY\nFOO=bar\n")
	issues, err := lintProfile(dir, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d: %v", len(issues), issues)
	}
}

func TestLintWhitespaceInKey(t *testing.T) {
	dir := setupTempDir(t)
	writeProfile(t, dir, "dev", "MY KEY=value\n")
	issues, err := lintProfile(dir, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) == 0 {
		t.Fatal("expected issue for whitespace in key")
	}
}

func TestLintLeadingWhitespaceInValue(t *testing.T) {
	dir := setupTempDir(t)
	writeProfile(t, dir, "dev", "KEY= value\n")
	issues, err := lintProfile(dir, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) == 0 {
		t.Fatal("expected issue for leading whitespace in value")
	}
}

func TestLintNonExistentProfile(t *testing.T) {
	dir := setupTempDir(t)
	_, err := lintProfile(dir, "ghost")
	if err == nil {
		t.Fatal("expected error for non-existent profile")
	}
}

func TestLintSkipsComments(t *testing.T) {
	dir := setupTempDir(t)
	writeProfile(t, dir, "dev", "# this is a comment\nKEY=val\n")
	issues, err := lintProfile(dir, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got: %v", issues)
	}
}
