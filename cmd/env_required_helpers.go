package cmd

import "strings"

// requiredKeySummary returns a human-readable summary of the assertion result.
func requiredKeySummary(profile string, missing, empty []string) string {
	var sb strings.Builder
	if len(missing) == 0 && len(empty) == 0 {
		sb.WriteString("OK: all required keys are present")
		if profile != "" {
			sb.WriteString(" in profile " + profile)
		}
		return sb.String()
	}
	if len(missing) > 0 {
		sb.WriteString("MISSING: " + strings.Join(missing, ", "))
	}
	if len(empty) > 0 {
		if sb.Len() > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString("EMPTY: " + strings.Join(empty, ", "))
	}
	return sb.String()
}

// partitionKeys splits a slice of KEY=VALUE or bare KEY tokens into
// key-only entries (bare) and a map of expected values (KEY=VALUE).
func partitionKeys(tokens []string) (bareKeys []string, exact map[string]string) {
	exact = make(map[string]string)
	for _, t := range tokens {
		if idx := strings.IndexByte(t, '='); idx > 0 {
			exact[t[:idx]] = t[idx+1:]
		} else {
			bareKeys = append(bareKeys, t)
		}
	}
	return bareKeys, exact
}
