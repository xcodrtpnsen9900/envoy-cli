package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	var outputProfile string

	cmd := &cobra.Command{
		Use:   "placeholder <profile>",
		Short: "Find keys with placeholder or empty values in a profile",
		Long:  `Scans a profile for keys whose values are placeholders like CHANGEME, TODO, FIXME, or empty strings.`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			results, err := findPlaceholders(projectDir(), args[0])
			if err != nil {
				fatalf("placeholder: %v", err)
			}
			if len(results) == 0 {
				fmt.Println("No placeholder values found.")
				return
			}
			if outputProfile != "" {
				if err := writePlaceholderReport(projectDir(), outputProfile, results); err != nil {
					fatalf("placeholder: %v", err)
				}
				fmt.Printf("Report written to profile: %s\n", outputProfile)
				return
			}
			for _, r := range results {
				fmt.Printf("  line %d: %s = %q\n", r.Line, r.Key, r.Value)
			}
		},
	}

	cmd.Flags().StringVarP(&outputProfile, "output", "o", "", "Write report keys to a new profile")
	rootCmd.AddCommand(cmd)
}

type placeholderResult struct {
	Line  int
	Key   string
	Value string
}

var placeholderTokens = []string{"CHANGEME", "TODO", "FIXME", "PLACEHOLDER", "YOUR_", "<", ">"}

func findPlaceholders(root, profile string) ([]placeholderResult, error) {
	path := profilePath(root, profile)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("profile %q not found", profile)
	}
	var results []placeholderResult
	for i, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.Index(line, "=")
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		if isPlaceholderValue(val) {
			results = append(results, placeholderResult{Line: i + 1, Key: key, Value: val})
		}
	}
	return results, nil
}

func isPlaceholderValue(val string) bool {
	if val == "" {
		return true
	}
	upper := strings.ToUpper(val)
	for _, token := range placeholderTokens {
		if strings.Contains(upper, strings.ToUpper(token)) {
			return true
		}
	}
	return false
}

func writePlaceholderReport(root, profile string, results []placeholderResult) error {
	path := profilePath(root, profile)
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("profile %q already exists", profile)
	}
	var sb strings.Builder
	sb.WriteString("# Placeholder report\n")
	for _, r := range results {
		fmt.Fprintf(&sb, "%s=%s\n", r.Key, r.Value)
	}
	return os.WriteFile(path, []byte(sb.String()), 0644)
}
