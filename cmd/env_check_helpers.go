package cmd

import (
	"fmt"
	"strings"
)

// requiredKeysFromString parses a comma-separated string of required keys.
func requiredKeysFromString(s string) []string {
	parts := strings.Split(s, ",")
	var keys []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			keys = append(keys, p)
		}
	}
	return keys
}

// formatMissingKeys returns a human-readable summary of missing keys.
func formatMissingKeys(profile string, missing []string) string {
	if len(missing) == 0 {
		return fmt.Sprintf("Profile '%s' has all required keys.", profile)
	}
	return fmt.Sprintf("Profile '%s' is missing: %s", profile, strings.Join(missing, ", "))
}
