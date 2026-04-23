package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	reorderCmd := &cobra.Command{
		Use:   "reorder <profile> <key1,key2,...>",
		Short: "Reorder keys in a profile to a specified order",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			profile := args[0]
			order := strings.Split(args[1], ",")
			if err := reorderProfile(projectDir(), profile, order); err != nil {
				fatalf("reorder: %v", err)
			}
			fmt.Printf("Profile '%s' keys reordered.\n", profile)
		},
	}
	rootCmd.AddCommand(reorderCmd)
}

// reorderProfile moves the specified keys to the top of the profile file,
// preserving all other lines (including comments) after them.
func reorderProfile(root, profile string, order []string) error {
	path := profilePath(root, profile)
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("profile '%s' not found", profile)
	}

	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	keyLineMap := make(map[string]string)
	var remaining []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			remaining = append(remaining, line)
			continue
		}
		parts := strings.SplitN(trimmed, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			if containsKey(order, key) {
				keyLineMap[key] = line
				continue
			}
		}
		remaining = append(remaining, line)
	}

	var result []string
	for _, key := range order {
		if line, ok := keyLineMap[key]; ok {
			result = append(result, line)
		}
	}
	result = append(result, remaining...)

	return os.WriteFile(path, []byte(strings.Join(result, "\n")+"\n"), 0644)
}

func containsKey(order []string, key string) bool {
	for _, k := range order {
		if strings.TrimSpace(k) == key {
			return true
		}
	}
	return false
}
