package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSetNewKey(t *testing.T) {
	dir := setupTempDir(t)
	addProfile(dir, "dev", "FOO=bar\n")

	if err := setProfileKeys(dir, "dev", []string{"BAZ=qux"}); err != nil {
		t.Fatal(err)
	}
	content := readProfile(t, dir, "dev")
	if !strings.Contains(content, "BAZ=qux") {
		t.Errorf("expected BAZ=qux in profile, got: %s", content)
	}
}

func TestSetExistingKey(t *testing.T) {
	dir := setupTempDir(t)
	addProfile(dir, "dev", "FOO=bar\nBAZ=old\n")

	if err := setProfileKeys(dir, "dev", []string{"BAZ=new"}); err != nil {
		t.Fatal(err)
	}
	content := readProfile(t, dir, "dev")
	if !strings.Contains(content, "BAZ=new") {
		t.Errorf("expected BAZ=new, got: %s", content)
	}
	if strings.Contains(content, "BAZ=old") {
		t.Errorf("old value should be replaced")
	}
}

func TestSetMultipleKeys(t *testing.T) {
	dir := setupTempDir(t)
	addProfile(dir, "dev", "A=1\n")

	if err := setProfileKeys(dir, "dev", []string{"A=2", "B=3"}); err != nil {
		t.Fatal(err)
	}
	content := readProfile(t, dir, "dev")
	if !strings.Contains(content, "A=2") || !strings.Contains(content, "B=3") {
		t.Errorf("unexpected content: %s", content)
	}
}

func TestSetNonExistentProfile(t *testing.T) {
	dir := setupTempDir(t)
	err := setProfileKeys(dir, "ghost", []string{"X=1"})
	if err == nil {
		t.Fatal("expected error for non-existent profile")
	}
}

func TestSetInvalidPair(t *testing.T) {
	dir := setupTempDir(t)
	addProfile(dir, "dev", "")
	err := setProfileKeys(dir, "dev", []string{"NOEQUALS"})
	if err == nil {
		t.Fatal("expected error for invalid pair")
	}
}

func readProfile(t *testing.T, root, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(root, ".envoy", "profiles", name+".env"))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
