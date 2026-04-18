package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	cloneCmd := &cobra.Command{
		Use:   "clone <profile> <new-name>",
		Short: "Clone a profile to a new project directory",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			targetDir, _ := cmd.Flags().GetString("target")
			if err := cloneProfile(args[0], args[1], targetDir); err != nil {
				fatalf("clone failed: %v", err)
			}
		},
	}
	cloneCmd.Flags().StringP("target", "t", "", "Target project directory (default: current project)")
	rootCmd.AddCommand(cloneCmd)
}

func cloneProfile(srcName, destName, targetDir string) error {
	srcPath := profilePath(srcName)
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("profile %q does not exist", srcName)
	}

	var destEnvoyDir string
	if targetDir == "" {
		destEnvoyDir = filepath.Join(projectDir(), ".envoy")
	} else {
		info, err := os.Stat(targetDir)
		if err != nil || !info.IsDir() {
			return fmt.Errorf("target directory %q does not exist or is not a directory", targetDir)
		}
		destEnvoyDir = filepath.Join(targetDir, ".envoy")
	}

	if err := os.MkdirAll(destEnvoyDir, 0755); err != nil {
		return fmt.Errorf("could not create target .envoy dir: %v", err)
	}

	destPath := filepath.Join(destEnvoyDir, destName+".env")
	if _, err := os.Stat(destPath); err == nil {
		return fmt.Errorf("profile %q already exists in target", destName)
	}

	if err := copyFile(srcPath, destPath); err != nil {
		return fmt.Errorf("could not clone profile: %v", err)
	}

	fmt.Printf("Cloned profile %q to %q in %s\n", srcName, destName, destEnvoyDir)
	return nil
}
