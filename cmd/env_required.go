package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	requiredCmd := &cobra.Command{
		Use:   "required <profile> <KEY1> [KEY2...]",
		Short: "Assert that required keys exist and are non-empty in a profile",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			profile := args[0]
			keys := args[1:]
			strict, _ := cmd.Flags().GetBool("strict")
			missing, empty := assertRequiredKeys(profile, keys, strict)
			if len(missing) > 0 {
				fmt.Fprintf(os.Stderr, "missing keys in profile %q: %s\n", profile, strings.Join(missing, ", "))
				os.Exit(1)
			}
			if len(empty) > 0 {
				fmt.Fprintf(os.Stderr, "empty keys in profile %q: %s\n", profile, strings.Join(empty, ", "))
				os.Exit(1)
			}
			fmt.Printf("all required keys present in profile %q\n", profile)
		},
	}
	requiredCmd.Flags().Bool("strict", false, "fail if any key exists but has an empty value")
	rootCmd.AddCommand(requiredCmd)
}

// assertRequiredKeys checks that all keys exist (and optionally are non-empty).
// Returns lists of missing keys and (if strict) empty-valued keys.
func assertRequiredKeys(profile string, keys []string, strict bool) (missing, empty []string) {
	path := profilePath(profile)
	envMap, err := readEnvMap(path)
	if err != nil {
		return keys, nil
	}
	for _, k := range keys {
		val, ok := envMap[k]
		if !ok {
			missing = append(missing, k)
		} else if strict && strings.TrimSpace(val) == "" {
			empty = append(empty, k)
		}
	}
	return missing, empty
}
