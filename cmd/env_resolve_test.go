package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func setupResolveDir(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	envoyDir := filepath.Join(tmp, ".envoy")
	if err := os.MkdirAll(envoyDir, 0755); err != nil {
		t.Fatal(err)
	}
	projectDir = tmp
	return tmp
}

func writeResolveProfile(t *testing.T, root, name, content string) {
	t.Helper()
	p := filepath.Join(root, ".envoy", name+".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestResolveProfileBasic(t *testing.T) {
	root := setupResolveDir(t)
	writeResolveProfile(t, root, "dev", "BASE=/usr/local\nBIN=${BASE}/bin\n")

	lines, err := resolveProfile("dev")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[1] != "BIN=/usr/local/bin" {
		t.Errorf("expected BIN=/usr/local/bin, got %s", lines[1])
	}
}

func TestResolveProfileSkipsComments(t *testing.T) {
	root := setupResolveDir(t)
	writeResolveProfile(t, root, "staging", "# comment\nHOST=localhost\nURL=http://${HOST}:8080\n")

	lines, err := resolveProfile("staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[1] != "URL=http://localhost:8080" {
		t.Errorf("unexpected URL value: %s", lines[1])
	}
}

func TestResolveProfileNonExistent(t *testing.T) {
	setupResolveDir(t)
	_, err := resolveProfile("ghost")
	if err == nil {
		t.Fatal("expected error for non-existent profile")
	}
}

func TestResolveProfileNoReferences(t *testing.T) {
	root := setupResolveDir(t)
	writeResolveProfile(t, root, "prod", "KEY=value\nOTHER=plain\n")

	lines, err := resolveProfile("prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lines[0] != "KEY=value" || lines[1] != "OTHER=plain" {
		t.Errorf("unexpected output: %v", lines)
	}
}

func TestExpandVars(t *testing.T) {
	env := map[string]string{"HOME": "/home/user", "NAME": "world"}
	cases := []struct {
		input    string
		expected string
	}{
		{"${HOME}/bin", "/home/user/bin"},
		{"hello $NAME", "hello world"},
		{"no-vars", "no-vars"},
		{"${MISSING}", ""},
	}
	for _, c := range cases {
		got := expandVars(c.input, env)
		if got != c.expected {
			t.Errorf("expandVars(%q) = %q, want %q", c.input, got, c.expected)
		}
	}
}
