package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(renameCmd)
}

var renameCmd = &cobra.Command{
	Use:   "rename [old] [new]",
	Short: "Rename an existing .env profile",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		oldName, newName := args[0], args[1]
		if err := renameProfile(projectDir, oldName, newName); err != nil {
			fatalf("error renaming profile: %v", err)
		}
		fmt.Printf("Profile '%s' renamed to '%s'.\n", oldName, newName)
	},
}

func renameProfile(dir, oldName, newName string) error {
	oldPath := profilePath(dir, oldName)
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return fmt.Errorf("profile '%s' does not exist", oldName)
	}

	newPath := profilePath(dir, newName)
	if _, err := os.Stat(newPath); err == nil {
		return fmt.Errorf("profile '%s' already exists", newName)
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to rename profile: %w", err)
	}

	// Update active file if the renamed profile was active
	activeFile := filepath.Join(dir, ".envoy", "active")
	content, err := os.ReadFile(activeFile)
	if err == nil && string(content) == oldName {
		if err := os.WriteFile(activeFile, []byte(newName), 0644); err != nil {
			return fmt.Errorf("failed to update active profile reference: %w", err)
		}
	}

	return nil
}
