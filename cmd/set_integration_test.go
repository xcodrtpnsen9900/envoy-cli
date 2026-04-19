package cmd

import (
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func newTestRootForSet(dir string) *cobra.Command {
	root := &cobra.Command{Use: "envoy"}
	set := &cobra.Command{
		Use:  "set [profile] [KEY=VALUE...]",
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return setProfileKeys(dir, args[0], args[1:])
		},
	}
	root.AddCommand(set)
	return root
}

func TestSetViaCommand(t *testing.T) {
	dir := setupTempDir(t)
	addProfile(dir, "staging", "HOST=localhost\n")

	root := newTestRootForSet(dir)
	root.SetArgs([]string{"set", "staging", "HOST=prod.example.com", "PORT=443"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}

	content := readProfile(t, dir, "staging")
	if !strings.Contains(content, "HOST=prod.example.com") {
		t.Errorf("HOST not updated: %s", content)
	}
	if !strings.Contains(content, "PORT=443") {
		t.Errorf("PORT not added: %s", content)
	}
}

func TestSetPreservesComments(t *testing.T) {
	dir := setupTempDir(t)
	addProfile(dir, "prod", "# production config\nDB=old\n")

	root := newTestRootForSet(dir)
	root.SetArgs([]string{"set", "prod", "DB=new"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}

	content := readProfile(t, dir, "prod")
	if !strings.Contains(content, "# production config") {
		t.Errorf("comment should be preserved: %s", content)
	}
	if !strings.Contains(content, "DB=new") {
		t.Errorf("DB not updated: %s", content)
	}
}
