package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitProject(t *testing.T) {
	dir := setupTempDir(t)

	if err := initProject(dir); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	envoyDir := filepath.Join(dir, ".envoy")
	if _, err := os.Stat(envoyDir); os.IsNotExist(err) {
		t.Error(".envoy directory was not created")
	}

	defaultEnv := filepath.Join(envoyDir, "default.env")
	if _, err := os.Stat(defaultEnv); os.IsNotExist(err) {
		t.Error("default.env profile was not created")
	}

	active, err := activeProfile(dir)
	if err != nil {
		t.Fatalf("expected active profile, got error: %v", err)
	}
	if active != "default" {
		t.Errorf("expected active profile 'default', got '%s'", active)
	}
}

func TestInitProjectAlreadyExists(t *testing.T) {
	dir := setupTempDir(t)

	if err := initProject(dir); err != nil {
		t.Fatalf("first init failed: %v", err)
	}
	if err := initProject(dir); err == nil {
		t.Error("expected error on second init, got nil")
	}
}

func TestActiveProfileMissing(t *testing.T) {
	dir := setupTempDir(t)

	_, err := activeProfile(dir)
	if err == nil {
		t.Error("expected error when .active file missing, got nil")
	}
}
