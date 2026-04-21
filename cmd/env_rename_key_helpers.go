package cmd

import (
	"bufio"
	"os"
	"strings"
)

// keyExistsInProfile returns true if the given key is present in the profile.
func keyExistsInProfile(profile, key string) (bool, error) {
	path := profilePath(profile)
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if strings.TrimSpace(parts[0]) == key {
			return true, nil
		}
	}
	return false, scanner.Err()
}

// listKeysInProfile returns all non-comment, non-empty keys from a profile.
func listKeysInProfile(profile string) ([]string, error) {
	path := profilePath(profile)
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var keys []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if k := strings.TrimSpace(parts[0]); k != "" {
			keys = append(keys, k)
		}
	}
	return keys, scanner.Err()
}
