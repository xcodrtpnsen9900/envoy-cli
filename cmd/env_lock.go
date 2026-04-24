package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

func lockFilePath(dir, profile string) string {
	return filepath.Join(dir, ".envoy", "locks", profile+".lock")
}

func init() {
	lockCmd := &cobra.Command{
		Use:   "lock [profile]",
		Short: "Lock a profile to prevent modifications",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := lockProfile(projectDir, args[0]); err != nil {
				fatalf("lock: %v", err)
			}
		},
	}

	unlockCmd := &cobra.Command{
		Use:   "unlock [profile]",
		Short: "Unlock a previously locked profile",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := unlockProfile(projectDir, args[0]); err != nil {
				fatalf("unlock: %v", err)
			}
		},
	}

	listLockedCmd := &cobra.Command{
		Use:   "list-locked",
		Short: "List all locked profiles",
		Run: func(cmd *cobra.Command, args []string) {
			if err := listLockedProfiles(projectDir); err != nil {
				fatalf("list-locked: %v", err)
			}
		},
	}

	rootCmd.AddCommand(lockCmd)
	rootCmd.AddCommand(unlockCmd)
	rootCmd.AddCommand(listLockedCmd)
}

func lockProfile(dir, profile string) error {
	pPath := profilePath(dir, profile)
	if _, err := os.Stat(pPath); os.IsNotExist(err) {
		return fmt.Errorf("profile %q does not exist", profile)
	}
	if isProfileLocked(dir, profile) {
		return fmt.Errorf("profile %q is already locked", profile)
	}
	lockDir := filepath.Join(dir, ".envoy", "locks")
	if err := os.MkdirAll(lockDir, 0755); err != nil {
		return fmt.Errorf("could not create locks directory: %w", err)
	}
	content := fmt.Sprintf("locked_at=%s\n", time.Now().Format(time.RFC3339))
	if err := os.WriteFile(lockFilePath(dir, profile), []byte(content), 0644); err != nil {
		return fmt.Errorf("could not write lock file: %w", err)
	}
	fmt.Printf("Profile %q locked.\n", profile)
	return nil
}

func unlockProfile(dir, profile string) error {
	if !isProfileLocked(dir, profile) {
		return fmt.Errorf("profile %q is not locked", profile)
	}
	if err := os.Remove(lockFilePath(dir, profile)); err != nil {
		return fmt.Errorf("could not remove lock file: %w", err)
	}
	fmt.Printf("Profile %q unlocked.\n", profile)
	return nil
}

func listLockedProfiles(dir string) error {
	lockDir := filepath.Join(dir, ".envoy", "locks")
	entries, err := os.ReadDir(lockDir)
	if os.IsNotExist(err) {
		fmt.Println("No locked profiles.")
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not read locks directory: %w", err)
	}
	count := 0
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".lock" {
			name := e.Name()[:len(e.Name())-5]
			fmt.Printf("  [locked] %s\n", name)
			count++
		}
	}
	if count == 0 {
		fmt.Println("No locked profiles.")
	}
	return nil
}
