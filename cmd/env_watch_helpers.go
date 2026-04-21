package cmd

import (
	"strings"
)

// watchDiffResult holds the categorized changes between two env file snapshots.
type watchDiffResult struct {
	Added   map[string]string
	Removed map[string]string
	Changed map[string][2]string // key -> [oldVal, newVal]
}

// computeWatchDiff compares two env content strings and returns structured diff.
func computeWatchDiff(oldContent, newContent string) watchDiffResult {
	result := watchDiffResult{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string][2]string),
	}

	oldMap := parseEnvContent(oldContent)
	newMap := parseEnvContent(newContent)

	for k, nv := range newMap {
		if ov, exists := oldMap[k]; !exists {
			result.Added[k] = nv
		} else if ov != nv {
			result.Changed[k] = [2]string{ov, nv}
		}
	}
	for k, ov := range oldMap {
		if _, exists := newMap[k]; !exists {
			result.Removed[k] = ov
		}
	}
	return result
}

// parseEnvContent parses env file content into a key-value map, skipping comments.
func parseEnvContent(content string) map[string]string {
	m := make(map[string]string)
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.Index(line, "=")
		if idx < 1 {
			continue
		}
		m[line[:idx]] = line[idx+1:]
	}
	return m
}

// hasWatchChanges returns true if the diff contains any changes.
func hasWatchChanges(d watchDiffResult) bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}
