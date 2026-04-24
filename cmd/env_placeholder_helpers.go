package cmd

import (
	"fmt"
	"strings"
)

// placeholderSummary returns a human-readable summary of placeholder findings.
func placeholderSummary(results []placeholderResult) string {
	if len(results) == 0 {
		return "No placeholders found."
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d placeholder(s) found:\n", len(results))
	for _, r := range results {
		fmt.Fprintf(&sb, "  [line %d] %s = %q\n", r.Line, r.Key, r.Value)
	}
	return sb.String()
}

// placeholderKeys returns just the key names from a result slice.
func placeholderKeys(results []placeholderResult) []string {
	keys := make([]string, 0, len(results))
	for _, r := range results {
		keys = append(keys, r.Key)
	}
	return keys
}

// filterPlaceholdersByToken returns only results whose value contains the given token (case-insensitive).
func filterPlaceholdersByToken(results []placeholderResult, token string) []placeholderResult {
	var out []placeholderResult
	upper := strings.ToUpper(token)
	for _, r := range results {
		if strings.Contains(strings.ToUpper(r.Value), upper) {
			out = append(out, r)
		}
	}
	return out
}

// countEmptyPlaceholders returns the number of results with an empty value.
func countEmptyPlaceholders(results []placeholderResult) int {
	count := 0
	for _, r := range results {
		if r.Value == "" {
			count++
		}
	}
	return count
}
