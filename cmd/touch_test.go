package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTouchCreatesProfile(t *testing.T) {
	dir := setupTempDir(t)
	projectDir = dir

	envoyDir := filepath.Join(dir, ".envoy")
	if err := os.MkdirAll(envoyDir, 0755); err != nil {
		t.Fatal(err)
	}

	touchProfile(nil, []string{"newprofile"})

	path := profilePath("newprofile")
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected profile to exist after touch: %v", err)
	}
}

func TestTouchExistingProfilePreservesContent(t *testing.T) {
	dir := setupTempDir(t)
	projectDir = dir

	if err := addProfile([]string{"existing"}, nil); err != nil {
		t.Fatal(err)
	}

	path := profilePath("existing")
	content := []byte("KEY=preserved\n")
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}

	touchProfile(nil, []string{"existing"})

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != string(content) {
		t.Errorf("content was overwritten: got %q", data)
	}
}

func TestTouchCreatesEmptyFile(t *testing.T) {
	dir := setupTempDir(t)
	projectDir = dir

	envoyDir := filepath.Join(dir, ".envoy")
	if err := os.MkdirAll(envoyDir, 0755); err != nil {
		t.Fatal(err)
	}

	touchProfile(nil, []string{"empty"})

	path := profilePath("empty")
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Size() != 0 {
		t.Errorf("expected empty file, got size %d", info.Size())
	}
}
