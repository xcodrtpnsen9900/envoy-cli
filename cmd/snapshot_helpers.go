package cmd

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// latestSnapshot returns the most recent snapshot filename (without .env)
// for the given profile, or empty string if none exist.
func latestSnapshot(profile string) (string, error) {
	dir := snapshotsDir(profile)
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".env") {
			names = append(names, e.Name())
		}
	}
	if len(names) == 0 {
		return "", nil
	}
	sort.Strings(names)
	return strings.TrimSuffix(names[len(names)-1], ".env"), nil
}

// deleteSnapshot removes a specific snapshot file.
func deleteSnapshot(profile, timestamp string) error {
	path := filepath.Join(snapshotsDir(profile), timestamp+".env")
	return os.Remove(path)
}

// snapshotCount returns the number of snapshots for a profile.
func snapshotCount(profile string) (int, error) {
	dir := snapshotsDir(profile)
	entries, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	count := 0
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".env") {
			count++
		}
	}
	return count, nil
}
