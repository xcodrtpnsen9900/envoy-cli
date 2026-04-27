package cmd

import (
	"os"
	"path/filepath"
	"strings"
)

// rollbackCount returns the number of snapshots available for a profile.
func rollbackCount(root, profile string) int {
	snaps, err := sortedSnapshotsForProfile(root, profile)
	if err != nil {
		return 0
	}
	return len(snaps)
}

// latestRollbackSnapshot returns the path to the most recent snapshot for a
// profile, or an empty string if none exist.
func latestRollbackSnapshot(root, profile string) string {
	snaps, err := sortedSnapshotsForProfile(root, profile)
	if err != nil || len(snaps) == 0 {
		return ""
	}
	return snaps[len(snaps)-1]
}

// nthRollbackSnapshot returns the nth most-recent snapshot (1 = latest).
// Returns empty string if out of range.
func nthRollbackSnapshot(root, profile string, n int) string {
	snaps, err := sortedSnapshotsForProfile(root, profile)
	if err != nil || len(snaps) == 0 {
		return ""
	}
	idx := len(snaps) - n
	if idx < 0 {
		idx = 0
	}
	return snaps[idx]
}

// profilesWithRollbackPoints returns all profile names that have at least one
// snapshot in the snapshots directory.
func profilesWithRollbackPoints(root string) ([]string, error) {
	dir := snapshotsDir(root)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	seen := map[string]struct{}{}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if idx := strings.LastIndex(name, "_"); idx > 0 {
			profile := name[:idx]
			seen[profile] = struct{}{}
		}
	}
	var result []string
	for p := range seen {
		result = append(result, p)
	}
	return result, nil
}

// rollbackSnapshotNames returns just the base filenames of snapshots for a
// profile, sorted oldest-first.
func rollbackSnapshotNames(root, profile string) []string {
	snaps, err := sortedSnapshotsForProfile(root, profile)
	if err != nil {
		return nil
	}
	names := make([]string, len(snaps))
	for i, s := range snaps {
		names[i] = filepath.Base(s)
	}
	return names
}
