package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupRenameKeyDir(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	t.Setenv("ENVOY_DIR", tmp)
	envoyDir := filepath.Join(tmp, ".envoy")
	if err := os.MkdirAll(envoyDir, 0755); err != nil {
		t.Fatal(err)
	}
	return tmp
}

func writeRenameKeyProfile(t *testing.T, profile, content string) {
	t.Helper()
	p := profilePath(profile)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestRenameKeyBasic(t *testing.T) {
	setupRenameKeyDir(t)
	writeRenameKeyProfile(t, "dev", "OLD_KEY=hello\nOTHER=world\n")

	if err := renameProfileKey("dev", "OLD_KEY", "NEW_KEY"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(profilePath("dev"))
	contents := string(data)
	if strings.Contains(contents, "OLD_KEY") {
		t.Error("old key still present after rename")
	}
	if !strings.Contains(contents, "NEW_KEY=hello") {
		t.Error("new key not found after rename")
	}
	if !strings.Contains(contents, "OTHER=world") {
		t.Error("unrelated key should be preserved")
	}
}

func TestRenameKeyNonExistentProfile(t *testing.T) {
	setupRenameKeyDir(t)
	err := renameProfileKey("ghost", "A", "B")
	if err == nil {
		t.Fatal("expected error for non-existent profile")
	}
}

func TestRenameKeyNotFound(t *testing.T) {
	setupRenameKeyDir(t)
	writeRenameKeyProfile(t, "dev", "EXISTING=value\n")

	err := renameProfileKey("dev", "MISSING", "NEW")
	if err == nil {
		t.Fatal("expected error when old key not found")
	}
}

func TestRenameKeyNewAlreadyExists(t *testing.T) {
	setupRenameKeyDir(t)
	writeRenameKeyProfile(t, "dev", "FOO=1\nBAR=2\n")

	err := renameProfileKey("dev", "FOO", "BAR")
	if err == nil {
		t.Fatal("expected error when new key already exists")
	}
}

func TestRenameKeyPreservesComments(t *testing.T) {
	setupRenameKeyDir(t)
	writeRenameKeyProfile(t, "dev", "# comment\nOLD=val\n")

	if err := renameProfileKey("dev", "OLD", "NEW"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, _ := os.ReadFile(profilePath("dev"))
	if !strings.Contains(string(data), "# comment") {
		t.Error("comment should be preserved after rename")
	}
}

func TestKeyExistsInProfile(t *testing.T) {
	setupRenameKeyDir(t)
	writeRenameKeyProfile(t, "dev", "ALPHA=1\n# comment\nBETA=2\n")

	ok, err := keyExistsInProfile("dev", "ALPHA")
	if err != nil || !ok {
		t.Error("expected ALPHA to exist")
	}
	ok, err = keyExistsInProfile("dev", "GAMMA")
	if err != nil || ok {
		t.Error("expected GAMMA to not exist")
	}
}

func TestListKeysInProfile(t *testing.T) {
	setupRenameKeyDir(t)
	writeRenameKeyProfile(t, "dev", "# comment\nA=1\nB=2\n")

	keys, err := listKeysInProfile("dev")
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != 2 || keys[0] != "A" || keys[1] != "B" {
		t.Errorf("unexpected keys: %v", keys)
	}
}
