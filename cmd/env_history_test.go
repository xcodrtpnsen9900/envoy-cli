package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func setupHistoryDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".envoy"), 0755); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestRecordAndReadHistory(t *testing.T) {
	dir := setupHistoryDir(t)

	if err := recordSwitchHistory(dir, "dev"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := recordSwitchHistory(dir, "staging"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, err := showHistory(dir, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	// most recent first
	if !strings.Contains(entries[0], "staging") {
		t.Errorf("expected first entry to contain 'staging', got %q", entries[0])
	}
}

func TestHistoryLimit(t *testing.T) {
	dir := setupHistoryDir(t)

	for _, p := range []string{"dev", "staging", "prod", "test", "local"} {
		if err := recordSwitchHistory(dir, p); err != nil {
			t.Fatal(err)
		}
	}

	entries, err := showHistory(dir, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("expected 3 entries with limit, got %d", len(entries))
	}
}

func TestHistoryEmptyFile(t *testing.T) {
	dir := setupHistoryDir(t)
	entries, err := showHistory(dir, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestHistoryEntryFormat(t *testing.T) {
	dir := setupHistoryDir(t)
	before := time.Now()
	if err := recordSwitchHistory(dir, "production"); err != nil {
		t.Fatal(err)
	}
	after := time.Now()

	entries, _ := showHistory(dir, 1)
	if len(entries) != 1 {
		t.Fatal("expected 1 entry")
	}
	parts := strings.SplitN(entries[0], " ", 2)
	if len(parts) != 2 {
		t.Fatalf("unexpected format: %q", entries[0])
	}
	ts, err := time.Parse(time.RFC3339, parts[0])
	if err != nil {
		t.Fatalf("invalid timestamp: %v", err)
	}
	if ts.Before(before) || ts.After(after.Add(time.Second)) {
		t.Errorf("timestamp out of range: %v", ts)
	}
	if parts[1] != "production" {
		t.Errorf("expected profile 'production', got %q", parts[1])
	}
}

func TestHistoryMaxEntries(t *testing.T) {
	dir := setupHistoryDir(t)
	for i := 0; i < maxHistoryEntries+10; i++ {
		if err := recordSwitchHistory(dir, "profile"); err != nil {
			t.Fatal(err)
		}
	}
	entries, err := showHistory(dir, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) > maxHistoryEntries {
		t.Errorf("expected at most %d entries, got %d", maxHistoryEntries, len(entries))
	}
}
