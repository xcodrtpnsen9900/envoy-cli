package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "rename [old] [new]",
		Short: "Rename an existing profile",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if err := renameProfile(args[0], args[1]); err != nil {
				fatalf("rename failed: %v", err)
			}
			fmt.Printf("Profile '%s' renamed to '%s'\n", args[0], args[1])
		},
	})
}

func renameProfile(oldName, newName string) error {
	oldPath := profilePath(oldName)
	newPath := profilePath(newName)

	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return fmt.Errorf("profile '%s' does not exist", oldName)
	}

	if _, err := os.Stat(newPath); err == nil {
		return fmt.Errorf("profile '%s' already exists", newName)
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("could not rename profile file: %w", err)
	}

	active, err := activeProfile()
	if err == nil && active == oldName {
		envLink := filepath.Join(projectDir, ".env")
		os.Remove(envLink)
		if err := os.Symlink(newPath, envLink); err != nil {
			return fmt.Errorf("could not update .env symlink: %w", err)
		}
		activeFile := filepath.Join(projectDir, ".envoy", "active")
		if err := os.WriteFile(activeFile, []byte(newName), 0644); err != nil {
			return fmt.Errorf("could not update active profile: %w", err)
		}
	}

	return nil
}
