package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsProfileLockedFalse(t *testing.T) {
	dir := setupLockDir(t)
	if isProfileLocked(dir, "dev") {
		t.Error("expected profile to not be locked")
	}
}

func TestIsProfileLockedTrue(t *testing.T) {
	dir := setupLockDir(t)
	_ = lockProfile(dir, "dev")
	if !isProfileLocked(dir, "dev") {
		t.Error("expected profile to be locked")
	}
}

func TestAssertNotLockedPasses(t *testing.T) {
	dir := setupLockDir(t)
	if err := assertNotLocked(dir, "dev"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAssertNotLockedFails(t *testing.T) {
	dir := setupLockDir(t)
	_ = lockProfile(dir, "dev")
	err := assertNotLocked(dir, "dev")
	if err == nil {
		t.Error("expected error for locked profile")
	}
}

func TestLockTimestamp(t *testing.T) {
	dir := setupLockDir(t)
	_ = lockProfile(dir, "dev")
	ts, err := lockTimestamp(dir, "dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ts == "" {
		t.Error("expected non-empty timestamp")
	}
}

func TestLockedProfileNames(t *testing.T) {
	dir := setupLockDir(t)
	envoyDir := filepath.Join(dir, ".envoy")
	_ = os.WriteFile(filepath.Join(envoyDir, "staging.env"), []byte("K=v\n"), 0644)
	_ = lockProfile(dir, "dev")
	_ = lockProfile(dir, "staging")
	names, err := lockedProfileNames(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(names) != 2 {
		t.Fatalf("expected 2 locked profiles, got %d", len(names))
	}
	all := strings.Join(names, ",")
	if !strings.Contains(all, "dev") || !strings.Contains(all, "staging") {
		t.Errorf("unexpected names: %v", names)
	}
}

func TestLockedProfileNamesNoDir(t *testing.T) {
	dir := t.TempDir()
	names, err := lockedProfileNames(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(names) != 0 {
		t.Errorf("expected no names, got %v", names)
	}
}
