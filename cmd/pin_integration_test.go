package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func newTestRootForPin(dir string) *cobra.Command {
	root := &cobra.Command{Use: "envoy"}
	root.AddCommand(&cobra.Command{
		Use:  "pin [profile]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return pinProfile(dir, args[0])
		},
	})
	root.AddCommand(&cobra.Command{
		Use:  "unpin [profile]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return unpinProfile(dir, args[0])
		},
	})
	return root
}

func TestPinViaCommand(t *testing.T) {
	dir := setupTempDir(t)
	initProject(dir)
	addProfile(dir, "staging", false)

	root := newTestRootForPin(dir)
	root.SetArgs([]string{"pin", "staging"})
	var buf bytes.Buffer
	root.SetOut(&buf)
	if err := root.Execute(); err != nil {
		t.Fatalf("pin command failed: %v", err)
	}
	if !isPinned(dir, "staging") {
		t.Error("expected staging to be pinned after command")
	}
}

func TestUnpinViaCommand(t *testing.T) {
	dir := setupTempDir(t)
	initProject(dir)
	addProfile(dir, "staging", false)
	_ = pinProfile(dir, "staging")

	root := newTestRootForPin(dir)
	root.SetArgs([]string{"unpin", "staging"})
	if err := root.Execute(); err != nil {
		t.Fatalf("unpin command failed: %v", err)
	}
	if isPinned(dir, "staging") {
		t.Error("expected staging to be unpinned after command")
	}
}

func TestPinNonExistentViaCommand(t *testing.T) {
	dir := setupTempDir(t)
	initProject(dir)

	root := newTestRootForPin(dir)
	root.SetArgs([]string{"pin", "ghost"})
	if err := root.Execute(); err == nil {
		t.Error("expected error pinning non-existent profile")
	}
}
