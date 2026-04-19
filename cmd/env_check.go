package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	envCheckCmd := &cobra.Command{
		Use:   "env-check <profile> <required-keys...>",
		Short: "Check that a profile contains all required keys",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			profile := args[0]
			required := args[1:]
			missing, err := checkRequiredKeys(profile, required)
			if err != nil {
				fatalf("error: %v", err)
			}
			if len(missing) > 0 {
				fmt.Fprintf(os.Stderr, "Missing keys in profile '%s':\n", profile)
				for _, k := range missing {
					fmt.Fprintf(os.Stderr, "  - %s\n", k)
				}
				os.Exit(1)
			}
			fmt.Printf("Profile '%s' contains all required keys.\n", profile)
		},
	}
	rootCmd.AddCommand(envCheckCmd)
}

func checkRequiredKeys(profile string, required []string) ([]string, error) {
	path := profilePath(profile)
	envMap, err := readEnvMap(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("profile '%s' not found", profile)
		}
		return nil, err
	}
	var missing []string
	for _, k := range required {
		k = strings.TrimSpace(k)
		if _, ok := envMap[k]; !ok {
			missing = append(missing, k)
		}
	}
	return missing, nil
}
