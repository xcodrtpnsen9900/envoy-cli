package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func setupPromoteDir(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	dir := filepath.Join(root, ".envoy", "profiles")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	return root
}

func writePromoteProfile(t *testing.T, root, name, content string) {
	t.Helper()
	p := filepath.Join(root, ".envoy", "profiles", name+".env")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

// readPromoteProfile is a test helper that reads the env map for a named profile
// under the given root directory, failing the test on any error.
func readPromoteProfile(t *testing.T, root, name string) map[string]string {
	t.Helper()
	m, err := readEnvMap(filepath.Join(root, ".envoy", "profiles", name+".env"))
	if err != nil {
		t.Fatalf("failed to read profile %q: %v", name, err)
	}
	return m
}

func TestPromoteAllKeys(t *testing.T) {
	root := setupPromoteDir(t)
	writePromoteProfile(t, root, "staging", "DB_HOST=staging-db\nAPI_KEY=abc123\n")
	writePromoteProfile(t, root, "prod", "DB_HOST=prod-db\n")

	err := promoteProfile(root, "staging", "prod", nil, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m := readPromoteProfile(t, root, "prod")
	if m["API_KEY"] != "abc123" {
		t.Errorf("expected API_KEY=abc123, got %q", m["API_KEY"])
	}
	if m["DB_HOST"] != "prod-db" {
		t.Errorf("expected DB_HOST unchanged, got %q", m["DB_HOST"])
	}
}

func TestPromoteWithOverwrite(t *testing.T) {
	root := setupPromoteDir(t)
	writePromoteProfile(t, root, "staging", "DB_HOST=staging-db\n")
	writePromoteProfile(t, root, "prod", "DB_HOST=prod-db\n")

	err := promoteProfile(root, "staging", "prod", nil, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m := readPromoteProfile(t, root, "prod")
	if m["DB_HOST"] != "staging-db" {
		t.Errorf("expected DB_HOST=staging-db, got %q", m["DB_HOST"])
	}
}

func TestPromoteSpecificKeys(t *testing.T) {
	root := setupPromoteDir(t)
	writePromoteProfile(t, root, "staging", "DB_HOST=staging-db\nAPI_KEY=abc\nSECRET=xyz\n")
	writePromoteProfile(t, root, "prod", "DB_HOST=prod-db\n")

	err := promoteProfile(root, "staging", "prod", []string{"API_KEY"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m := readPromoteProfile(t, root, "prod")
	if m["API_KEY"] != "abc" {
		t.Errorf("expected API_KEY=abc, got %q", m["API_KEY"])
	}
	if _, ok := m["SECRET"]; ok {
		t.Error("SECRET should not have been promoted")
	}
}

func TestPromoteNonExistentSource(t *testing.T) {
	root := setupPromoteDir(t)
	writePromoteProfile(t, root, "prod", "DB_HOST=prod-db\n")

	err := promoteProfile(root, "ghost", "prod", nil, false)
	if err == nil {
		t.Error("expected error for missing source profile")
	}
}

func TestPromoteNonExistentTarget(t *testing.T) {
	root := setupPromoteDir(t)
	writePromoteProfile(t, root, "staging", "DB_HOST=staging-db\n")

	err := promoteProfile(root, "staging", "ghost", nil, false)
	if err == nil {
		t.Error("expected error for missing target profile")
	}
}
