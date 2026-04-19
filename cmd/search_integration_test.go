package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func runSearchCommand(t *testing.T, dir string, args ...string) (string, error) {
	t.Helper()
	prevDir := projectDir
	projectDir = dir
	defer func() { projectDir = prevDir }()

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	allArgs := append([]string{"search"}, args...)
	rootCmd.SetArgs(allArgs)

	var execErr error
	for _, c := range rootCmd.Commands() {
		if c.Use == "search [profile] [key]" {
			c.SetOut(buf)
			_ = c
		}
	}
	_ = execErr
	return buf.String(), nil
}

func TestSearchOutputContainsLineNumber(t *testing.T) {
	dir := t.TempDir()
	envoyDir := filepath.Join(dir, ".envoy")
	os.MkdirAll(envoyDir, 0755)
	os.WriteFile(filepath.Join(envoyDir, "prod.env"), []byte("HOST=example.com\nPORT=443\n"), 0644)

	results, err := searchProfile(dir, "prod", "HOST", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected at least one result")
	}
	if !strings.HasPrefix(results[0], "line ") {
		t.Errorf("expected result to start with 'line ', got: %s", results[0])
	}
}

func TestSearchEmptyProfile(t *testing.T) {
	dir := t.TempDir()
	envoyDir := filepath.Join(dir, ".envoy")
	os.MkdirAll(envoyDir, 0755)
	os.WriteFile(filepath.Join(envoyDir, "empty.env"), []byte(""), 0644)

	results, err := searchProfile(dir, "empty", "KEY", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results for empty profile, got %d", len(results))
	}
}

var _ = cobra.Command{}
