package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func loadSchema(dir string) (*EnvSchema, error) {
	path := schemaFilePath(dir)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &EnvSchema{Keys: map[string]SchemaKey{}}, nil
	}
	if err != nil {
		return nil, err
	}
	var schema EnvSchema
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, fmt.Errorf("invalid schema file: %w", err)
	}
	if schema.Keys == nil {
		schema.Keys = map[string]SchemaKey{}
	}
	return &schema, nil
}

func saveSchema(dir string, schema *EnvSchema) error {
	path := schemaFilePath(dir)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0644)
}

func defineSchemaKey(dir, key string, required bool, desc, def string) error {
	schema, err := loadSchema(dir)
	if err != nil {
		return err
	}
	schema.Keys[key] = SchemaKey{
		Required:    required,
		Description: desc,
		Default:     def,
	}
	return saveSchema(dir, schema)
}

func validateAgainstSchema(dir, profile string) []string {
	schema, err := loadSchema(dir)
	if err != nil {
		return []string{fmt.Sprintf("could not load schema: %v", err)}
	}
	if len(schema.Keys) == 0 {
		return nil
	}
	pPath := profilePath(dir, profile)
	data, err := os.ReadFile(pPath)
	if err != nil {
		return []string{fmt.Sprintf("could not read profile '%s': %v", profile, err)}
	}
	present := map[string]bool{}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if idx := strings.Index(line, "="); idx > 0 {
			present[strings.TrimSpace(line[:idx])] = true
		}
	}
	var errs []string
	for k, meta := range schema.Keys {
		if meta.Required && !present[k] {
			errs = append(errs, fmt.Sprintf("required key '%s' is missing", k))
		}
	}
	return errs
}
