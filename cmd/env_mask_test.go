package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func setupMaskDir(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	envoyDir := filepath.Join(tmp, ".envoy")
	if err := os.MkdirAll(envoyDir, 0755); err != nil {
		t.Fatal(err)
	}
	return tmp
}

func writeMaskProfile(t *testing.T, dir, name, content string) {
	t.Helper()
	p := filepath.Join(dir, ".envoy", name+".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestMaskValueShort(t *testing.T) {
	if got := maskValue("ab"); got != "**" {
		t.Errorf("expected **, got %s", got)
	}
}

func TestMaskValueLong(t *testing.T) {
	got := maskValue("mysecret")
	if got[0:2] != "my" || got[len(got)-2:] != "et" {
		t.Errorf("unexpected mask result: %s", got)
	}
	if len(got) != len("mysecret") {
		t.Errorf("length mismatch: %d vs %d", len(got), len("mysecret"))
	}
}

func TestMaskValueEmpty(t *testing.T) {
	if got := maskValue(""); got != "" {
		t.Errorf("expected empty, got %s", got)
	}
}

func TestIsSensitiveKey(t *testing.T) {
	keys := defaultSensitiveKeys()
	if !isSensitiveKey("DB_PASSWORD", keys) {
		t.Error("expected DB_PASSWORD to be sensitive")
	}
	if !isSensitiveKey("API_TOKEN", keys) {
		t.Error("expected API_TOKEN to be sensitive")
	}
	if isSensitiveKey("APP_NAME", keys) {
		t.Error("expected APP_NAME to not be sensitive")
	}
}

func TestMaskProfileMasksSensitiveKeys(t *testing.T) {
	tmp := setupMaskDir(t)
	writeMaskProfile(t, tmp, "prod", "APP_NAME=myapp\nDB_PASSWORD=supersecret\n# comment\nAPI_KEY=abc123\n")

	old := projectDir
	projectDir = tmp
	defer func() { projectDir = old }()

	path := filepath.Join(tmp, ".envoy", "prod.env")
	err := maskProfile(path, defaultSensitiveKeys(), false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMaskProfileNonExistent(t *testing.T) {
	err := maskProfile("/nonexistent/.envoy/ghost.env", defaultSensitiveKeys(), false)
	if err == nil {
		t.Error("expected error for non-existent profile")
	}
}

func TestMaskProfileReveal(t *testing.T) {
	tmp := setupMaskDir(t)
	writeMaskProfile(t, tmp, "dev", "SECRET_KEY=topsecret\n")

	path := filepath.Join(tmp, ".envoy", "dev.env")
	err := maskProfile(path, defaultSensitiveKeys(), true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
