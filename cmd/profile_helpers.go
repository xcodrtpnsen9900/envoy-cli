package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

// addProfile creates a new profile file. Returns an error for use in tests.
func addProfile(args []string, _ interface{}) error {
	name := args[0]
	path := profilePath(name)

	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("profile %q already exists", name)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create profile: %v", err)
	}
	f.Close()
	return nil
}
