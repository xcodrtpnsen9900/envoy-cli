package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRenameProfile(t *testing.T) {
	dir := setupTempDir(t)
	projectDir = dir

	if err := addProfile("staging", []string{}); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if err := renameProfile("staging", "uat"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	profiles, _ := listProfiles()
	for _, p := range profiles {
		if p == "staging" {
			t.Error("old profile name should not exist")
		}
	}

	found := false
	for _, p := range profiles {
		if p == "uat" {
			found = true
		}
	}
	if !found {
		t.Error("new profile name should exist")
	}
}

func TestRenameNonExistentProfile(t *testing.T) {
	dir := setupTempDir(t)
	projectDir = dir

	if err := renameProfile("ghost", "phantom"); err == nil {
		t.Error("expected error renaming non-existent profile")
	}
}

func TestRenameToExistingProfile(t *testing.T) {
	dir := setupTempDir(t)
	projectDir = dir

	addProfile("alpha", []string{})
	addProfile("beta", []string{})

	if err := renameProfile("alpha", "beta"); err == nil {
		t.Error("expected error when renaming to existing profile name")
	}
}

func TestRenameUpdatesActiveProfile(t *testing.T) {
	dir := setupTempDir(t)
	projectDir = dir

	addProfile("dev", []string{})
	switchProfile("dev")

	if err := renameProfile("dev", "development"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	active, err := activeProfile()
	if err != nil {
		t.Fatalf("unexpected error reading active profile: %v", err)
	}
	if active != "development" {
		t.Errorf("expected active profile to be 'development', got '%s'", active)
	}

	envPath := filepath.Join(dir, ".env")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		t.Error("expected .env symlink to still exist after rename")
	}
}
