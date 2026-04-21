package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndReloadSchema(t *testing.T) {
	dir := t.TempDir()
	_ = os.MkdirAll(filepath.Join(dir, ".envoy"), 0755)
	schema := &EnvSchema{
		Keys: map[string]SchemaKey{
			"PORT":    {Required: true, Description: "Server port", Default: "8080"},
			"DEBUG":   {Required: false, Description: "Debug mode"},
		},
	}
	if err := saveSchema(dir, schema); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	loaded, err := loadSchema(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(loaded.Keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(loaded.Keys))
	}
	port := loaded.Keys["PORT"]
	if port.Default != "8080" {
		t.Errorf("expected default '8080', got '%s'", port.Default)
	}
}

func TestDefineSchemaKeyOverwrite(t *testing.T) {
	dir := t.TempDir()
	_ = os.MkdirAll(filepath.Join(dir, ".envoy"), 0755)
	_ = defineSchemaKey(dir, "TOKEN", true, "Auth token", "")
	_ = defineSchemaKey(dir, "TOKEN", false, "Updated description", "default-token")
	schema, _ := loadSchema(dir)
	key := schema.Keys["TOKEN"]
	if key.Required {
		t.Error("expected Required to be false after overwrite")
	}
	if key.Default != "default-token" {
		t.Errorf("expected default 'default-token', got '%s'", key.Default)
	}
	if key.Description != "Updated description" {
		t.Errorf("unexpected description: %s", key.Description)
	}
}

func TestValidateAgainstSchemaNoSchema(t *testing.T) {
	dir := t.TempDir()
	_ = os.MkdirAll(filepath.Join(dir, ".envoy"), 0755)
	p := filepath.Join(dir, ".envoy", "local.env")
	_ = os.WriteFile(p, []byte("FOO=bar\n"), 0644)
	errs := validateAgainstSchema(dir, "local")
	if len(errs) != 0 {
		t.Errorf("expected no errors with empty schema, got: %v", errs)
	}
}

func TestValidateAgainstSchemaNonExistentProfile(t *testing.T) {
	dir := t.TempDir()
	_ = os.MkdirAll(filepath.Join(dir, ".envoy"), 0755)
	_ = defineSchemaKey(dir, "KEY", true, "", "")
	errs := validateAgainstSchema(dir, "ghost")
	if len(errs) == 0 {
		t.Fatal("expected error for non-existent profile")
	}
}
