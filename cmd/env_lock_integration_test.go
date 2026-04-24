package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func newTestRootForLock(t *testing.T, dir string) *cobra.Command {
	t.Helper()
	root := &cobra.Command{Use: "envoy"}

	lock := &cobra.Command{
		Use:  "lock [profile]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return lockProfile(dir, args[0])
		},
	}
	unlock := &cobra.Command{
		Use:  "unlock [profile]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return unlockProfile(dir, args[0])
		},
	}
	listLocked := &cobra.Command{
		Use: "list-locked",
		RunE: func(cmd *cobra.Command, args []string) error {
			return listLockedProfiles(dir)
		},
	}
	root.AddCommand(lock, unlock, listLocked)
	return root
}

func TestLockViaCommand(t *testing.T) {
	dir := setupLockDir(t)
	root := newTestRootForLock(t, dir)
	root.SetArgs([]string{"lock", "dev"})
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !isProfileLocked(dir, "dev") {
		t.Error("expected dev to be locked after command")
	}
}

func TestUnlockViaCommand(t *testing.T) {
	dir := setupLockDir(t)
	_ = lockProfile(dir, "dev")
	root := newTestRootForLock(t, dir)
	root.SetArgs([]string{"unlock", "dev"})
	if err := root.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if isProfileLocked(dir, "dev") {
		t.Error("expected dev to be unlocked after command")
	}
}

func TestListLockedViaCommand(t *testing.T) {
	dir := setupLockDir(t)
	envoyDir := filepath.Join(dir, ".envoy")
	_ = os.WriteFile(filepath.Join(envoyDir, "prod.env"), []byte("K=v\n"), 0644)
	_ = lockProfile(dir, "dev")
	_ = lockProfile(dir, "prod")

	var buf bytes.Buffer
	root := newTestRootForLock(t, dir)
	root.SetOut(&buf)
	root.SetArgs([]string{"list-locked"})
	_ = root.Execute()

	names, _ := lockedProfileNames(dir)
	if len(names) != 2 {
		t.Errorf("expected 2 locked profiles, got %d", len(names))
	}
	for _, n := range names {
		if !strings.Contains(n, "dev") && !strings.Contains(n, "prod") {
			t.Errorf("unexpected locked profile name: %s", n)
		}
	}
}
