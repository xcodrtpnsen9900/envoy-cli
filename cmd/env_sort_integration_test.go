package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func newTestRootForSort(t *testing.T) (root *cobra.Command, envoyDir string) {
	t.Helper()
	tmp := t.TempDir()
	envoyDir = filepath.Join(tmp, ".envoy")
	if err := os.MkdirAll(envoyDir, 0755); err != nil {
		t.Fatal(err)
	}
	origDir := projectDir
	projectDir = tmp
	t.Cleanup(func() { projectDir = origDir })
	return rootCmd, envoyDir
}

func runSortCommand(t *testing.T, args ...string) (string, error) {
	t.Helper()
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	return buf.String(), err
}

func TestSortOutputIsAlphabetical(t *testing.T) {
	_, envoyDir := newTestRootForSort(t)
	profile := filepath.Join(envoyDir, "staging.env")
	content := "ZEBRA=z\nAPPLE=a\nMIDDLE=m\n"
	if err := os.WriteFile(profile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	out, err := runSortCommand(t, "sort", "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v\noutput: %s", err, out)
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) < 3 {
		t.Fatalf("expected at least 3 lines, got: %v", lines)
	}
	if extractKey(lines[0]) != "APPLE" {
		t.Errorf("first line should be APPLE, got: %s", lines[0])
	}
	if extractKey(lines[2]) != "ZEBRA" {
		t.Errorf("last line should be ZEBRA, got: %s", lines[2])
	}
}

func TestSortNonExistentProfile(t *testing.T) {
	newTestRootForSort(t)
	_, err := runSortCommand(t, "sort", "ghost")
	if err == nil {
		t.Error("expected error for non-existent profile")
	}
}
