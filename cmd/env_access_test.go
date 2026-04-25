package cmd

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func setupAccessDir(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".envoy"), 0755); err != nil {
		t.Fatal(err)
	}
	return root
}

func TestRecordAndReadAccessLog(t *testing.T) {
	root := setupAccessDir(t)
	if err := recordAccess(root, "production", "read"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := recordAccess(root, "staging", "switch"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries, err := readAccessLog(root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Profile != "production" || entries[0].Action != "read" {
		t.Errorf("unexpected first entry: %+v", entries[0])
	}
}

func TestReadAccessLogEmpty(t *testing.T) {
	root := setupAccessDir(t)
	entries, err := readAccessLog(root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty log, got %d entries", len(entries))
	}
}

func TestAccessCount(t *testing.T) {
	root := setupAccessDir(t)
	_ = recordAccess(root, "dev", "read")
	_ = recordAccess(root, "dev", "switch")
	_ = recordAccess(root, "prod", "read")
	if got := accessCount(root); got != 3 {
		t.Errorf("expected 3, got %d", got)
	}
	if got := accessCountForProfile(root, "dev"); got != 2 {
		t.Errorf("expected 2 for dev, got %d", got)
	}
}

func TestLastAccessEntry(t *testing.T) {
	root := setupAccessDir(t)
	if e := lastAccessEntry(root); e != nil {
		t.Errorf("expected nil for empty log")
	}
	_ = recordAccess(root, "alpha", "read")
	time.Sleep(2 * time.Millisecond)
	_ = recordAccess(root, "beta", "switch")
	e := lastAccessEntry(root)
	if e == nil {
		t.Fatal("expected entry, got nil")
	}
	if e.Profile != "beta" {
		t.Errorf("expected last entry profile=beta, got %s", e.Profile)
	}
}

func TestClearAccessLog(t *testing.T) {
	root := setupAccessDir(t)
	_ = recordAccess(root, "dev", "read")
	if err := clearAccessLog(root); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := accessCount(root); got != 0 {
		t.Errorf("expected 0 after clear, got %d", got)
	}
}
