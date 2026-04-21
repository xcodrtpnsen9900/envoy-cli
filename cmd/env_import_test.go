package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setupImportDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	envoyDir := filepath.Join(dir, ".envoy")
	if err := os.MkdirAll(envoyDir, 0755); err != nil {
		t.Fatal(err)
	}
	return dir
}

func writeImportFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestImportProfileAddsNewKeys(t *testing.T) {
	dir := setupImportDir(t)
	t.Setenv("ENVOY_DIR", dir)

	profileFile := filepath.Join(dir, ".envoy", "dev.env")
	writeImportFile(t, profileFile, "EXISTING=value1\n")

	srcFile := filepath.Join(dir, "external.env")
	writeImportFile(t, srcFile, "NEW_KEY=hello\nANOTHER=world\n")

	importProfile("dev", srcFile, false)

	m, err := readImportEnvMap(profileFile)
	if err != nil {
		t.Fatal(err)
	}
	if m["NEW_KEY"] != "hello" {
		t.Errorf("expected NEW_KEY=hello, got %q", m["NEW_KEY"])
	}
	if m["ANOTHER"] != "world" {
		t.Errorf("expected ANOTHER=world, got %q", m["ANOTHER"])
	}
}

func TestImportProfileSkipsExistingWithoutOverwrite(t *testing.T) {
	dir := setupImportDir(t)
	t.Setenv("ENVOY_DIR", dir)

	profileFile := filepath.Join(dir, ".envoy", "dev.env")
	writeImportFile(t, profileFile, "KEY=original\n")

	srcFile := filepath.Join(dir, "external.env")
	writeImportFile(t, srcFile, "KEY=override\n")

	importProfile("dev", srcFile, false)

	m, _ := readImportEnvMap(profileFile)
	if m["KEY"] != "original" {
		t.Errorf("expected KEY=original (no overwrite), got %q", m["KEY"])
	}
}

func TestImportProfileOverwritesExistingKeys(t *testing.T) {
	dir := setupImportDir(t)
	t.Setenv("ENVOY_DIR", dir)

	profileFile := filepath.Join(dir, ".envoy", "dev.env")
	writeImportFile(t, profileFile, "KEY=original\n")

	srcFile := filepath.Join(dir, "external.env")
	writeImportFile(t, srcFile, "KEY=override\n")

	importProfile("dev", srcFile, true)

	m, _ := readImportEnvMap(profileFile)
	if m["KEY"] != "override" {
		t.Errorf("expected KEY=override, got %q", m["KEY"])
	}
}

func TestImportProfileSkipsComments(t *testing.T) {
	dir := setupImportDir(t)
	t.Setenv("ENVOY_DIR", dir)

	profileFile := filepath.Join(dir, ".envoy", "staging.env")
	writeImportFile(t, profileFile, "")

	srcFile := filepath.Join(dir, "external.env")
	writeImportFile(t, srcFile, "# this is a comment\nVALID=yes\n")

	importProfile("staging", srcFile, false)

	m, _ := readImportEnvMap(profileFile)
	if _, ok := m["# this is a comment"]; ok {
		t.Error("comment should not be imported as a key")
	}
	if m["VALID"] != "yes" {
		t.Errorf("expected VALID=yes, got %q", m["VALID"])
	}
}

func TestReadImportEnvMapMissingFile(t *testing.T) {
	_, err := readImportEnvMap("/nonexistent/path.env")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestWriteImportEnvMap(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.env")
	m := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if err := writeImportEnvMap(path, m); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(path)
	content := string(data)
	if !strings.Contains(content, "FOO=bar") {
		t.Errorf("expected FOO=bar in output, got: %s", content)
	}
}
