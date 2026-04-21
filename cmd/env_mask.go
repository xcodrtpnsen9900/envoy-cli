package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	var reveal bool
	var keys []string

	maskCmd := &cobra.Command{
		Use:   "mask <profile>",
		Short: "Display a profile with sensitive values masked",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			path := profilePath(args[0])
			if err := maskProfile(path, keys, reveal); err != nil {
				fatalf("%v", err)
			}
		},
	}

	maskCmd.Flags().BoolVar(&reveal, "reveal", false, "Show full values (no masking)")
	maskCmd.Flags().StringSliceVar(&keys, "keys", defaultSensitiveKeys(), "Keys to mask (comma-separated)")

	rootCmd.AddCommand(maskCmd)
}

func defaultSensitiveKeys() []string {
	return []string{"PASSWORD", "SECRET", "TOKEN", "KEY", "API", "PRIVATE", "PASS", "CREDENTIAL"}
}

func isSensitiveKey(key string, sensitiveKeys []string) bool {
	upper := strings.ToUpper(key)
	for _, s := range sensitiveKeys {
		if strings.Contains(upper, strings.ToUpper(s)) {
			return true
		}
	}
	return false
}

func maskValue(value string) string {
	if len(value) == 0 {
		return ""
	}
	if len(value) <= 4 {
		return strings.Repeat("*", len(value))
	}
	return value[:2] + strings.Repeat("*", len(value)-4) + value[len(value)-2:]
}

func maskProfile(path string, sensitiveKeys []string, reveal bool) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("profile not found: %s", path)
	}

	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			fmt.Println(line)
			continue
		}
		idx := strings.Index(line, "=")
		if idx < 0 {
			fmt.Println(line)
			continue
		}
		key := line[:idx]
		value := line[idx+1:]
		if !reveal && isSensitiveKey(key, sensitiveKeys) {
			fmt.Printf("%s=%s\n", key, maskValue(value))
		} else {
			fmt.Println(line)
		}
	}
	return nil
}
