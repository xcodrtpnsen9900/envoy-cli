package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func newTestRootForEnvCheck(dir string) *cobra.Command {
	projectDir = dir
	root := &cobra.Command{Use: "envoy"}
	root.AddCommand(&cobra.Command{
		Use:  "env-check <profile> <required-keys...>",
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			missing, err := checkRequiredKeys(args[0], args[1:])
			if err != nil {
				return err
			}
			if len(missing) > 0 {
				cmd.Printf("Missing: %s", strings.Join(missing, ", "))
				return nil
			}
			cmd.Printf("Profile '%s' contains all required keys.", args[0])
			return nil
		},
	})
	return root
}

func TestEnvCheckAllPresentViaCommand(t *testing.T) {
	dir := setupTempDir(t)
	envoyDir := filepath.Join(dir, ".envoy")
	os.MkdirAll(envoyDir, 0755)
	os.WriteFile(filepath.Join(envoyDir, "prod.env"), []byte("APP_KEY=abc\nDB=pg\n"), 0644)

	root := newTestRootForEnvCheck(dir)
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"env-check", "prod", "APP_KEY", "DB"})
	root.Execute()

	if !strings.Contains(buf.String(), "all required keys") {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestEnvCheckMissingViaCommand(t *testing.T) {
	dir := setupTempDir(t)
	envoyDir := filepath.Join(dir, ".envoy")
	os.MkdirAll(envoyDir, 0755)
	os.WriteFile(filepath.Join(envoyDir, "prod.env"), []byte("APP_KEY=abc\n"), 0644)

	root := newTestRootForEnvCheck(dir)
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"env-check", "prod", "APP_KEY", "MISSING_KEY"})
	root.Execute()

	if !strings.Contains(buf.String(), "Missing") {
		t.Errorf("expected missing output, got: %s", buf.String())
	}
}
