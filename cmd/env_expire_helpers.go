package cmd

import (
	"fmt"
	"time"
)

// isExpired returns true if the profile has an expiry entry that is in the past.
func isExpired(dir, profile string) (bool, error) {
	entries, err := loadExpiry(dir)
	if err != nil {
		return false, err
	}
	entry, ok := entries[profile]
	if !ok {
		return false, nil
	}
	return time.Now().After(entry.ExpiresAt), nil
}

// clearExpiry removes the expiry entry for a profile.
func clearExpiry(dir, profile string) error {
	entries, err := loadExpiry(dir)
	if err != nil {
		return err
	}
	if _, ok := entries[profile]; !ok {
		return fmt.Errorf("no expiry set for profile %q", profile)
	}
	delete(entries, profile)
	return saveExpiry(dir, entries)
}

// expiredProfiles returns the names of all profiles whose expiry has passed.
func expiredProfiles(dir string) ([]string, error) {
	entries, err := loadExpiry(dir)
	if err != nil {
		return nil, err
	}
	var expired []string
	for name, entry := range entries {
		if time.Now().After(entry.ExpiresAt) {
			expired = append(expired, name)
		}
	}
	return expired, nil
}

// expiryCount returns the number of profiles with an expiry set.
func expiryCount(dir string) (int, error) {
	entries, err := loadExpiry(dir)
	if err != nil {
		return 0, err
	}
	return len(entries), nil
}

// timeUntilExpiry returns the duration remaining until a profile expires.
// Returns an error if no expiry is set for the profile.
func timeUntilExpiry(dir, profile string) (time.Duration, error) {
	entries, err := loadExpiry(dir)
	if err != nil {
		return 0, err
	}
	entry, ok := entries[profile]
	if !ok {
		return 0, fmt.Errorf("no expiry set for profile %q", profile)
	}
	return time.Until(entry.ExpiresAt), nil
}
