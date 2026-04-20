package cmd

import (
	"strings"
)

// extractKey returns the key portion of a KEY=VALUE line.
// Returns empty string if the line has no '=' separator.
func extractKey(line string) string {
	parts := strings.SplitN(line, "=", 2)
	if len(parts) < 2 {
		return ""
	}
	return strings.TrimSpace(parts[0])
}

// isSorted reports whether the given slice of kv lines is already
// sorted alphabetically (case-insensitive) by key.
func isSorted(lines []string) bool {
	for i := 1; i < len(lines); i++ {
		prev := strings.ToLower(extractKey(lines[i-1]))
		curr := strings.ToLower(extractKey(lines[i]))
		if prev > curr {
			return false
		}
	}
	return true
}

// deduplicateKeys removes duplicate KEY entries from a slice of kv lines,
// keeping the last occurrence of each key.
func deduplicateKeys(lines []string) []string {
	seen := make(map[string]int)
	for i, line := range lines {
		key := extractKey(line)
		if key != "" {
			seen[key] = i
		}
	}
	result := make([]string, 0, len(seen))
	for i, line := range lines {
		key := extractKey(line)
		if key == "" {
			continue
		}
		if seen[key] == i {
			result = append(result, line)
		}
	}
	return result
}
