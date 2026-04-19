package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "lint [profile]",
		Short: "Lint a profile for common issues",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			issues, err := lintProfile(projectDir(), args[0])
			if err != nil {
				fatalf("lint error: %v", err)
			}
			if len(issues) == 0 {
				fmt.Println("No issues found.")
				return
			}
			for _, issue := range issues {
				fmt.Println(issue)
			}
			os.Exit(1)
		},
	})
}

func lintProfile(dir, name string) ([]string, error) {
	path := filepath.Join(dir, ".envoy", "profiles", name+".env")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("profile %q not found", name)
	}

	var issues []string
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		lineno := i + 1
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		idx := strings.Index(trimmed, "=")
		if idx < 0 {
			issues = append(issues, fmt.Sprintf("line %d: missing '=' in %q", lineno, trimmed))
			continue
		}
		key := strings.TrimSpace(trimmed[:idx])
		val := trimmed[idx+1:]
		if key == "" {
			issues = append(issues, fmt.Sprintf("line %d: empty key", lineno))
		}
		if strings.Contains(key, " ") {
			issues = append(issues, fmt.Sprintf("line %d: key %q contains whitespace", lineno, key))
		}
		if strings.HasPrefix(val, " ") || strings.HasSuffix(val, " ") {
			issues = append(issues, fmt.Sprintf("line %d: value for %q has leading/trailing whitespace", lineno, key))
		}
		if strings.Contains(val, "\t") {
			issues = append(issues, fmt.Sprintf("line %d: value for %q contains tab character", lineno, key))
		}
	}
	return issues, nil
}
