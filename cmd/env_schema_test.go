package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func setupSchemaDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".envoy"), 0755); err != nil {
		t.Fatal(err)
	}
	return dir
}

func TestDefineAndLoadSchema(t *testing.T) {
	dir := setupSchemaDir(t)
	if err := defineSchemaKey(dir, "DB_URL", true, "Database URL", ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	schema, err := loadSchema(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	key, ok := schema.Keys["DB_URL"]
	if !ok {
		t.Fatal("expected DB_URL in schema")
	}
	if !key.Required {
		t.Error("expected DB_URL to be required")
	}
	if key.Description != "Database URL" {
		t.Errorf("expected description 'Database URL', got '%s'", key.Description)
	}
}

func TestLoadSchemaMissingFile(t *testing.T) {
	dir := setupSchemaDir(t)
	schema, err := loadSchema(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(schema.Keys) != 0 {
		t.Errorf("expected empty schema, got %d keys", len(schema.Keys))
	}
}

func TestValidateAgainstSchemaAllPresent(t *testing.T) {
	dir := setupSchemaDir(t)
	_ = defineSchemaKey(dir, "API_KEY", true, "", "")
	p := filepath.Join(dir, ".envoy", "production.env")
	_ = os.WriteFile(p, []byte("API_KEY=secret\n"), 0644)
	errs := validateAgainstSchema(dir, "production")
	if len(errs) != 0 {
		t.Errorf("expected no errors, got: %v", errs)
	}
}

func TestValidateAgainstSchemaMissingRequired(t *testing.T) {
	dir := setupSchemaDir(t)
	_ = defineSchemaKey(dir, "SECRET_KEY", true, "", "")
	p := filepath.Join(dir, ".envoy", "staging.env")
	_ = os.WriteFile(p, []byte("OTHER=value\n"), 0644)
	errs := validateAgainstSchema(dir, "staging")
	if len(errs) == 0 {
		t.Fatal("expected validation error for missing required key")
	}
}

func TestValidateAgainstSchemaOptionalKeyMissing(t *testing.T) {
	dir := setupSchemaDir(t)
	_ = defineSchemaKey(dir, "LOG_LEVEL", false, "Logging level", "info")
	p := filepath.Join(dir, ".envoy", "dev.env")
	_ = os.WriteFile(p, []byte("DB_HOST=localhost\n"), 0644)
	errs := validateAgainstSchema(dir, "dev")
	if len(errs) != 0 {
		t.Errorf("expected no errors for optional key, got: %v", errs)
	}
}
