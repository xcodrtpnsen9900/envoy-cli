package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func setupLockDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	envoyDir := filepath.Join(dir, ".envoy")
	if err := os.MkdirAll(envoyDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(envoyDir, "active"), []byte("dev"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(envoyDir, "dev.env"), []byte("KEY=value\n"), 0644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestLockProfile(t *testing.T) {
	dir := setupLockDir(t)
	if err := lockProfile(dir, "dev"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !isProfileLocked(dir, "dev") {
		t.Error("expected profile to be locked")
	}
}

func TestLockNonExistentProfile(t *testing.T) {
	dir := setupLockDir(t)
	err := lockProfile(dir, "ghost")
	if err == nil {
		t.Error("expected error for non-existent profile")
	}
}

func TestLockAlreadyLocked(t *testing.T) {
	dir := setupLockDir(t)
	_ = lockProfile(dir, "dev")
	err := lockProfile(dir, "dev")
	if err == nil {
		t.Error("expected error when locking already-locked profile")
	}
}

func TestUnlockProfile(t *testing.T) {
	dir := setupLockDir(t)
	_ = lockProfile(dir, "dev")
	if err := unlockProfile(dir, "dev"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if isProfileLocked(dir, "dev") {
		t.Error("expected profile to be unlocked")
	}
}

func TestUnlockNotLocked(t *testing.T) {
	dir := setupLockDir(t)
	err := unlockProfile(dir, "dev")
	if err == nil {
		t.Error("expected error when unlocking non-locked profile")
	}
}

func TestListLockedProfilesEmpty(t *testing.T) {
	dir := setupLockDir(t)
	if err := listLockedProfiles(dir); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
