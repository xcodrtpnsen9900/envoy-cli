package cmd

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// backupCount returns the number of backup files for a given profile.
func backupCount(profile string) int {
	dir := filepath.Join(backupsDir(), profile)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	count := 0
	for _, e := range entries {
		if !e.IsDir() {
			count++
		}
	}
	return count
}

// latestBackup returns the path to the most recently created backup file
// for a profile, or an empty string if none exist.
func latestBackup(profile string) string {
	dir := filepath.Join(backupsDir(), profile)
	entries, err := os.ReadDir(dir)
	if err != nil || len(entries) == 0 {
		return ""
	}
	names := []string{}
	for _, e := range entries {
		if !e.IsDir() {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	return filepath.Join(dir, names[len(names)-1])
}

// purgeOldBackups removes all but the most recent `keep` backups for a profile.
func purgeOldBackups(profile string, keep int) error {
	dir := filepath.Join(backupsDir(), profile)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil // nothing to purge
	}
	names := []string{}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".env") {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	if len(names) <= keep {
		return nil
	}
	for _, name := range names[:len(names)-keep] {
		if err := os.Remove(filepath.Join(dir, name)); err != nil {
			return err
		}
	}
	return nil
}
