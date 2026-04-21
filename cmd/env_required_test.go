package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func setupRequiredDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("ENVOY_DIR", dir)
	envoyDir := filepath.Join(dir, ".envoy")
	if err := os.MkdirAll(envoyDir, 0755); err != nil {
		t.Fatal(err)
	}
	return dir
}

func writeRequiredProfile(t *testing.T, dir, name, content string) {
	t.Helper()
	p := filepath.Join(dir, ".envoy", name+".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestAssertRequiredKeysAllPresent(t *testing.T) {
	dir := setupRequiredDir(t)
	writeRequiredProfile(t, dir, "dev", "FOO=bar\nBAZ=qux\n")
	missing, empty := assertRequiredKeys("dev", []string{"FOO", "BAZ"}, false)
	if len(missing) != 0 || len(empty) != 0 {
		t.Fatalf("expected no issues, got missing=%v empty=%v", missing, empty)
	}
}

func TestAssertRequiredKeysMissing(t *testing.T) {
	dir := setupRequiredDir(t)
	writeRequiredProfile(t, dir, "dev", "FOO=bar\n")
	missing, _ := assertRequiredKeys("dev", []string{"FOO", "MISSING_KEY"}, false)
	if len(missing) != 1 || missing[0] != "MISSING_KEY" {
		t.Fatalf("expected MISSING_KEY to be missing, got %v", missing)
	}
}

func TestAssertRequiredKeysStrictEmpty(t *testing.T) {
	dir := setupRequiredDir(t)
	writeRequiredProfile(t, dir, "dev", "FOO=\nBAR=hello\n")
	_, empty := assertRequiredKeys("dev", []string{"FOO", "BAR"}, true)
	if len(empty) != 1 || empty[0] != "FOO" {
		t.Fatalf("expected FOO to be empty, got %v", empty)
	}
}

func TestAssertRequiredKeysNonExistentProfile(t *testing.T) {
	setupRequiredDir(t)
	missing, _ := assertRequiredKeys("ghost", []string{"A", "B"}, false)
	if len(missing) != 2 {
		t.Fatalf("expected both keys missing for non-existent profile, got %v", missing)
	}
}

func TestRequiredKeySummaryOK(t *testing.T) {
	s := requiredKeySummary("prod", nil, nil)
	if s == "" {
		t.Fatal("expected non-empty summary")
	}
}

func TestPartitionKeys(t *testing.T) {
	bare, exact := partitionKeys([]string{"FOO", "BAR=hello", "BAZ"})
	if len(bare) != 2 {
		t.Fatalf("expected 2 bare keys, got %d", len(bare))
	}
	if exact["BAR"] != "hello" {
		t.Fatalf("expected exact[BAR]=hello, got %q", exact["BAR"])
	}
}
