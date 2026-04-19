package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPinProfile(t *testing.T) {
	dir := setupTempDir(t)
	initProject(dir)
	addProfile(dir, "staging", false)

	if err := pinProfile(dir, "staging"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !isPinned(dir, "staging") {
		t.Error("expected staging to be pinned")
	}
}

func TestPinNonExistentProfile(t *testing.T) {
	dir := setupTempDir(t)
	initProject(dir)

	if err := pinProfile(dir, "ghost"); err == nil {
		t.Error("expected error for non-existent profile")
	}
}

func TestPinAlreadyPinned(t *testing.T) {
	dir := setupTempDir(t)
	initProject(dir)
	addProfile(dir, "prod", false)

	_ = pinProfile(dir, "prod")
	if err := pinProfile(dir, "prod"); err != nil {
		t.Fatalf("expected no error on double pin, got %v", err)
	}
	pinned, _ := readPinned(dir)
	count := 0
	for _, p := range pinned {
		if p == "prod" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected prod to appear once, got %d", count)
	}
}

func TestUnpinProfile(t *testing.T) {
	dir := setupTempDir(t)
	initProject(dir)
	addProfile(dir, "dev", false)
	_ = pinProfile(dir, "dev")

	if err := unpinProfile(dir, "dev"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if isPinned(dir, "dev") {
		t.Error("expected dev to be unpinned")
	}
}

func TestUnpinNotPinned(t *testing.T) {
	dir := setupTempDir(t)
	initProject(dir)
	addProfile(dir, "dev", false)

	if err := unpinProfile(dir, "dev"); err == nil {
		t.Error("expected error when unpinning a non-pinned profile")
	}
}

func TestPinFileStoredCorrectly(t *testing.T) {
	dir := setupTempDir(t)
	initProject(dir)
	addProfile(dir, "alpha", false)
	addProfile(dir, "beta", false)
	_ = pinProfile(dir, "alpha")
	_ = pinProfile(dir, "beta")

	data, err := os.ReadFile(filepath.Join(dir, ".envoy", "pinned"))
	if err != nil {
		t.Fatalf("expected pinned file to exist: %v", err)
	}
	content := string(data)
	if content == "" {
		t.Error("expected pinned file to have content")
	}
}
