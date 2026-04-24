package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

type ExpiryEntry struct {
	Profile   string    `json:"profile"`
	ExpiresAt time.Time `json:"expires_at"`
}

func expiryFilePath(dir string) string {
	return filepath.Join(dir, ".envoy", "expiry.json")
}

func init() {
	setExpiryCmd := &cobra.Command{
		Use:   "expire <profile> <duration>",
		Short: "Set an expiry duration on a profile (e.g. 24h, 7d)",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if err := setProfileExpiry(projectDir, args[0], args[1]); err != nil {
				fatalf("%v", err)
			}
		},
	}

	checkExpiryCmd := &cobra.Command{
		Use:   "expire-check",
		Short: "Check for expired or soon-to-expire profiles",
		Run: func(cmd *cobra.Command, args []string) {
			if err := checkProfileExpiry(projectDir); err != nil {
				fatalf("%v", err)
			}
		},
	}

	rootCmd.AddCommand(setExpiryCmd)
	rootCmd.AddCommand(checkExpiryCmd)
}

func setProfileExpiry(dir, profile, duration string) error {
	p := profilePath(dir, profile)
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return fmt.Errorf("profile %q does not exist", profile)
	}
	d, err := time.ParseDuration(duration)
	if err != nil {
		return fmt.Errorf("invalid duration %q: %w", duration, err)
	}
	entries, _ := loadExpiry(dir)
	entries[profile] = ExpiryEntry{Profile: profile, ExpiresAt: time.Now().Add(d)}
	if err := saveExpiry(dir, entries); err != nil {
		return err
	}
	fmt.Printf("Profile %q will expire at %s\n", profile, entries[profile].ExpiresAt.Format(time.RFC3339))
	return nil
}

func checkProfileExpiry(dir string) error {
	entries, err := loadExpiry(dir)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		fmt.Println("No expiry entries found.")
		return nil
	}
	now := time.Now()
	for name, entry := range entries {
		remaining := time.Until(entry.ExpiresAt)
		if remaining <= 0 {
			fmt.Printf("[EXPIRED] %s (expired %s ago)\n", name, (-remaining).Round(time.Second))
		} else if remaining < 24*time.Hour {
			fmt.Printf("[WARNING] %s expires in %s\n", name, remaining.Round(time.Second))
		} else {
			fmt.Printf("[OK]      %s expires at %s\n", name, entry.ExpiresAt.Format(time.RFC3339))
		}
		_ = now
	}
	return nil
}

func loadExpiry(dir string) (map[string]ExpiryEntry, error) {
	path := expiryFilePath(dir)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return map[string]ExpiryEntry{}, nil
	}
	if err != nil {
		return nil, err
	}
	var entries map[string]ExpiryEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func saveExpiry(dir string, entries map[string]ExpiryEntry) error {
	path := expiryFilePath(dir)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
