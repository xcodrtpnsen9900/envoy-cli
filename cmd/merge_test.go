package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMergeProfiles(t *testing.T) {
	dir := setupTempDir(t)

	os.WriteFile(filepath.Join(dir, ".envoy", "profiles", "base.env"), []byte("FOO=1\nBAR=2\n"), 0644)
	os.WriteFile(filepath.Join(dir, ".envoy", "profiles", "overlay.env"), []byte("BAR=99\nBAZ=3\n"), 0644)

	err := mergeProfiles("base", "overlay", "merged")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, ".envoy", "profiles", "merged.env"))
	if err != nil {
		t.Fatalf("merged profile not created: %v", err)
	}

	content := string(data)
	if !strings.Contains(content, "BAR=99") {
		t.Errorf("expected overlay value BAR=99, got: %s", content)
	}
	if !strings.Contains(content, "FOO=1") {
		t.Errorf("expected base value FOO=1, got: %s", content)
	}
	if !strings.Contains(content, "BAZ=3") {
		t.Errorf("expected overlay value BAZ=3, got: %s", content)
	}
}

func TestMergeNonExistentBase(t *testing.T) {
	setupTempDir(t)

	err := mergeProfiles("ghost", "overlay", "out")
	if err == nil {
		t.Fatal("expected error for missing base profile")
	}
}

func TestMergeOutputAlreadyExists(t *testing.T) {
	dir := setupTempDir(t)

	os.WriteFile(filepath.Join(dir, ".envoy", "profiles", "a.env"), []byte("X=1\n"), 0644)
	os.WriteFile(filepath.Join(dir, ".envoy", "profiles", "b.env"), []byte("Y=2\n"), 0644)
	os.WriteFile(filepath.Join(dir, ".envoy", "profiles", "out.env"), []byte("Z=3\n"), 0644)

	err := mergeProfiles("a", "b", "out")
	if err == nil {
		t.Fatal("expected error when output profile already exists")
	}
}

func TestReadEnvMapSkipsComments(t *testing.T) {
	dir := setupTempDir(t)
	path := filepath.Join(dir, ".envoy", "profiles", "test.env")
	os.WriteFile(path, []byte("# comment\nKEY=val\n\nOTHER=x\n"), 0644)

	m, err := readEnvMap(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(m) != 2 {
		t.Errorf("expected 2 entries, got %d", len(m))
	}
	if m["KEY"] != "val" {
		t.Errorf("expected KEY=val, got %s", m["KEY"])
	}
}
