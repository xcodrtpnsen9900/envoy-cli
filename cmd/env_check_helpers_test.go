package cmd

import (
	"strings"
	"testing"
)

func TestRequiredKeysFromString(t *testing.T) {
	keys := requiredKeysFromString("DB_HOST, DB_PORT , API_KEY")
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	if keys[1] != "DB_PORT" {
		t.Errorf("expected DB_PORT, got %s", keys[1])
	}
}

func TestRequiredKeysFromStringEmpty(t *testing.T) {
	keys := requiredKeysFromString("")
	if len(keys) != 0 {
		t.Errorf("expected empty slice, got %v", keys)
	}
}

func TestFormatMissingKeysNone(t *testing.T) {
	out := formatMissingKeys("prod", nil)
	if !strings.Contains(out, "all required keys") {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestFormatMissingKeysSome(t *testing.T) {
	out := formatMissingKeys("prod", []string{"DB_HOST", "API_KEY"})
	if !strings.Contains(out, "DB_HOST") || !strings.Contains(out, "API_KEY") {
		t.Errorf("unexpected output: %s", out)
	}
}
