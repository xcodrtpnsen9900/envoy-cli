package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCompareProfilesOnlyInA(t *testing.T) {
	dir := setupTempDir(t)
	projectDir = dir

	envoyDir := filepath.Join(dir, ".envoy")
	os.MkdirAll(envoyDir, 0755)

	os.WriteFile(filepath.Join(envoyDir, "dev.env"), []byte("FOO=1\nBAR=2\n"), 0644)
	os.WriteFile(filepath.Join(envoyDir, "prod.env"), []byte("FOO=1\n"), 0644)

	if err := compareProfiles("dev", "prod"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCompareProfilesIdentical(t *testing.T) {
	dir := setupTempDir(t)
	projectDir = dir

	envoyDir := filepath.Join(dir, ".envoy")
	os.MkdirAll(envoyDir, 0755)

	content := []byte("FOO=1\nBAR=2\n")
	os.WriteFile(filepath.Join(envoyDir, "a.env"), content, 0644)
	os.WriteFile(filepath.Join(envoyDir, "b.env"), content, 0644)

	if err := compareProfiles("a", "b"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCompareNonExistentProfile(t *testing.T) {
	dir := setupTempDir(t)
	projectDir = dir

	envoyDir := filepath.Join(dir, ".envoy")
	os.MkdirAll(envoyDir, 0755)
	os.WriteFile(filepath.Join(envoyDir, "dev.env"), []byte("FOO=1\n"), 0644)

	if err := compareProfiles("dev", "ghost"); err == nil {
		t.Fatal("expected error for non-existent profile")
	}
}

func TestKeysOnlyIn(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2"}
	b := map[string]string{"B": "2", "C": "3"}

	only := keysOnlyIn(a, b)
	if len(only) != 1 || only[0] != "A" {
		t.Fatalf("expected [A], got %v", only)
	}
}

func TestSharedKeys(t *testing.T) {
	a := map[string]string{"A": "1", "B": "2"}
	b := map[string]string{"B": "2", "C": "3"}

	shared := sharedKeys(a, b)
	if len(shared) != 1 || shared[0] != "B" {
		t.Fatalf("expected [B], got %v", shared)
	}
}
