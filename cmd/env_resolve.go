package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	var outputFormat string

	resolveCmd := &cobra.Command{
		Use:   "resolve <profile>",
		Short: "Resolve a profile by expanding values that reference other env vars",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			profile := args[0]
			resolved, err := resolveProfile(profile)
			if err != nil {
				fatalf("resolve: %v", err)
			}
			if outputFormat == "export" {
				for _, kv := range resolved {
					fmt.Printf("export %s\n", kv)
				}
			} else {
				for _, kv := range resolved {
					fmt.Println(kv)
				}
			}
		},
	}

	resolveCmd.Flags().StringVarP(&outputFormat, "format", "f", "", "Output format: export")
	rootCmd.AddCommand(resolveCmd)
}

// resolveProfile reads a profile and expands any values referencing ${VAR} or $VAR
// using other keys defined within the same profile (in order).
func resolveProfile(profile string) ([]string, error) {
	path := profilePath(profile)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("profile %q not found", profile)
	}

	env := map[string]string{}
	var lines []string

	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		idx := strings.IndexByte(trimmed, '=')
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(trimmed[:idx])
		val := strings.TrimSpace(trimmed[idx+1:])
		resolved := expandVars(val, env)
		env[key] = resolved
		lines = append(lines, key+"="+resolved)
	}

	return lines, nil
}

// expandVars replaces ${VAR} and $VAR occurrences in s using the provided map.
func expandVars(s string, env map[string]string) string {
	return os.Expand(s, func(key string) string {
		if v, ok := env[key]; ok {
			return v
		}
		// Fall back to actual OS environment
		return os.Getenv(key)
	})
}
