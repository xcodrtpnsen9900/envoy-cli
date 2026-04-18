package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize envoy in the current project directory",
	Long:  `Creates the .envoy directory structure and an optional default profile.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := initProject(projectDir); err != nil {
			fatalf("init failed: %v", err)
		}
		fmt.Println("Initialized envoy in", projectDir)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

// initProject sets up the .envoy directory and a default profile if none exist.
func initProject(dir string) error {
	envoyDir := filepath.Join(dir, ".envoy")
	if _, err := os.Stat(envoyDir); err == nil {
		return fmt.Errorf(".envoy directory already exists")
	}
	if err := os.MkdirAll(envoyDir, 0755); err != nil {
		return fmt.Errorf("could not create .envoy directory: %w", err)
	}

	defaultProfile := filepath.Join(envoyDir, "default.env")
	if err := os.WriteFile(defaultProfile, []byte("# Default profile\n"), 0644); err != nil {
		return fmt.Errorf("could not create default profile: %w", err)
	}

	activeFile := filepath.Join(envoyDir, ".active")
	if err := os.WriteFile(activeFile, []byte("default"), 0644); err != nil {
		return fmt.Errorf("could not set active profile: %w", err)
	}

	return nil
}
