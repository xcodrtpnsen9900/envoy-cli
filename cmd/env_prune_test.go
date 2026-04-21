package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupPruneDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	envoyDir := filepath.Join(dir, ".envoy")
	if err := os.MkdirAll(envoyDir, 0755); err != nil {
		t.Fatalf("failed to create .envoy dir: %v", err)
	}
	projectDir = dir
	return dir
}

func writePruneProfile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, ".envoy", name+".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write profile: %v", err)
	}
}

func readPruneProfile(t *testing.T, dir, name string) string {
	t.Helper()
	path := filepath.Join(dir, ".envoy", name+".env")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read profile: %v", err)
	}
	return string(data)
}

func TestPruneRemovesDuplicateKeys(t *testing.T) {
	dir := setupPruneDir(t)
	writePruneProfile(t, dir, "dev", "FOO=1\nBAR=2\nFOO=3\n")

	pruneProfile("dev", false)

	content := readPruneProfile(t, dir, "dev")
	lines := strings.Split(strings.TrimSpace(content), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines after prune, got %d: %v", len(lines), lines)
	}
	if !strings.Contains(content, "FOO=1") {
		t.Error("expected first occurrence of FOO=1 to be kept")
	}
	if strings.Contains(content, "FOO=3") {
		t.Error("expected duplicate FOO=3 to be removed")
	}
}

func TestPruneRemovesEmptyKeyLines(t *testing.T) {
	dir := setupPruneDir(t)
	writePruneProfile(t, dir, "dev", "FOO=bar\n=empty\nBAZ=qux\n")

	pruneProfile("dev", false)

	content := readPruneProfile(t, dir, "dev")
	if strings.Contains(content, "=empty") {
		t.Error("expected empty key line to be removed")
	}
	if !strings.Contains(content, "FOO=bar") || !strings.Contains(content, "BAZ=qux") {
		t.Error("expected valid keys to be preserved")
	}
}

func TestPrunePreservesCommentsAndBlanks(t *testing.T) {
	dir := setupPruneDir(t)
	writePruneProfile(t, dir, "dev", "# comment\nFOO=1\n\nBAR=2\n")

	pruneProfile("dev", false)

	content := readPruneProfile(t, dir, "dev")
	if !strings.Contains(content, "# comment") {
		t.Error("expected comment to be preserved")
	}
	if !strings.Contains(content, "FOO=1") || !strings.Contains(content, "BAR=2") {
		t.Error("expected keys to be preserved")
	}
}

func TestPruneNonExistentProfile(t *testing.T) {
	setupPruneDir(t)
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected fatalf to be called for non-existent profile")
		}
	}()
	pruneProfile("ghost", false)
}

func TestPruneDryRunDoesNotModify(t *testing.T) {
	dir := setupPruneDir(t)
	original := "FOO=1\nFOO=2\nBAR=3\n"
	writePruneProfile(t, dir, "dev", original)

	pruneProfile("dev", true)

	content := readPruneProfile(t, dir, "dev")
	if content != original {
		t.Errorf("dry-run should not modify file; got %q", content)
	}
}
