package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCloneProfileSameDir(t *testing.T) {
	dir := setupTempDir(t)
	envoyDir := filepath.Join(dir, ".envoy")
	os.MkdirAll(envoyDir, 0755)

	src := filepath.Join(envoyDir, "staging.env")
	os.WriteFile(src, []byte("KEY=value\n"), 0644)

	if err := cloneProfile("staging", "staging-copy", ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dest := filepath.Join(envoyDir, "staging-copy.env")
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		t.Fatal("cloned profile file does not exist")
	}

	data, _ := os.ReadFile(dest)
	if string(data) != "KEY=value\n" {
		t.Errorf("expected cloned content to match, got %q", string(data))
	}
}

func TestCloneNonExistentProfile(t *testing.T) {
	dir := setupTempDir(t)
	envoyDir := filepath.Join(dir, ".envoy")
	os.MkdirAll(envoyDir, 0755)

	err := cloneProfile("ghost", "ghost-copy", "")
	if err == nil {
		t.Fatal("expected error for non-existent profile")
	}
}

func TestCloneProfileAlreadyExists(t *testing.T) {
	dir := setupTempDir(t)
	envoyDir := filepath.Join(dir, ".envoy")
	os.MkdirAll(envoyDir, 0755)

	os.WriteFile(filepath.Join(envoyDir, "prod.env"), []byte("A=1\n"), 0644)
	os.WriteFile(filepath.Join(envoyDir, "prod-clone.env"), []byte("B=2\n"), 0644)

	err := cloneProfile("prod", "prod-clone", "")
	if err == nil {
		t.Fatal("expected error when destination profile already exists")
	}
}

func TestCloneProfileToTargetDir(t *testing.T) {
	dir := setupTempDir(t)
	envoyDir := filepath.Join(dir, ".envoy")
	os.MkdirAll(envoyDir, 0755)
	os.WriteFile(filepath.Join(envoyDir, "dev.env"), []byte("ENV=dev\n"), 0644)

	targetDir := t.TempDir()
	if err := cloneProfile("dev", "dev", targetDir); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dest := filepath.Join(targetDir, ".envoy", "dev.env")
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		t.Fatal("expected cloned profile in target dir")
	}
}
