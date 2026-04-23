package cmd

import (
	"strings"
)

// reorderLines takes a slice of env file lines and a desired key order,
// returning lines with ordered keys first, followed by the rest.
func reorderLines(lines []string, order []string) []string {
	keyLineMap := make(map[string]string)
	var remaining []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			remaining = append(remaining, line)
			continue
		}
		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			if containsKey(order, key) {
				keyLineMap[key] = line
				continue
			}
		}
		remaining = append(remaining, line)
	}

	var result []string
	for _, key := range order {
		if line, ok := keyLineMap[key]; ok {
			result = append(result, line)
		}
	}
	return append(result, remaining...)
}

// keysInOrder returns the list of keys from lines that appear in order slice.
func keysInOrder(lines []string) []string {
	var keys []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) == 2 {
			keys = append(keys, strings.TrimSpace(parts[0]))
		}
	}
	return keys
}
