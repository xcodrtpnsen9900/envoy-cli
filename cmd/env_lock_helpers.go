package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// isProfileLocked returns true if the given profile has an active lock file.
func isProfileLocked(dir, profile string) bool {
	_, err := os.Stat(lockFilePath(dir, profile))
	return err == nil
}

// assertNotLocked returns an error if the profile is locked.
func assertNotLocked(dir, profile string) error {
	if isProfileLocked(dir, profile) {
		return fmt.Errorf("profile %q is locked; unlock it before making changes", profile)
	}
	return nil
}

// lockTimestamp reads the locked_at timestamp from a lock file.
func lockTimestamp(dir, profile string) (string, error) {
	f, err := os.Open(lockFilePath(dir, profile))
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "locked_at=") {
			return strings.TrimPrefix(line, "locked_at="), nil
		}
	}
	return "", fmt.Errorf("timestamp not found in lock file")
}

// lockedProfileNames returns all profile names that are currently locked.
func lockedProfileNames(dir string) ([]string, error) {
	lockDir := filepath.Join(dir, ".envoy", "locks")
	entries, err := os.ReadDir(lockDir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".lock" {
			names = append(names, e.Name()[:len(e.Name())-5])
		}
	}
	return names, nil
}
