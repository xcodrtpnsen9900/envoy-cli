package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupFmtDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("ENVOY_DIR", filepath.Join(dir, ".envoy"))
	if err := os.MkdirAll(filepath.Join(dir, ".envoy"), 0755); err != nil {
		t.Fatal(err)
	}
	return dir
}

func writeFmtProfile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, ".envoy", name+".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func readFmtProfile(t *testing.T, dir, name string) string {
	t.Helper()
	path := filepath.Join(dir, ".envoy", name+".env")
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestFmtProfileTrimsWhitespace(t *testing.T) {
	dir := setupFmtDir(t)
	writeFmtProfile(t, dir, "dev", "  KEY = value  \nFOO=bar\n")
	if err := fmtProfile("dev", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := readFmtProfile(t, dir, "dev")
	if !strings.Contains(out, "KEY=value") {
		t.Errorf("expected KEY=value, got: %s", out)
	}
	if !strings.Contains(out, "FOO=bar") {
		t.Errorf("expected FOO=bar, got: %s", out)
	}
}

func TestFmtProfilePreservesComments(t *testing.T) {
	dir := setupFmtDir(t)
	writeFmtProfile(t, dir, "dev", "# this is a comment\nKEY=value\n")
	if err := fmtProfile("dev", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := readFmtProfile(t, dir, "dev")
	if !strings.Contains(out, "# this is a comment") {
		t.Errorf("expected comment preserved, got: %s", out)
	}
}

func TestFmtProfileNonExistent(t *testing.T) {
	setupFmtDir(t)
	err := fmtProfile("ghost", false)
	if err == nil {
		t.Fatal("expected error for non-existent profile")
	}
}

func TestFmtProfileDryRun(t *testing.T) {
	dir := setupFmtDir(t)
	original := "  KEY = value  \n"
	writeFmtProfile(t, dir, "dev", original)
	if err := fmtProfile("dev", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// File should remain unchanged in dry-run mode
	out := readFmtProfile(t, dir, "dev")
	if out != original {
		t.Errorf("dry-run should not modify file; got: %q", out)
	}
}

func TestFormatEnvLine(t *testing.T) {
	cases := []struct{ input, want string }{
		{"KEY=value", "KEY=value"},
		{"  KEY = value  ", "KEY=value"},
		{"# comment", "# comment"},
		{"", ""},
		{"NOEQUALS", "NOEQUALS"},
	}
	for _, c := range cases {
		got := formatEnvLine(c.input)
		if got != c.want {
			t.Errorf("formatEnvLine(%q) = %q, want %q", c.input, got, c.want)
		}
	}
}
