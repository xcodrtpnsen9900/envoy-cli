package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func writeSortTestProfile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name+".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestSortProfileBasic(t *testing.T) {
	dir := t.TempDir()
	path := writeSortTestProfile(t, dir, "dev", "ZEBRA=1\nAPPLE=2\nMIDDLE=3\n")

	kvLines, _, err := sortProfile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(kvLines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(kvLines))
	}
	if extractKey(kvLines[0]) != "APPLE" {
		t.Errorf("expected APPLE first, got %s", kvLines[0])
	}
	if extractKey(kvLines[2]) != "ZEBRA" {
		t.Errorf("expected ZEBRA last, got %s", kvLines[2])
	}
}

func TestSortProfilePreservesComments(t *testing.T) {
	dir := t.TempDir()
	path := writeSortTestProfile(t, dir, "dev", "# header\nZOO=1\nANT=2\n")

	_, comments, err := sortProfile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(comments) != 1 || comments[0] != "# header" {
		t.Errorf("expected comment preserved, got %v", comments)
	}
}

func TestSortProfileNonExistent(t *testing.T) {
	_, _, err := sortProfile("/nonexistent/path/profile.env")
	if err == nil {
		t.Error("expected error for non-existent profile")
	}
}

func TestWriteSortedProfile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.env")

	comments := []string{"# comment"}
	kvLines := []string{"A=1", "B=2"}
	if err := writeSortedProfile(path, comments, kvLines); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(path)
	content := string(data)
	if content != "# comment\nA=1\nB=2\n" {
		t.Errorf("unexpected file content: %q", content)
	}
}
