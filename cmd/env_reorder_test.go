package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupReorderDir(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".envoy"), 0755); err != nil {
		t.Fatal(err)
	}
	return root
}

func writeReorderProfile(t *testing.T, root, profile, content string) {
	t.Helper()
	p := filepath.Join(root, ".envoy", profile+".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func readReorderProfile(t *testing.T, root, profile string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, ".envoy", profile+".env"))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestReorderProfileBasic(t *testing.T) {
	root := setupReorderDir(t)
	writeReorderProfile(t, root, "dev", "APP_NAME=myapp\nDEBUG=true\nPORT=8080\n")

	if err := reorderProfile(root, "dev", []string{"PORT", "DEBUG"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content := readReorderProfile(t, root, "dev")
	lines := strings.Split(strings.TrimRight(content, "\n"), "\n")
	if lines[0] != "PORT=8080" {
		t.Errorf("expected PORT first, got %q", lines[0])
	}
	if lines[1] != "DEBUG=true" {
		t.Errorf("expected DEBUG second, got %q", lines[1])
	}
}

func TestReorderProfileNonExistent(t *testing.T) {
	root := setupReorderDir(t)
	err := reorderProfile(root, "ghost", []string{"KEY"})
	if err == nil {
		t.Fatal("expected error for non-existent profile")
	}
}

func TestReorderProfilePreservesComments(t *testing.T) {
	root := setupReorderDir(t)
	writeReorderProfile(t, root, "dev", "# comment\nAPP=foo\nSECRET=bar\n")

	if err := reorderProfile(root, "dev", []string{"SECRET"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content := readReorderProfile(t, root, "dev")
	if !strings.Contains(content, "# comment") {
		t.Error("expected comment to be preserved")
	}
	if !strings.HasPrefix(content, "SECRET=bar") {
		t.Errorf("expected SECRET first, got: %q", content)
	}
}

func TestReorderProfileUnknownKeyIgnored(t *testing.T) {
	root := setupReorderDir(t)
	writeReorderProfile(t, root, "dev", "APP=foo\nPORT=3000\n")

	if err := reorderProfile(root, "dev", []string{"MISSING", "PORT"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content := readReorderProfile(t, root, "dev")
	if !strings.Contains(content, "PORT=3000") {
		t.Error("expected PORT to be present")
	}
}
