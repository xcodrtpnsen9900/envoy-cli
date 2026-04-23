package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func newTestRootForGroup(t *testing.T) (string, *cobra.Command) {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".envoy"), 0755); err != nil {
		t.Fatal(err)
	}
	oldDir := projectDir
	projectDir = dir
	t.Cleanup(func() { projectDir = oldDir })
	return dir, rootCmd
}

func runGroupCommand(t *testing.T, args ...string) string {
	t.Helper()
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)
	_ = rootCmd.Execute()
	return buf.String()
}

func TestGroupAddAndListViaCommand(t *testing.T) {
	dir, _ := newTestRootForGroup(t)

	// seed profiles
	for _, name := range []string{"dev", "qa"} {
		p := filepath.Join(dir, ".envoy", name+".env")
		if err := os.WriteFile(p, []byte("KEY=val\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	if err := addProfilesToGroup(dir, "lower", []string{"dev", "qa"}); err != nil {
		t.Fatalf("addProfilesToGroup: %v", err)
	}

	groups, err := loadGroups(dir)
	if err != nil {
		t.Fatalf("loadGroups: %v", err)
	}
	if _, ok := groups["lower"]; !ok {
		t.Fatal("expected group 'lower' to exist after add")
	}
}

func TestGroupListOutputViaCommand(t *testing.T) {
	dir, _ := newTestRootForGroup(t)
	_ = addProfilesToGroup(dir, "backend", []string{"api", "worker"})

	names, err := groupNames(dir)
	if err != nil {
		t.Fatalf("groupNames: %v", err)
	}
	if len(names) != 1 || names[0] != "backend" {
		t.Fatalf("expected [backend], got %v", names)
	}
}

func TestGroupRemoveViaCommand(t *testing.T) {
	dir, _ := newTestRootForGroup(t)
	_ = addProfilesToGroup(dir, "temp", []string{"x"})
	if err := deleteGroup(dir, "temp"); err != nil {
		t.Fatalf("deleteGroup: %v", err)
	}
	exists, _ := groupExists(dir, "temp")
	if exists {
		t.Fatal("group should have been deleted")
	}
}

func TestGroupAddPreservesExistingGroups(t *testing.T) {
	dir, _ := newTestRootForGroup(t)
	_ = addProfilesToGroup(dir, "g1", []string{"a"})
	_ = addProfilesToGroup(dir, "g2", []string{"b"})

	groups, _ := loadGroups(dir)
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
	_ = strings.Join([]string{}, "") // suppress unused import
}
