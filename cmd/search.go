package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	var caseSensitive bool

	searchCmd := &cobra.Command{
		Use:   "search [profile] [key]",
		Short: "Search for a key in a profile",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			results, err := searchProfile(projectDir, args[0], args[1], caseSensitive)
			if err != nil {
				fatalf("%v", err)
			}
			if len(results) == 0 {
				fmt.Printf("No matches found for '%s' in profile '%s'\n", args[1], args[0])
				return
			}
			for _, r := range results {
				fmt.Println(r)
			}
		},
	}

	searchCmd.Flags().BoolVarP(&caseSensitive, "case-sensitive", "c", false, "Enable case-sensitive search")
	rootCmd.AddCommand(searchCmd)
}

func searchProfile(dir, profile, key string, caseSensitive bool) ([]string, error) {
	path := filepath.Join(dir, ".envoy", profile+".env")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("profile '%s' does not exist", profile)
		}
		return nil, err
	}

	var matches []string
	lines := strings.Split(string(data), "\n")
	searchKey := key
	if !caseSensitive {
		searchKey = strings.ToLower(key)
	}

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		compare := trimmed
		if !caseSensitive {
			compare = strings.ToLower(trimmed)
		}
		if strings.Contains(compare, searchKey) {
			matches = append(matches, fmt.Sprintf("line %d: %s", i+1, trimmed))
		}
	}
	return matches, nil
}
