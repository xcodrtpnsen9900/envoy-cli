package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	backupCmd := &cobra.Command{
		Use:   "backup [profile]",
		Short: "Create a timestamped backup of a profile",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dest, _ := cmd.Flags().GetString("output")
			if err := backupProfile(args[0], dest); err != nil {
				fatalf("backup failed: %v", err)
			}
		},
	}
	backupCmd.Flags().StringP("output", "o", "", "Output directory for backup file (default: .envoy/backups)")

	restoreBackupCmd := &cobra.Command{
		Use:   "restore-backup [profile] [backup-file]",
		Short: "Restore a profile from a backup file",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if err := restoreBackup(args[0], args[1]); err != nil {
				fatalf("restore failed: %v", err)
			}
		},
	}

	listBackupsCmd := &cobra.Command{
		Use:   "list-backups [profile]",
		Short: "List all backups for a profile",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := listBackups(args[0]); err != nil {
				fatalf("list-backups failed: %v", err)
			}
		},
	}

	rootCmd.AddCommand(backupCmd, restoreBackupCmd, listBackupsCmd)
}

func backupsDir() string {
	return filepath.Join(projectDir(), ".envoy", "backups")
}

func backupProfile(profile, outputDir string) error {
	src := profilePath(profile)
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return fmt.Errorf("profile %q does not exist", profile)
	}

	dir := outputDir
	if dir == "" {
		dir = filepath.Join(backupsDir(), profile)
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("could not create backup directory: %w", err)
	}

	timestamp := time.Now().Format("20060102T150405")
	destName := fmt.Sprintf("%s.%s.env", profile, timestamp)
	dest := filepath.Join(dir, destName)

	if err := copyFileBackup(src, dest); err != nil {
		return fmt.Errorf("could not write backup: %w", err)
	}

	fmt.Printf("Backup created: %s\n", dest)
	writeAuditEntry("backup", profile)
	return nil
}

func restoreBackup(profile, backupFile string) error {
	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		return fmt.Errorf("backup file %q not found", backupFile)
	}
	dest := profilePath(profile)
	if err := copyFileBackup(backupFile, dest); err != nil {
		return fmt.Errorf("could not restore backup: %w", err)
	}
	fmt.Printf("Profile %q restored from %s\n", profile, backupFile)
	writeAuditEntry("restore-backup", profile)
	return nil
}

func listBackups(profile string) error {
	dir := filepath.Join(backupsDir(), profile)
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		fmt.Printf("No backups found for profile %q\n", profile)
		return nil
	}
	if err != nil {
		return fmt.Errorf("could not read backup directory: %w", err)
	}
	if len(entries) == 0 {
		fmt.Printf("No backups found for profile %q\n", profile)
		return nil
	}
	fmt.Printf("Backups for profile %q:\n", profile)
	for _, e := range entries {
		if !e.IsDir() {
			fmt.Printf("  %s\n", e.Name())
		}
	}
	return nil
}

func copyFileBackup(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
