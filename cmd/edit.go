package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "edit [profile]",
		Short: "Open a profile in your default editor",
		Args:  cobra.ExactArgs(1),
		Run:   editProfile,
	})
}

func editProfile(cmd *cobra.Command, args []string) {
	name := args[0]
	path := profilePath(name)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fatalf("profile %q does not exist", name)
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		editor = "vi"
	}

	c := exec.Command(editor, path)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		fatalf("editor exited with error: %v", err)
	}

	fmt.Printf("Profile %q saved.\n", name)
}
