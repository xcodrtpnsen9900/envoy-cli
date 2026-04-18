package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCopyProfile(t *testing.T) {
	dir := setupTempDir(t)

	if err := addProfile(dir, "staging", map[string]string{"APP_ENV": "staging", "PORT": "8080"}); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	if err := copyProfile(dir, "staging", "staging-backup"); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	profiles, err := listProfiles(dir)
	if err != nil {
		t.Fatalf("listProfiles failed: %v", err)
	}

	found := false
	for _, p := range profiles {
		if p == "staging-backup" {
			found = true
		}
	}
	if !found {
		t.Error("copied profile 'staging-backup' not found in list")
	}

	srcContent, _ := os.ReadFile(profilePath(dir, "staging"))
	dstContent, _ := os.ReadFile(profilePath(dir, "staging-backup"))
	if string(srcContent) != string(dstContent) {
		t.Error("copied profile content does not match source")
	}
}

func TestCopyNonExistentProfile(t *testing.T) {
	dir := setupTempDir(t)

	err := copyProfile(dir, "ghost", "ghost-copy")
	if err == nil {
		t.Error("expected error when copying non-existent profile")
	}
}

func TestCopyProfileAlreadyExists(t *testing.T) {
	dir := setupTempDir(t)

	for _, name := range []string{"prod", "prod-copy"} {
		if err := addProfile(dir, name, map[string]string{"ENV": name}); err != nil {
			t.Fatalf("setup failed: %v", err)
		}
	}

	err := copyProfile(dir, "prod", "prod-copy")
	if err == nil {
		t.Error("expected error when destination profile already exists")
	}
}

func TestCopyFileCreatesDestDir(t *testing.T) {
	dir := setupTempDir(t)

	src := filepath.Join(dir, "source.env")
	if err := os.WriteFile(src, []byte("KEY=value"), 0644); err != nil {
		t.Fatalf("failed to write source: %v", err)
	}

	dst := filepath.Join(dir, "nested", "dest.env")
	if err := copyFile(src, dst); err != nil {
		t.Fatalf("copyFile failed: %v", err)
	}

	if _, err := os.Stat(dst); os.IsNotExist(err) {
		t.Error("destination file was not created")
	}
}
