package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckRequiredKeysAllPresent(t *testing.T) {
	dir := setupTempDir(t)
	projectDir = dir
	envoyDir := filepath.Join(dir, ".envoy")
	os.MkdirAll(envoyDir, 0755)

	content := "DB_HOST=localhost\nDB_PORT=5432\nAPI_KEY=secret\n"
	os.WriteFile(filepath.Join(envoyDir, "prod.env"), []byte(content), 0644)

	missing, err := checkRequiredKeys("prod", []string{"DB_HOST", "DB_PORT", "API_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(missing) != 0 {
		t.Errorf("expected no missing keys, got %v", missing)
	}
}

func TestCheckRequiredKeysMissing(t *testing.T) {
	dir := setupTempDir(t)
	projectDir = dir
	envoyDir := filepath.Join(dir, ".envoy")
	os.MkdirAll(envoyDir, 0755)

	content := "DB_HOST=localhost\n"
	os.WriteFile(filepath.Join(envoyDir, "staging.env"), []byte(content), 0644)

	missing, err := checkRequiredKeys("staging", []string{"DB_HOST", "DB_PORT", "API_KEY"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(missing) != 2 {
		t.Errorf("expected 2 missing keys, got %v", missing)
	}
}

func TestCheckRequiredKeysNonExistentProfile(t *testing.T) {
	dir := setupTempDir(t)
	projectDir = dir
	envoyDir := filepath.Join(dir, ".envoy")
	os.MkdirAll(envoyDir, 0755)

	_, err := checkRequiredKeys("ghost", []string{"KEY"})
	if err == nil {
		t.Fatal("expected error for non-existent profile")
	}
}

func TestCheckRequiredKeysSkipsComments(t *testing.T) {
	dir := setupTempDir(t)
	projectDir = dir
	envoyDir := filepath.Join(dir, ".envoy")
	os.MkdirAll(envoyDir, 0755)

	content := "# DB_PORT=5432\nDB_HOST=localhost\n"
	os.WriteFile(filepath.Join(envoyDir, "dev.env"), []byte(content), 0644)

	missing, err := checkRequiredKeys("dev", []string{"DB_PORT"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(missing) != 1 {
		t.Errorf("expected DB_PORT to be missing (it's a comment), got %v", missing)
	}
}
