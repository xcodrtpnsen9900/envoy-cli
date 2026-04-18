package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the currently active profile",
	Run: func(cmd *cobra.Command, args []string) {
		active, err := activeProfile(projectDir)
		if err != nil {
			fatalf("could not get active profile: %v", err)
		}
		fmt.Printf("Active profile: %s\n", active)
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

// activeProfile reads the currently active profile name from .envoy/.active.
func activeProfile(dir string) (string, error) {
	activeFile := filepath.Join(dir, ".envoy", ".active")
	data, err := os.ReadFile(activeFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("no active profile found; run 'envoy init' first")
		}
		return "", err
	}
	name := string(data)
	if name == "" {
		return "", fmt.Errorf("active profile file is empty")
	}
	return name, nil
}
