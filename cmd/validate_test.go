package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateValidProfile(t *testing.T) {
	dir := setupTempDir(t)
	os.WriteFile(filepath.Join(dir, ".envoy", "profiles", "good.env"), []byte("FOO=bar\nBAZ=123\n# comment\n"), 0644)

	if err := validateProfile("good"); err != nil {
		t.Fatalf("expected valid profile, got error: %v", err)
	}
}

func TestValidateMissingEquals(t *testing.T) {
	dir := setupTempDir(t)
	os.WriteFile(filepath.Join(dir, ".envoy", "profiles", "bad.env"), []byte("FOOBAR\nBAZ=ok\n"), 0644)

	if err := validateProfile("bad"); err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestValidateEmptyKey(t *testing.T) {
	dir := setupTempDir(t)
	os.WriteFile(filepath.Join(dir, ".envoy", "profiles", "emptykey.env"), []byte("=value\n"), 0644)

	if err := validateProfile("emptykey"); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestValidateNonExistentProfile(t *testing.T) {
	setupTempDir(t)

	if err := validateProfile("ghost"); err == nil {
		t.Fatal("expected error for non-existent profile")
	}
}

func TestValidateKeyWithWhitespace(t *testing.T) {
	dir := setupTempDir(t)
	os.WriteFile(filepath.Join(dir, ".envoy", "profiles", "ws.env"), []byte("MY KEY=value\n"), 0644)

	if err := validateProfile("ws"); err == nil {
		t.Fatal("expected error for key with whitespace")
	}
}
