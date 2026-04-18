package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestActiveProfileAfterInit(t *testing.T) {
	dir := setupTempDir(t)

	if err := initProject(dir); err != nil {
		t.Fatalf("init failed: %v", err)
	}

	active, err := activeProfile(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if active != "default" {
		t.Errorf("expected 'default', got '%s'", active)
	}
}

func TestActiveProfileAfterSwitch(t *testing.T) {
	dir := setupTempDir(t)

	if err := initProject(dir); err != nil {
		t.Fatalf("init failed: %v", err)
	}
	if err := addProfile(dir, "staging"); err != nil {
		t.Fatalf("addProfile failed: %v", err)
	}
	if err := switchProfile(dir, "staging"); err != nil {
		t.Fatalf("switchProfile failed: %v", err)
	}

	active, err := activeProfile(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if active != "staging" {
		t.Errorf("expected 'staging', got '%s'", active)
	}
}

func TestActiveProfileEmptyFile(t *testing.T) {
	dir := setupTempDir(t)

	envoyDir := filepath.Join(dir, ".envoy")
	if err := os.MkdirAll(envoyDir, 0755); err != nil {
		t.Fatal(err)
	}
	activeFile := filepath.Join(envoyDir, ".active")
	if err := os.WriteFile(activeFile, []byte(""), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := activeProfile(dir)
	if err == nil {
		t.Error("expected error for empty active file, got nil")
	}
}
