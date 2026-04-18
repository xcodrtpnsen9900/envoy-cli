package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestShowProfile(t *testing.T) {
	dir := setupTempDir(t)
	profDir := filepath.Join(dir, ".envoy", "profiles")
	if err := os.MkdirAll(profDir, 0755); err != nil {
		t.Fatal(err)
	}
	content := "KEY=value\nFOO=bar\n"
	if err := os.WriteFile(filepath.Join(profDir, "dev.env"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	if err := showProfile("dev"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestShowNonExistentProfile(t *testing.T) {
	dir := setupTempDir(t)
	if err := os.MkdirAll(filepath.Join(dir, ".envoy", "profiles"), 0755); err != nil {
		t.Fatal(err)
	}
	err := showProfile("ghost")
	if err == nil {
		t.Fatal("expected error for non-existent profile")
	}
}

func TestShowProfileEmptyFile(t *testing.T) {
	dir := setupTempDir(t)
	profDir := filepath.Join(dir, ".envoy", "profiles")
	if err := os.MkdirAll(profDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(profDir, "empty.env"), []byte{}, 0644); err != nil {
		t.Fatal(err)
	}
	if err := showProfile("empty"); err != nil {
		t.Fatalf("expected no error for empty profile, got %v", err)
	}
}
