package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	renameKeyCmd := &cobra.Command{
		Use:   "rename-key [profile] [old-key] [new-key]",
		Short: "Rename a key within a profile",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			profile, oldKey, newKey := args[0], args[1], args[2]
			if err := renameProfileKey(profile, oldKey, newKey); err != nil {
				fatalf("rename-key: %v", err)
			}
			fmt.Printf("Renamed key '%s' to '%s' in profile '%s'\n", oldKey, newKey, profile)
		},
	}
	rootCmd.AddCommand(renameKeyCmd)
}

// renameProfileKey renames oldKey to newKey inside the given profile file.
// It returns an error if the profile does not exist, oldKey is not found,
// or newKey already exists.
func renameProfileKey(profile, oldKey, newKey string) error {
	path := profilePath(profile)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("profile '%s' does not exist", profile)
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var lines []string
	foundOld := false
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			lines = append(lines, line)
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		if key == newKey {
			return fmt.Errorf("key '%s' already exists in profile '%s'", newKey, profile)
		}
		if key == oldKey {
			foundOld = true
			value := ""
			if len(parts) == 2 {
				value = parts[1]
			}
			lines = append(lines, newKey+"="+value)
		} else {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if !foundOld {
		return fmt.Errorf("key '%s' not found in profile '%s'", oldKey, profile)
	}

	out := strings.Join(lines, "\n")
	if len(lines) > 0 {
		out += "\n"
	}
	return os.WriteFile(path, []byte(out), 0644)
}
