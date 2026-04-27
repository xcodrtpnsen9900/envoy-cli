package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rollbackCmd := &cobra.Command{
		Use:   "rollback [profile]",
		Short: "Roll back a profile to its previous snapshot",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			profile := args[0]
			steps, _ := cmd.Flags().GetInt("steps")
			if steps < 1 {
				steps = 1
			}
			if err := rollbackProfile(projectDir(), profile, steps); err != nil {
				fatalf("rollback failed: %v", err)
			}
		},
	}
	rollbackCmd.Flags().IntP("steps", "n", 1, "Number of snapshots to roll back")
	rootCmd.AddCommand(rollbackCmd)

	rollbackListCmd := &cobra.Command{
		Use:   "rollback-list [profile]",
		Short: "List available rollback points for a profile",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			profile := args[0]
			if err := listRollbackPoints(projectDir(), profile); err != nil {
				fatalf("rollback-list failed: %v", err)
			}
		},
	}
	rootCmd.AddCommand(rollbackListCmd)
}

func rollbackProfile(root, profile string, steps int) error {
	snaps, err := sortedSnapshotsForProfile(root, profile)
	if err != nil {
		return err
	}
	if len(snaps) == 0 {
		return fmt.Errorf("no snapshots found for profile %q", profile)
	}
	if steps > len(snaps) {
		steps = len(snaps)
	}
	target := snaps[len(snaps)-steps]
	dst := profilePath(root, profile)
	if err := copyFileContents(target, dst); err != nil {
		return fmt.Errorf("could not restore snapshot: %w", err)
	}
	fmt.Printf("Rolled back profile %q to snapshot: %s\n", profile, filepath.Base(target))
	return nil
}

func listRollbackPoints(root, profile string) error {
	snaps, err := sortedSnapshotsForProfile(root, profile)
	if err != nil {
		return err
	}
	if len(snaps) == 0 {
		fmt.Printf("No rollback points found for profile %q\n", profile)
		return nil
	}
	fmt.Printf("Rollback points for %q (oldest → newest):\n", profile)
	for i, s := range snaps {
		fmt.Printf("  [%d] %s\n", i+1, filepath.Base(s))
	}
	return nil
}

func sortedSnapshotsForProfile(root, profile string) ([]string, error) {
	dir := snapshotsDir(root)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	prefix := profile + "_"
	var matches []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasPrefix(e.Name(), prefix) {
			matches = append(matches, filepath.Join(dir, e.Name()))
		}
	}
	sort.Strings(matches)
	return matches, nil
}

func copyFileContents(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}
