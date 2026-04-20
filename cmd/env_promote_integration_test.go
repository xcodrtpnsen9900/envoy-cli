package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func newTestRootForPromote(t *testing.T) (root string, execute func(args ...string) (string, error)) {
	t.Helper()
	root = t.TempDir()
	dir := filepath.Join(root, ".envoy", "profiles")
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}

	execute = func(args ...string) (string, error) {
		buf := &bytes.Buffer{}
		cmd := &cobra.Command{Use: "envoy"}
		cmd.PersistentFlags().String("project", root, "")

		promoteCmd := &cobra.Command{
			Use:  "promote <src> <dst>",
			Args: cobra.ExactArgs(2),
			RunE: func(c *cobra.Command, a []string) error {
				overwrite, _ := c.Flags().GetBool("overwrite")
				keys, _ := c.Flags().GetStringSlice("keys")
				return promoteProfile(root, a[0], a[1], keys, overwrite)
			},
		}
		promoteCmd.Flags().Bool("overwrite", false, "")
		promoteCmd.Flags().StringSlice("keys", nil, "")
		cmd.AddCommand(promoteCmd)
		cmd.SetOut(buf)
		cmd.SetErr(buf)
		cmd.SetArgs(args)
		err := cmd.Execute()
		return buf.String(), err
	}
	return
}

func TestPromoteViaCommand(t *testing.T) {
	root, execute := newTestRootForPromote(t)
	writePromoteProfile(t, root, "dev", "APP_ENV=development\nDEBUG=true\n")
	writePromoteProfile(t, root, "prod", "APP_ENV=production\n")

	_, err := execute("promote", "dev", "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m, _ := readEnvMap(filepath.Join(root, ".envoy", "profiles", "prod.env"))
	if m["DEBUG"] != "true" {
		t.Errorf("expected DEBUG=true in prod, got %q", m["DEBUG"])
	}
	if m["APP_ENV"] != "production" {
		t.Errorf("APP_ENV should not be overwritten, got %q", m["APP_ENV"])
	}
}

func TestPromoteOutputMessages(t *testing.T) {
	root, _ := newTestRootForPromote(t)
	writePromoteProfile(t, root, "dev", "NEW_KEY=hello\n")
	writePromoteProfile(t, root, "prod", "")

	// Call directly to capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	promoteProfile(root, "dev", "prod", nil, false)
	w.Close()
	os.Stdout = old

	buf := &bytes.Buffer{}
	buf.ReadFrom(r)
	out := buf.String()

	if !strings.Contains(out, "promoted: NEW_KEY") {
		t.Errorf("expected promotion message, got: %s", out)
	}
}
