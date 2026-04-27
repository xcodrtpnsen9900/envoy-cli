package cmd

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func setupRollbackDir(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	envoyDir := filepath.Join(root, ".envoy")
	if err := os.MkdirAll(envoyDir, 0755); err != nil {
		t.Fatal(err)
	}
	return root
}

func writeRollbackProfile(t *testing.T, root, profile, content string) {
	t.Helper()
	p := profilePath(root, profile)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func makeSnapshot(t *testing.T, root, profile, content string) {
	t.Helper()
	dir := snapshotsDir(root)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	ts := time.Now().UnixNano()
	name := filepath.Join(dir, profile+"_"+time.Unix(0, ts).Format("20060102150405")+".env")
	if err := os.WriteFile(name, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Millisecond)
}

func TestRollbackNoSnapshots(t *testing.T) {
	root := setupRollbackDir(t)
	writeRollbackProfile(t, root, "dev", "KEY=original\n")
	err := rollbackProfile(root, "dev", 1)
	if err == nil {
		t.Fatal("expected error for missing snapshots")
	}
}

func TestRollbackOneStep(t *testing.T) {
	root := setupRollbackDir(t)
	writeRollbackProfile(t, root, "dev", "KEY=current\n")
	makeSnapshot(t, root, "dev", "KEY=previous\n")
	if err := rollbackProfile(root, "dev", 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(profilePath(root, "dev"))
	if string(data) != "KEY=previous\n" {
		t.Errorf("expected rolled-back content, got: %s", data)
	}
}

func TestRollbackTwoSteps(t *testing.T) {
	root := setupRollbackDir(t)
	writeRollbackProfile(t, root, "dev", "KEY=current\n")
	makeSnapshot(t, root, "dev", "KEY=oldest\n")
	makeSnapshot(t, root, "dev", "KEY=newer\n")
	if err := rollbackProfile(root, "dev", 2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(profilePath(root, "dev"))
	if string(data) != "KEY=oldest\n" {
		t.Errorf("expected oldest content, got: %s", data)
	}
}

func TestRollbackCountHelper(t *testing.T) {
	root := setupRollbackDir(t)
	writeRollbackProfile(t, root, "staging", "A=1\n")
	makeSnapshot(t, root, "staging", "A=1\n")
	makeSnapshot(t, root, "staging", "A=2\n")
	if c := rollbackCount(root, "staging"); c != 2 {
		t.Errorf("expected 2 snapshots, got %d", c)
	}
}

func TestLatestRollbackSnapshotHelper(t *testing.T) {
	root := setupRollbackDir(t)
	writeRollbackProfile(t, root, "prod", "X=1\n")
	makeSnapshot(t, root, "prod", "X=1\n")
	makeSnapshot(t, root, "prod", "X=2\n")
	latest := latestRollbackSnapshot(root, "prod")
	if latest == "" {
		t.Fatal("expected a snapshot path")
	}
	data, _ := os.ReadFile(latest)
	if string(data) != "X=2\n" {
		t.Errorf("expected latest content X=2, got: %s", data)
	}
}

func TestProfilesWithRollbackPoints(t *testing.T) {
	root := setupRollbackDir(t)
	writeRollbackProfile(t, root, "dev", "K=1\n")
	writeRollbackProfile(t, root, "prod", "K=2\n")
	makeSnapshot(t, root, "dev", "K=1\n")
	makeSnapshot(t, root, "prod", "K=2\n")
	profiles, err := profilesWithRollbackPoints(root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(profiles) != 2 {
		t.Errorf("expected 2 profiles with rollback points, got %d", len(profiles))
	}
}
