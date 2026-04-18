package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExportProfileBash(t *testing.T) {
	dir := setupTempDir(t)
	setProjectDir(dir)

	profDir := filepath.Join(dir, ".envoy", "profiles")
	os.MkdirAll(profDir, 0755)
	os.WriteFile(filepath.Join(profDir, "dev.env"), []byte("FOO=bar\nBAZ=qux\n"), 0644)

	out, err := exportProfile("dev", "bash")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "export FOO=bar") {
		t.Errorf("expected export FOO=bar, got: %s", out)
	}
	if !strings.Contains(out, "export BAZ=qux") {
		t.Errorf("expected export BAZ=qux, got: %s", out)
	}
}

func TestExportProfileFish(t *testing.T) {
	dir := setupTempDir(t)
	setProjectDir(dir)

	profDir := filepath.Join(dir, ".envoy", "profiles")
	os.MkdirAll(profDir, 0755)
	os.WriteFile(filepath.Join(profDir, "prod.env"), []byte("KEY=value\n"), 0644)

	out, err := exportProfile("prod", "fish")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out, "set -x KEY value;") {
		t.Errorf("expected fish syntax, got: %s", out)
	}
}

func TestExportProfileSkipsComments(t *testing.T) {
	dir := setupTempDir(t)
	setProjectDir(dir)

	profDir := filepath.Join(dir, ".envoy", "profiles")
	os.MkdirAll(profDir, 0755)
	os.WriteFile(filepath.Join(profDir, "test.env"), []byte("# comment\nFOO=1\n"), 0644)

	out, err := exportProfile("test", "bash")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if strings.Contains(out, "comment") {
		t.Errorf("comments should be skipped, got: %s", out)
	}
	if !strings.Contains(out, "export FOO=1") {
		t.Errorf("expected export FOO=1, got: %s", out)
	}
}

func TestExportNonExistentProfile(t *testing.T) {
	dir := setupTempDir(t)
	setProjectDir(dir)
	os.MkdirAll(filepath.Join(dir, ".envoy", "profiles"), 0755)

	_, err := exportProfile("ghost", "bash")
	if err == nil {
		t.Fatal("expected error for non-existent profile")
	}
}
