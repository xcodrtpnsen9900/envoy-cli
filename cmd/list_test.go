package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListProfilesEmpty(t *testing.T) {
	dir := setupTempDir(t)
	profilesDir := filepath.Join(dir, ".envoy", "profiles")
	if err := os.MkdirAll(profilesDir, 0755); err != nil {
		t.Fatal(err)
	}
	profiles, err := listProfilesDetailed(dir, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 0 {
		t.Errorf("expected 0 profiles, got %d", len(profiles))
	}
}

func TestListProfilesReturnsNames(t *testing.T) {
	dir := setupTempDir(t)
	profilesDir := filepath.Join(dir, ".envoy", "profiles")
	if err := os.MkdirAll(profilesDir, 0755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"dev", "staging", "prod"} {
		if err := os.WriteFile(filepath.Join(profilesDir, name), []byte("KEY=val\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	profiles, err := listProfilesDetailed(dir, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 3 {
		t.Errorf("expected 3 profiles, got %d", len(profiles))
	}
}

func TestListProfilesVerbose(t *testing.T) {
	dir := setupTempDir(t)
	profilesDir := filepath.Join(dir, ".envoy", "profiles")
	if err := os.MkdirAll(profilesDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(profilesDir, "dev"), []byte("KEY=value\n"), 0644); err != nil {
		t.Fatal(err)
	}
	profiles, err := listProfilesDetailed(dir, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 1 {
		t.Fatalf("expected 1 profile, got %d", len(profiles))
	}
	if profiles[0].size == 0 {
		t.Error("expected non-zero size in verbose mode")
	}
	if profiles[0].modTime == "" {
		t.Error("expected non-empty modTime in verbose mode")
	}
}

func TestListProfilesNoEnvoyDir(t *testing.T) {
	dir := setupTempDir(t)
	profiles, err := listProfilesDetailed(dir, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if profiles != nil {
		t.Errorf("expected nil profiles when dir missing, got %v", profiles)
	}
}
