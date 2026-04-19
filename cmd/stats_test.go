package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProfileStats(t *testing.T) {
	dir := setupTempDir(t)
	projectDir = dir

	envoyDir := filepath.Join(dir, ".envoy")
	os.MkdirAll(envoyDir, 0755)

	content := "# comment\nFOO=bar\nBAR=baz\n\nBAZ=qux\n"
	os.WriteFile(filepath.Join(envoyDir, "dev.env"), []byte(content), 0644)

	if err := profileStats("dev"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProfileStatsNonExistent(t *testing.T) {
	dir := setupTempDir(t)
	projectDir = dir

	envoyDir := filepath.Join(dir, ".envoy")
	os.MkdirAll(envoyDir, 0755)

	if err := profileStats("ghost"); err == nil {
		t.Fatal("expected error for non-existent profile")
	}
}

func TestProfileStatsEmptyFile(t *testing.T) {
	dir := setupTempDir(t)
	projectDir = dir

	envoyDir := filepath.Join(dir, ".envoy")
	os.MkdirAll(envoyDir, 0755)
	os.WriteFile(filepath.Join(envoyDir, "empty.env"), []byte(""), 0644)

	if err := profileStats("empty"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProfileStatsOnlyComments(t *testing.T) {
	dir := setupTempDir(t)
	projectDir = dir

	envoyDir := filepath.Join(dir, ".envoy")
	os.MkdirAll(envoyDir, 0755)
	os.WriteFile(filepath.Join(envoyDir, "notes.env"), []byte("# line1\n# line2\n"), 0644)

	if err := profileStats("notes"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
