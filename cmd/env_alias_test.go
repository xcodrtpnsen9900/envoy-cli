package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func setupAliasDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".envoy"), 0755); err != nil {
		t.Fatal(err)
	}
	return dir
}

func writeAliasProfile(t *testing.T, dir, name string) {
	t.Helper()
	f := profilePath(dir, name)
	if err := os.WriteFile(f, []byte("KEY=value\n"), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestLoadAliasesEmpty(t *testing.T) {
	dir := setupAliasDir(t)
	aliases, err := loadAliases(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(aliases) != 0 {
		t.Errorf("expected empty aliases, got %v", aliases)
	}
}

func TestSaveAndLoadAliases(t *testing.T) {
	dir := setupAliasDir(t)
	aliases := map[string]string{"prod": "production", "dev": "development"}
	if err := saveAliases(dir, aliases); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	loaded, err := loadAliases(dir)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if loaded["prod"] != "production" || loaded["dev"] != "development" {
		t.Errorf("unexpected aliases: %v", loaded)
	}
}

func TestResolveAliasFound(t *testing.T) {
	dir := setupAliasDir(t)
	_ = saveAliases(dir, map[string]string{"p": "production"})
	resolved, err := resolveAlias(dir, "p")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved != "production" {
		t.Errorf("expected 'production', got %q", resolved)
	}
}

func TestResolveAliasNotFound(t *testing.T) {
	dir := setupAliasDir(t)
	resolved, err := resolveAlias(dir, "staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resolved != "staging" {
		t.Errorf("expected passthrough 'staging', got %q", resolved)
	}
}

func TestAliasesForProfile(t *testing.T) {
	dir := setupAliasDir(t)
	_ = saveAliases(dir, map[string]string{"p": "production", "prod": "production", "dev": "development"})
	aliases, err := aliasesForProfile(dir, "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(aliases) != 2 {
		t.Errorf("expected 2 aliases, got %d: %v", len(aliases), aliases)
	}
}

func TestAliasCount(t *testing.T) {
	dir := setupAliasDir(t)
	_ = saveAliases(dir, map[string]string{"a": "alpha", "b": "beta"})
	count, err := aliasCount(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2, got %d", count)
	}
}

func TestAliasExists(t *testing.T) {
	dir := setupAliasDir(t)
	_ = saveAliases(dir, map[string]string{"x": "xray"})
	ok, err := aliasExists(dir, "x")
	if err != nil || !ok {
		t.Errorf("expected alias 'x' to exist")
	}
	ok, err = aliasExists(dir, "z")
	if err != nil || ok {
		t.Errorf("expected alias 'z' to not exist")
	}
}
