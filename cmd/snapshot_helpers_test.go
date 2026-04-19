package cmd

import (
	"os"
	"testing"
)

func TestLatestSnapshotNoSnapshots(t *testing.T) {
	dir := setupTempDir(t)
	defer os.RemoveAll(dir)

	ts, err := latestSnapshot("dev")
	if err != nil {
		t.Fatal(err)
	}
	if ts != "" {
		t.Fatalf("expected empty, got %q", ts)
	}
}

func TestLatestSnapshotReturnsNewest(t *testing.T) {
	dir := setupTempDir(t)
	defer os.RemoveAll(dir)

	if err := addProfile("dev", dir); err != nil {
		t.Fatal(err)
	}
	p := profilePath("dev")
	os.WriteFile(p, []byte("K=1\n"), 0644)

	snapshotProfile("dev")
	snapshotProfile("dev")

	ts, err := latestSnapshot("dev")
	if err != nil {
		t.Fatal(err)
	}
	if ts == "" {
		t.Fatal("expected a timestamp")
	}
}

func TestSnapshotCount(t *testing.T) {
	dir := setupTempDir(t)
	defer os.RemoveAll(dir)

	if err := addProfile("qa", dir); err != nil {
		t.Fatal(err)
	}
	p := profilePath("qa")
	os.WriteFile(p, []byte("X=1\n"), 0644)

	snapshotProfile("qa")
	snapshotProfile("qa")

	count, err := snapshotCount("qa")
	if err != nil {
		t.Fatal(err)
	}
	if count < 1 {
		t.Fatalf("expected at least 1 snapshot, got %d", count)
	}
}

func TestDeleteSnapshot(t *testing.T) {
	dir := setupTempDir(t)
	defer os.RemoveAll(dir)

	if err := addProfile("dev", dir); err != nil {
		t.Fatal(err)
	}
	p := profilePath("dev")
	os.WriteFile(p, []byte("K=v\n"), 0644)
	snapshotProfile("dev")

	ts, _ := latestSnapshot("dev")
	if err := deleteSnapshot("dev", ts); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	count, _ := snapshotCount("dev")
	if count != 0 {
		t.Fatalf("expected 0 snapshots after delete, got %d", count)
	}
}
