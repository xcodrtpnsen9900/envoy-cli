package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEncryptAndDecryptProfile(t *testing.T) {
	dir := setupTempDir(t)
	defer os.RemoveAll(dir)

	profilesDir := filepath.Join(dir, "profiles")
	if err := os.MkdirAll(profilesDir, 0755); err != nil {
		t.Fatal(err)
	}

	profile := "staging"
	content := "DB_HOST=localhost\nDB_PASS=secret\n"
	if err := os.WriteFile(filepath.Join(profilesDir, profile+".env"), []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	password := "strongpassword"
	if err := encryptProfile(profile, password); err != nil {
		t.Fatalf("encryptProfile failed: %v", err)
	}

	encFile := filepath.Join(profilesDir, profile+".enc")
	if _, err := os.Stat(encFile); os.IsNotExist(err) {
		t.Fatal("encrypted file not created")
	}

	if err := os.Remove(filepath.Join(profilesDir, profile+".env")); err != nil {
		t.Fatal(err)
	}

	if err := decryptProfile(profile, password); err != nil {
		t.Fatalf("decryptProfile failed: %v", err)
	}

	got, err := os.ReadFile(filepath.Join(profilesDir, profile+".env"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != content {
		t.Errorf("expected %q, got %q", content, string(got))
	}
}

func TestEncryptNonExistentProfile(t *testing.T) {
	dir := setupTempDir(t)
	defer os.RemoveAll(dir)

	if err := os.MkdirAll(filepath.Join(dir, "profiles"), 0755); err != nil {
		t.Fatal(err)
	}

	err := encryptProfile("ghost", "pass")
	if err == nil {
		t.Fatal("expected error for non-existent profile")
	}
}

func TestDecryptWrongPassword(t *testing.T) {
	dir := setupTempDir(t)
	defer os.RemoveAll(dir)

	profilesDir := filepath.Join(dir, "profiles")
	if err := os.MkdirAll(profilesDir, 0755); err != nil {
		t.Fatal(err)
	}

	profile := "prod"
	if err := os.WriteFile(filepath.Join(profilesDir, profile+".env"), []byte("KEY=val\n"), 0600); err != nil {
		t.Fatal(err)
	}

	if err := encryptProfile(profile, "correctpass"); err != nil {
		t.Fatal(err)
	}

	err := decryptProfile(profile, "wrongpass")
	if err == nil {
		t.Fatal("expected error for wrong password")
	}
}

func TestDecryptNonExistentEncFile(t *testing.T) {
	dir := setupTempDir(t)
	defer os.RemoveAll(dir)

	if err := os.MkdirAll(filepath.Join(dir, "profiles"), 0755); err != nil {
		t.Fatal(err)
	}

	err := decryptProfile("missing", "pass")
	if err == nil {
		t.Fatal("expected error for missing encrypted file")
	}
}
