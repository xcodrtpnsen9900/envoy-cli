package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSnapshotProfile(t *testing.T) {
	dir := setupTempDir(t)
	defer os.RemoveAll(dir)

	if err := addProfile("dev", dir); err != nil {
		t.Fatal(err)
	}
	p := profilePath("dev")
	os.WriteFile(p, []byte("KEY=value\n"), 0644)

	if err := snapshotProfile("dev"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	entries, err := os.ReadDir(snapshotsDir("dev"))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 snapshot, got %d", len(entries))
	}
}

func TestSnapshotNonExistentProfile(t *testing.T) {
	dir := setupTempDir(t)
	defer os.RemoveAll(dir)

	err := snapshotProfile("ghost")
	if err == nil {
		t.Fatal("expected error for non-existent profile")
	}
	_ = dir
}

func TestListSnapshotsEmpty(t *testing.T) {
	dir := setupTempDir(t)
	defer os.RemoveAll(dir)

	// Should not error even if no snapshots dir
	if err := listSnapshots("dev"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = dir
}

func TestRestoreSnapshot(t *testing.T) {
	dir := setupTempDir(t)
	defer os.RemoveAll(dir)

	if err := addProfile("staging", dir); err != nil {
		t.Fatal(err)
	}
	p := profilePath("staging")
	os.WriteFile(p, []byte("A=1\n"), 0644)

	if err := snapshotProfile("staging"); err != nil {
		t.Fatal(err)
	}

	// Overwrite profile
	os.WriteFile(p, []byte("A=999\n"), 0644)

	// Get snapshot timestamp
	entries, _ := os.ReadDir(snapshotsDir("staging"))
	ts := strings.TrimSuffix(entries[0].Name(), ".env")

	if err := restoreSnapshot("staging", ts); err != nil {
		t.Fatalf("restore failed: %v", err)
	}

	data, _ := os.ReadFile(p)
	if string(data) != "A=1\n" {
		t.Fatalf("expected restored content, got %q", string(data))
	}
}

func TestRestoreNonExistentSnapshot(t *testing.T) {
	dir := setupTempDir(t)
	defer os.RemoveAll(dir)

	err := restoreSnapshot("dev", "00000000T000000")
	if err == nil {
		t.Fatal("expected error")
	}
	_ = filepath.Join(dir, ".envoy")
}
