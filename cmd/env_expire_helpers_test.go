package cmd

import (
	"testing"
	"time"
)

func TestClearExpiry(t *testing.T) {
	dir := setupExpireDir(t)
	_ = setProfileExpiry(dir, "dev", "1h")
	if err := clearExpiry(dir, "dev"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	count, _ := expiryCount(dir)
	if count != 0 {
		t.Errorf("expected 0 entries after clear, got %d", count)
	}
}

func TestClearExpiryNotSet(t *testing.T) {
	dir := setupExpireDir(t)
	err := clearExpiry(dir, "dev")
	if err == nil {
		t.Error("expected error when clearing unset expiry")
	}
}

func TestExpiredProfiles(t *testing.T) {
	dir := setupExpireDir(t)
	entries := map[string]ExpiryEntry{
		"dev":  {Profile: "dev", ExpiresAt: time.Now().Add(-1 * time.Hour)},
		"prod": {Profile: "prod", ExpiresAt: time.Now().Add(1 * time.Hour)},
	}
	_ = saveExpiry(dir, entries)
	expired, err := expiredProfiles(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(expired) != 1 || expired[0] != "dev" {
		t.Errorf("expected [dev], got %v", expired)
	}
}

func TestExpiryCount(t *testing.T) {
	dir := setupExpireDir(t)
	_ = setProfileExpiry(dir, "dev", "1h")
	count, err := expiryCount(dir)
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("expected count 1, got %d", count)
	}
}

func TestTimeUntilExpiry(t *testing.T) {
	dir := setupExpireDir(t)
	_ = setProfileExpiry(dir, "dev", "2h")
	remaining, err := timeUntilExpiry(dir, "dev")
	if err != nil {
		t.Fatal(err)
	}
	if remaining <= 0 || remaining > 2*time.Hour+time.Second {
		t.Errorf("unexpected remaining duration: %v", remaining)
	}
}

func TestTimeUntilExpiryNotSet(t *testing.T) {
	dir := setupExpireDir(t)
	_, err := timeUntilExpiry(dir, "dev")
	if err == nil {
		t.Error("expected error when no expiry is set")
	}
}
