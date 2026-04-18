package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete [profile]",
	Short: "Delete an existing .env profile",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		if err := deleteProfile(projectDir, name); err != nil {
			fatalf("error deleting profile: %v", err)
		}
		fmt.Printf("Profile '%s' deleted.\n", name)
	},
}

func deleteProfile(dir, name string) error {
	active, err := activeProfile(dir)
	if err == nil && active == name {
		return fmt.Errorf("cannot delete the active profile '%s'; switch to another profile first", name)
	}

	path := profilePath(dir, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("profile '%s' does not exist", name)
	}

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to remove profile file: %w", err)
	}

	// Remove from active file if it somehow references this profile
	activeFile := filepath.Join(dir, ".envoy", "active")
	content, err := os.ReadFile(activeFile)
	if err == nil && string(content) == name {
		os.WriteFile(activeFile, []byte(""), 0644)
	}

	return nil
}
