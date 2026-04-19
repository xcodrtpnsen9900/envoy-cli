package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	pinCmd := &cobra.Command{
		Use:   "pin [profile]",
		Short: "Pin a profile to prevent accidental switching or deletion",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := pinProfile(projectDir, args[0]); err != nil {
				fatalf("%v", err)
			}
		},
	}

	unpinCmd := &cobra.Command{
		Use:   "unpin [profile]",
		Short: "Unpin a previously pinned profile",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := unpinProfile(projectDir, args[0]); err != nil {
				fatalf("%v", err)
			}
		},
	}

	rootCmd.AddCommand(pinCmd)
	rootCmd.AddCommand(unpinCmd)
}

func pinFilePath(dir string) string {
	return filepath.Join(dir, ".envoy", "pinned")
}

func pinProfile(dir, name string) error {
	profiles, err := listProfiles(dir)
	if err != nil || !contains(profiles, name) {
		return fmt.Errorf("profile %q does not exist", name)
	}
	pinned, _ := readPinned(dir)
	for _, p := range pinned {
		if p == name {
			fmt.Printf("Profile %q is already pinned.\n", name)
			return nil
		}
	}
	pinned = append(pinned, name)
	if err := writePinned(dir, pinned); err != nil {
		return err
	}
	fmt.Printf("Profile %q pinned.\n", name)
	return nil
}

func unpinProfile(dir, name string) error {
	pinned, _ := readPinned(dir)
	newList := []string{}
	found := false
	for _, p := range pinned {
		if p == name {
			found = true
			continue
		}
		newList = append(newList, p)
	}
	if !found {
		return fmt.Errorf("profile %q is not pinned", name)
	}
	if err := writePinned(dir, newList); err != nil {
		return err
	}
	fmt.Printf("Profile %q unpinned.\n", name)
	return nil
}

func readPinned(dir string) ([]string, error) {
	data, err := os.ReadFile(pinFilePath(dir))
	if err != nil {
		return []string{}, nil
	}
	return splitLines(string(data)), nil
}

func writePinned(dir string, profiles []string) error {
	content := ""
	for _, p := range profiles {
		content += p + "\n"
	}
	return os.WriteFile(pinFilePath(dir), []byte(content), 0644)
}

func isPinned(dir, name string) bool {
	pinned, _ := readPinned(dir)
	return contains(pinned, name)
}

func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}
