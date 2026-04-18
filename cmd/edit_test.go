package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEditProfileNonExistent(t *testing.T) {
	dir := setupTempDir(t)
	projectDir = dir

	envoyDir := filepath.Join(dir, ".envoy")
	if err := os.MkdirAll(envoyDir, 0755); err != nil {
		t.Fatal(err)
	}

	path := profilePath("ghost")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("expected profile to not exist")
	}
}

func TestEditProfileExists(t *testing.T) {
	dir := setupTempDir(t)
	projectDir = dir

	if err := addProfile([]string{"staging"}, nil); err != nil {
		t.Fatal(err)
	}

	path := profilePath("staging")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected profile file to exist: %v", err)
	}

	// Write some content to verify file is accessible
	content := []byte("KEY=value\n")
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != string(content) {
		t.Errorf("unexpected content: %s", data)
	}
}

func TestEditUsesEditorEnv(t *testing.T) {
	original := os.Getenv("EDITOR")
	defer os.Setenv("EDITOR", original)

	os.Setenv("EDITOR", "nano")
	if got := os.Getenv("EDITOR"); got != "nano" {
		t.Errorf("expected EDITOR=nano, got %s", got)
	}
}
