package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTagProfile(t *testing.T) {
	dir := setupTempDir(t)
	if err := initProject(dir); err != nil {
		t.Fatal(err)
	}
	if err := addProfile(dir, "staging", nil); err != nil {
		t.Fatal(err)
	}
	if err := tagProfile(dir, "staging", []string{"cloud", "prod-like"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tags, err := getTags(dir, "staging")
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 2 || tags[0] != "cloud" || tags[1] != "prod-like" {
		t.Errorf("expected [cloud prod-like], got %v", tags)
	}
}

func TestTagProfileNonExistent(t *testing.T) {
	dir := setupTempDir(t)
	if err := initProject(dir); err != nil {
		t.Fatal(err)
	}
	err := tagProfile(dir, "ghost", []string{"x"})
	if err == nil {
		t.Fatal("expected error for non-existent profile")
	}
}

func TestTagProfileDeduplicates(t *testing.T) {
	dir := setupTempDir(t)
	if err := initProject(dir); err != nil {
		t.Fatal(err)
	}
	if err := addProfile(dir, "dev", nil); err != nil {
		t.Fatal(err)
	}
	_ = tagProfile(dir, "dev", []string{"local", "debug"})
	_ = tagProfile(dir, "dev", []string{"debug", "verbose"})
	tags, _ := getTags(dir, "dev")
	if len(tags) != 3 {
		t.Errorf("expected 3 unique tags, got %d: %v", len(tags), tags)
	}
}

func TestGetTagsNoFile(t *testing.T) {
	dir := setupTempDir(t)
	if err := initProject(dir); err != nil {
		t.Fatal(err)
	}
	tags, err := getTags(dir, "missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tags != nil {
		t.Errorf("expected nil tags, got %v", tags)
	}
}

func TestTagFileStoredCorrectly(t *testing.T) {
	dir := setupTempDir(t)
	if err := initProject(dir); err != nil {
		t.Fatal(err)
	}
	if err := addProfile(dir, "ci", nil); err != nil {
		t.Fatal(err)
	}
	_ = tagProfile(dir, "ci", []string{"automated"})
	path := filepath.Join(dir, ".envoy", "ci.tags")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("expected tag file to exist")
	}
}
