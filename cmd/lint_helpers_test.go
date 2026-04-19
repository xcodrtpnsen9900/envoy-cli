package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

// writeProfile writes a named profile for lint tests reusing the profiles dir layout.
func writeProfile(t *testing.T, dir, name, content string) {
	t.Helper()
	profilesDir := filepath.Join(dir, ".envoy", "profiles")
	if err := os.MkdirAll(profilesDir, 0755); err != nil {
		t.Fatalf("mkdir profiles: %v", err)
	}
	path := filepath.Join(profilesDir, name+".env")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write profile: %v", err)
	}
}
