package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "touch [profile]",
		Short: "Create an empty profile if it does not already exist",
		Args:  cobra.ExactArgs(1),
		Run:   touchProfile,
	})
}

func touchProfile(cmd *cobra.Command, args []string) {
	name := args[0]
	path := profilePath(name)

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fatalf("failed to create profiles directory: %v", err)
	}

	if _, err := os.Stat(path); err == nil {
		fmt.Printf("Profile %q already exists, skipping.\n", name)
		return
	}

	f, err := os.Create(path)
	if err != nil {
		fatalf("failed to create profile: %v", err)
	}
	f.Close()

	fmt.Printf("Created empty profile %q.\n", name)
}
