package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func newTestRootForSchema(t *testing.T, dir string) *cobra.Command {
	t.Helper()
	root := &cobra.Command{Use: "envoy"}
	schema := &cobra.Command{Use: "schema"}

	defineCmd := &cobra.Command{
		Use:  "define <key>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			required, _ := cmd.Flags().GetBool("required")
			desc, _ := cmd.Flags().GetString("desc")
			def, _ := cmd.Flags().GetString("default")
			return defineSchemaKey(dir, args[0], required, desc, def)
		},
	}
	defineCmd.Flags().Bool("required", false, "")
	defineCmd.Flags().String("desc", "", "")
	defineCmd.Flags().String("default", "", "")

	validateCmd := &cobra.Command{
		Use:  "validate <profile>",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			errs := validateAgainstSchema(dir, args[0])
			for _, e := range errs {
				cmd.Println(e)
			}
			return nil
		},
	}

	schema.AddCommand(defineCmd, validateCmd)
	root.AddCommand(schema)
	return root
}

func TestSchemaDefineAndValidateViaCommand(t *testing.T) {
	dir := t.TempDir()
	_ = os.MkdirAll(filepath.Join(dir, ".envoy"), 0755)
	p := filepath.Join(dir, ".envoy", "prod.env")
	_ = os.WriteFile(p, []byte("DB_URL=postgres://localhost/db\n"), 0644)

	root := newTestRootForSchema(t, dir)
	root.SetArgs([]string{"schema", "define", "DB_URL", "--required"})
	if err := root.Execute(); err != nil {
		t.Fatalf("define failed: %v", err)
	}

	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"schema", "validate", "prod"})
	if err := root.Execute(); err != nil {
		t.Fatalf("validate failed: %v", err)
	}
	if strings.Contains(buf.String(), "missing") {
		t.Errorf("unexpected validation error: %s", buf.String())
	}
}

func TestSchemaValidateMissingKeyViaCommand(t *testing.T) {
	dir := t.TempDir()
	_ = os.MkdirAll(filepath.Join(dir, ".envoy"), 0755)
	p := filepath.Join(dir, ".envoy", "dev.env")
	_ = os.WriteFile(p, []byte("OTHER=value\n"), 0644)

	root := newTestRootForSchema(t, dir)
	root.SetArgs([]string{"schema", "define", "REQUIRED_KEY", "--required"})
	_ = root.Execute()

	var buf bytes.Buffer
	root.SetOut(&buf)
	root.SetArgs([]string{"schema", "validate", "dev"})
	_ = root.Execute()
	if !strings.Contains(buf.String(), "REQUIRED_KEY") {
		t.Errorf("expected missing key in output, got: %s", buf.String())
	}
}
