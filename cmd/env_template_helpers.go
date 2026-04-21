package cmd

import (
	"regexp"
	"sort"
)

// extractTemplateVars returns all unique {{VAR}} placeholders found in the template string.
func extractTemplateVars(tmpl string) []string {
	matches := templateVarRe.FindAllStringSubmatch(tmpl, -1)
	seen := make(map[string]bool)
	var vars []string
	for _, m := range matches {
		if len(m) > 1 && !seen[m[1]] {
			seen[m[1]] = true
			vars = append(vars, m[1])
		}
	}
	sort.Strings(vars)
	return vars
}

// missingTemplateKeys returns keys referenced in the template but absent from envMap.
func missingTemplateKeys(tmpl string, envMap map[string]string) []string {
	vars := extractTemplateVars(tmpl)
	var missing []string
	for _, v := range vars {
		if _, ok := envMap[v]; !ok {
			missing = append(missing, v)
		}
	}
	return missing
}

// templateVarPattern returns the compiled regex (exported for tests).
func templateVarPattern() *regexp.Regexp {
	return templateVarRe
}
