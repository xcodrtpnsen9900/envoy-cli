package cmd

import (
	"testing"
)

func TestProfilesByTag(t *testing.T) {
	dir := setupTempDir(t)
	if err := initProject(dir); err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{"dev", "staging", "prod"} {
		if err := addProfile(dir, p, nil); err != nil {
			t.Fatal(err)
		}
	}
	_ = tagProfile(dir, "dev", []string{"local"})
	_ = tagProfile(dir, "staging", []string{"cloud", "local"})
	_ = tagProfile(dir, "prod", []string{"cloud"})

	results, err := profilesByTag(dir, "local")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 profiles with tag 'local', got %d: %v", len(results), results)
	}
}

func TestProfilesByTagCaseInsensitive(t *testing.T) {
	dir := setupTempDir(t)
	if err := initProject(dir); err != nil {
		t.Fatal(err)
	}
	if err := addProfile(dir, "dev", nil); err != nil {
		t.Fatal(err)
	}
	_ = tagProfile(dir, "dev", []string{"Cloud"})
	results, err := profilesByTag(dir, "cloud")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0] != "dev" {
		t.Errorf("expected [dev], got %v", results)
	}
}

func TestProfilesByTagNoMatch(t *testing.T) {
	dir := setupTempDir(t)
	if err := initProject(dir); err != nil {
		t.Fatal(err)
	}
	if err := addProfile(dir, "dev", nil); err != nil {
		t.Fatal(err)
	}
	_ = tagProfile(dir, "dev", []string{"local"})
	results, err := profilesByTag(dir, "remote")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Errorf("expected no results, got %v", results)
	}
}

func TestProfilesByTagNotInitialized(t *testing.T) {
	dir := setupTempDir(t)
	_, err := profilesByTag(dir, "any")
	if err == nil {
		t.Fatal("expected error when envoy not initialized")
	}
}

func TestProfilesByTagEmptyTag(t *testing.T) {
	dir := setupTempDir(t)
	if err := initProject(dir); err != nil {
		t.Fatal(err)
	}
	if err := addProfile(dir, "dev", nil); err != nil {
		t.Fatal(err)
	}
	results, err := profilesByTag(dir, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 0 {
		t.Errorf("expected no results for empty tag, got %v", results)
	}
}
