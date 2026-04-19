package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	snapshotCmd := &cobra.Command{
		Use:   "snapshot [profile]",
		Short: "Save a timestamped snapshot of a profile",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := snapshotProfile(args[0]); err != nil {
				fatalf("%v", err)
			}
		},
	}

	listSnapshotsCmd := &cobra.Command{
		Use:   "snapshots [profile]",
		Short: "List all snapshots for a profile",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := listSnapshots(args[0]); err != nil {
				fatalf("%v", err)
			}
		},
	}

	restoreSnapshotCmd := &cobra.Command{
		Use:   "restore [profile] [timestamp]",
		Short: "Restore a profile from a snapshot",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if err := restoreSnapshot(args[0], args[1]); err != nil {
				fatalf("%v", err)
			}
		},
	}

	rootCmd.AddCommand(snapshotCmd)
	rootCmd.AddCommand(listSnapshotsCmd)
	rootCmd.AddCommand(restoreSnapshotCmd)
}

func snapshotsDir(profile string) string {
	return filepath.Join(projectDir(), ".envoy", "snapshots", profile)
}

func snapshotProfile(profile string) error {
	src := profilePath(profile)
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return fmt.Errorf("profile %q does not exist", profile)
	}

	dir := snapshotsDir(profile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102T150405")
	dst := filepath.Join(dir, timestamp+".env")
	if err := copyFile(src, dst); err != nil {
		return err
	}
	fmt.Printf("Snapshot saved: %s\n", timestamp)
	return nil
}

func listSnapshots(profile string) error {
	dir := snapshotsDir(profile)
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		fmt.Printf("No snapshots for profile %q\n", profile)
		return nil
	}
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		fmt.Printf("No snapshots for profile %q\n", profile)
		return nil
	}
	for _, e := range entries {
		fmt.Println(e.Name())
	}
	return nil
}

func restoreSnapshot(profile, timestamp string) error {
	src := filepath.Join(snapshotsDir(profile), timestamp+".env")
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return fmt.Errorf("snapshot %q not found for profile %q", timestamp, profile)
	}
	dst := profilePath(profile)
	if err := copyFile(src, dst); err != nil {
		return err
	}
	fmt.Printf("Profile %q restored from snapshot %s\n", profile, timestamp)
	return nil
}
