package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSearchProfileFindsKey(t *testing.T) {
	dir := t.TempDir()
	envoyDir := filepath.Join(dir, ".envoy")
	os.MkdirAll(envoyDir, 0755)
	os.WriteFile(filepath.Join(envoyDir, "dev.env"), []byte("DB_HOST=localhost\nDB_PORT=5432\nAPI_KEY=secret\n"), 0644)

	results, err := searchProfile(dir, "dev", "DB", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestSearchProfileCaseInsensitive(t *testing.T) {
	dir := t.TempDir()
	envoyDir := filepath.Join(dir, ".envoy")
	os.MkdirAll(envoyDir, 0755)
	os.WriteFile(filepath.Join(envoyDir, "dev.env"), []byte("DB_HOST=localhost\nAPI_KEY=secret\n"), 0644)

	results, err := searchProfile(dir, "dev", "db_host", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestSearchProfileCaseSensitiveNoMatch(t *testing.T) {
	dir := t.TempDir()
	envoyDir := filepath.Join(dir, ".envoy")
	os.MkdirAll(envoyDir, 0755)
	os.WriteFile(filepath.Join(envoyDir, "dev.env"), []byte("DB_HOST=localhost\n"), 0644)

	results, err := searchProfile(dir, "dev", "db_host", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestSearchProfileSkipsComments(t *testing.T) {
	dir := t.TempDir()
	envoyDir := filepath.Join(dir, ".envoy")
	os.MkdirAll(envoyDir, 0755)
	os.WriteFile(filepath.Join(envoyDir, "dev.env"), []byte("# DB_HOST comment\nDB_HOST=localhost\n"), 0644)

	results, err := searchProfile(dir, "dev", "DB_HOST", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result (comment skipped), got %d", len(results))
	}
}

func TestSearchNonExistentProfile(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".envoy"), 0755)

	_, err := searchProfile(dir, "ghost", "KEY", false)
	if err == nil {
		t.Fatal("expected error for non-existent profile")
	}
}
