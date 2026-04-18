package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	var copyCmd = &cobra.Command{
		Use:   "copy <source> <destination>",
		Short: "Copy an existing profile to a new profile",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if err := copyProfile(projectDir, args[0], args[1]); err != nil {
				fatalf("Error: %v", err)
			}
			fmt.Printf("Profile '%s' copied to '%s'\n", args[0], args[1])
		},
	}
	rootCmd.AddCommand(copyCmd)
}

func copyProfile(dir, src, dst string) error {
	srcPath := profilePath(dir, src)
	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("profile '%s' does not exist", src)
	}

	dstPath := profilePath(dir, dst)
	if _, err := os.Stat(dstPath); err == nil {
		return fmt.Errorf("profile '%s' already exists", dst)
	}

	if err := copyFile(srcPath, dstPath); err != nil {
		return fmt.Errorf("failed to copy profile: %w", err)
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
