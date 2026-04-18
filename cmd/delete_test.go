package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDeleteProfile(t *testing.T) {
	dir := setupTempDir(t)

	if err := initProject(dir); err != nil {
		t.Fatalf("initProject failed: %v", err)
	}
	if err := addProfile(dir, "staging", map[string]string{"APP_ENV": "staging"}); err != nil {
		t.Fatalf("addProfile failed: %v", err)
	}
	if err := switchProfile(dir, "staging"); err != nil {
		t.Fatalf("switchProfile failed: %v", err)
	}
	if err := addProfile(dir, "production", map[string]string{"APP_ENV": "production"}); err != nil {
		t.Fatalf("addProfile failed: %v", err)
	}

	// Cannot delete active profile
	if err := deleteProfile(dir, "staging"); err == nil {
		t.Error("expected error when deleting active profile, got nil")
	}

	// Can delete non-active profile
	if err := deleteProfile(dir, "production"); err != nil {
		t.Errorf("unexpected error deleting non-active profile: %v", err)
	}

	profiles, err := listProfiles(dir)
	if err != nil {
		t.Fatalf("listProfiles failed: %v", err)
	}
	for _, p := range profiles {
		if p == "production" {
			t.Error("expected 'production' to be deleted")
		}
	}
}

func TestDeleteNonExistentProfile(t *testing.T) {
	dir := setupTempDir(t)

	if err := initProject(dir); err != nil {
		t.Fatalf("initProject failed: %v", err)
	}

	if err := deleteProfile(dir, "ghost"); err == nil {
		t.Error("expected error deleting non-existent profile, got nil")
	}
}

func TestDeleteProfileRemovesFile(t *testing.T) {
	dir := setupTempDir(t)

	if err := initProject(dir); err != nil {
		t.Fatalf("initProject failed: %v", err)
	}
	if err := addProfile(dir, "dev", map[string]string{"APP_ENV": "dev"}); err != nil {
		t.Fatalf("addProfile failed: %v", err)
	}
	if err := addProfile(dir, "staging", map[string]string{"APP_ENV": "staging"}); err != nil {
		t.Fatalf("addProfile failed: %v", err)
	}
	if err := switchProfile(dir, "staging"); err != nil {
		t.Fatalf("switchProfile failed: %v", err)
	}
	if err := deleteProfile(dir, "dev"); err != nil {
		t.Fatalf("deleteProfile failed: %v", err)
	}

	path := filepath.Join(dir, ".envoy", "profiles", "dev.env")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("expected profile file to be removed from disk")
	}
}
