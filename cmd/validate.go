package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	validateCmd := &cobra.Command{
		Use:   "validate <profile>",
		Short: "Validate an env profile for syntax errors",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := validateProfile(args[0]); err != nil {
				fatalf("validation failed: %v", err)
			}
			fmt.Printf("Profile %q is valid.\n", args[0])
		},
	}
	rootCmd.AddCommand(validateCmd)
}

func validateProfile(name string) error {
	dir := projectDir()
	path := profilePath(dir, name)

	lines, err := readLines(path)
	if err != nil {
		return fmt.Errorf("cannot read profile %q: %w", name, err)
	}

	var errs []string
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if !strings.Contains(trimmed, "=") {
			errs = append(errs, fmt.Sprintf("line %d: missing '=' in %q", i+1, trimmed))
			continue
		}
		parts := strings.SplitN(trimmed, "=", 2)
		key := strings.TrimSpace(parts[0])
		if key == "" {
			errs = append(errs, fmt.Sprintf("line %d: empty key", i+1))
		}
		if strings.ContainsAny(key, " \t") {
			errs = append(errs, fmt.Sprintf("line %d: key %q contains whitespace", i+1, key))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%d issue(s) found:\n  %s", len(errs), strings.Join(errs, "\n  "))
	}
	return nil
}

func readLines(path string) ([]string, error) {
	data, err := readEnvMap(path)
	_ = data
	// reuse file reading via os directly
	import_workaround_use_os_readfile(path)
	return nil, err
}
