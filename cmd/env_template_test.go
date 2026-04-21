package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTemplateProfile(t *testing.T, dir, name, content string) {
	t.Helper()
	envoyDir := filepath.Join(dir, ".envoy")
	_ = os.MkdirAll(envoyDir, 0755)
	err := os.WriteFile(filepath.Join(envoyDir, name+".env"), []byte(content), 0644)
	if err != nil {
		t.Fatalf("writeTemplateProfile: %v", err)
	}
}

func TestRenderTemplateBasic(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVOY_DIR", dir)

	writeTemplateProfile(t, dir, "prod", "HOST=example.com\nPORT=8080\n")

	tmplFile := filepath.Join(dir, "app.conf.tmpl")
	os.WriteFile(tmplFile, []byte("server={{HOST}}:{{PORT}}\n"), 0644)

	outFile := filepath.Join(dir, "app.conf")
	err := renderTemplate(tmplFile, "prod", outFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(outFile)
	if !strings.Contains(string(data), "server=example.com:8080") {
		t.Errorf("expected rendered output, got: %s", data)
	}
}

func TestRenderTemplateUnknownVarPreserved(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVOY_DIR", dir)

	writeTemplateProfile(t, dir, "dev", "HOST=localhost\n")

	tmplFile := filepath.Join(dir, "tmpl.txt")
	os.WriteFile(tmplFile, []byte("{{HOST}} {{UNKNOWN}}"), 0644)

	outFile := filepath.Join(dir, "out.txt")
	err := renderTemplate(tmplFile, "dev", outFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, _ := os.ReadFile(outFile)
	if !strings.Contains(string(data), "{{UNKNOWN}}") {
		t.Errorf("expected {{UNKNOWN}} preserved, got: %s", data)
	}
}

func TestRenderTemplateNonExistentProfile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVOY_DIR", dir)
	_ = os.MkdirAll(filepath.Join(dir, ".envoy"), 0755)

	tmplFile := filepath.Join(dir, "tmpl.txt")
	os.WriteFile(tmplFile, []byte("{{KEY}}"), 0644)

	err := renderTemplate(tmplFile, "ghost", "")
	if err == nil {
		t.Fatal("expected error for non-existent profile")
	}
}

func TestRenderTemplateNonExistentTemplateFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("ENVOY_DIR", dir)
	writeTemplateProfile(t, dir, "dev", "KEY=val\n")

	err := renderTemplate(filepath.Join(dir, "missing.tmpl"), "dev", "")
	if err == nil {
		t.Fatal("expected error for missing template file")
	}
}
