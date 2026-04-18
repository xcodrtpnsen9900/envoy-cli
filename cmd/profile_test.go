package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTempDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "envoy-test-*")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	return dir
}

func TestAddAndListProfiles(t *testing.T) {
	dir := setupTempDir(t)

	if err := addProfile(dir, "development"); err != nil {
		t.Fatalf("addProfile failed: %v", err)
	}
	if err := addProfile(dir, "production"); err != nil {
		t.Fatalf("addProfile failed: %v", err)
	}

	profiles, err := listProfiles(dir)
	if err != nil {
		t.Fatalf("listProfiles failed: %v", err)
	}
	if len(profiles) != 2 {
		t.Fatalf("expected 2 profiles, got %d", len(profiles))
	}
}

func TestAddDuplicateProfile(t *testing.T) {
	dir := setupTempDir(t)
	_ = addProfile(dir, "staging")
	if err := addProfile(dir, "staging"); err == nil {
		t.Fatal("expected error for duplicate profile, got nil")
	}
}

func TestSwitchProfile(t *testing.T) {
	dir := setupTempDir(t)

	envContent := []byte("APP_ENV=development\nDEBUG=true\n")
	src := filepath.Join(dir, activeLink)
	if err := os.WriteFile(src, envContent, 0644); err != nil {
		t.Fatal(err)
	}

	if err := addProfile(dir, "dev"); err != nil {
		t.Fatal(err)
	}
	if err := switchProfile(dir, "dev"); err != nil {
		t.Fatalf("switchProfile failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, activeLink))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != string(envContent) {
		t.Fatalf("active .env content mismatch")
	}
}

func TestSwitchNonExistentProfile(t *testing.T) {
	dir := setupTempDir(t)
	if err := switchProfile(dir, "ghost"); err == nil {
		t.Fatal("expected error for missing profile, got nil")
	}
}
