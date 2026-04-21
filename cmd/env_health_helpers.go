package cmd

import "sort"

// healthSummary returns a human-readable one-line summary of a health report.
func healthSummary(r *healthReport) string {
	if r.OK {
		return "healthy"
	}
	issues := []string{}
	if len(r.EmptyVals) > 0 {
		issues = append(issues, "empty values")
	}
	if len(r.Duplicates) > 0 {
		issues = append(issues, "duplicate keys")
	}
	result := ""
	for i, s := range issues {
		if i > 0 {
			result += ", "
		}
		result += s
	}
	return result
}

// sortedKeys returns the keys of a string map in sorted order.
func sortedKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// healthIssueCount returns the total number of issues found in a report.
func healthIssueCount(r *healthReport) int {
	return len(r.EmptyVals) + len(r.Duplicates)
}
