package cmd

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func setupExpireDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	envoyDir := filepath.Join(dir, ".envoy")
	if err := os.MkdirAll(envoyDir, 0755); err != nil {
		t.Fatal(err)
	}
	// create a test profile
	if err := os.WriteFile(filepath.Join(envoyDir, "dev.env"), []byte("KEY=value\n"), 0644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestSetProfileExpiry(t *testing.T) {
	dir := setupExpireDir(t)
	if err := setProfileExpiry(dir, "dev", "1h"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries, err := loadExpiry(dir)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := entries["dev"]; !ok {
		t.Error("expected expiry entry for 'dev'")
	}
}

func TestSetProfileExpiryNonExistent(t *testing.T) {
	dir := setupExpireDir(t)
	err := setProfileExpiry(dir, "ghost", "1h")
	if err == nil {
		t.Error("expected error for non-existent profile")
	}
}

func TestSetProfileExpiryInvalidDuration(t *testing.T) {
	dir := setupExpireDir(t)
	err := setProfileExpiry(dir, "dev", "notaduration")
	if err == nil {
		t.Error("expected error for invalid duration")
	}
}

func TestCheckProfileExpiryNoEntries(t *testing.T) {
	dir := setupExpireDir(t)
	if err := checkProfileExpiry(dir); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIsExpiredFalse(t *testing.T) {
	dir := setupExpireDir(t)
	_ = setProfileExpiry(dir, "dev", "1h")
	expired, err := isExpired(dir, "dev")
	if err != nil {
		t.Fatal(err)
	}
	if expired {
		t.Error("expected profile to not be expired")
	}
}

func TestIsExpiredTrue(t *testing.T) {
	dir := setupExpireDir(t)
	entries := map[string]ExpiryEntry{
		"dev": {Profile: "dev", ExpiresAt: time.Now().Add(-1 * time.Hour)},
	}
	if err := saveExpiry(dir, entries); err != nil {
		t.Fatal(err)
	}
	expired, err := isExpired(dir, "dev")
	if err != nil {
		t.Fatal(err)
	}
	if !expired {
		t.Error("expected profile to be expired")
	}
}

func TestIsExpiredNoEntry(t *testing.T) {
	dir := setupExpireDir(t)
	expired, err := isExpired(dir, "dev")
	if err != nil {
		t.Fatal(err)
	}
	if expired {
		t.Error("expected false when no expiry entry")
	}
}
