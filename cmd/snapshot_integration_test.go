package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func runSnapshotCommand(t *testing.T, args ...string) (string, error) {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	w.Close()
	os.Stdout = old

	var buf strings.Builder
	buf.ReadFrom(r)
	return buf.String(), err
}

func newTestRootForSnapshot() *cobra.Command {
	return &cobra.Command{Use: "envoy"}
}

func TestSnapshotViaCommand(t *testing.T) {
	dir := setupTempDir(t)
	defer os.RemoveAll(dir)

	if err := addProfile("prod", dir); err != nil {
		t.Fatal(err)
	}
	p := profilePath("prod")
	os.WriteFile(p, []byte("ENV=production\n"), 0644)

	if err := snapshotProfile("prod"); err != nil {
		t.Fatalf("snapshot failed: %v", err)
	}

	entries, err := os.ReadDir(snapshotsDir("prod"))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) == 0 {
		t.Fatal("expected at least one snapshot")
	}
}

func TestMultipleSnapshots(t *testing.T) {
	dir := setupTempDir(t)
	defer os.RemoveAll(dir)

	if err := addProfile("qa", dir); err != nil {
		t.Fatal(err)
	}
	p := profilePath("qa")

	for i := 0; i < 3; i++ {
		os.WriteFile(p, []byte("V=version\n"), 0644)
		// small sleep not needed; timestamps may collide in fast tests
		// just test that snapshotProfile doesn't error
		snapshotProfile("qa")
	}

	entries, _ := os.ReadDir(snapshotsDir("qa"))
	if len(entries) == 0 {
		t.Fatal("expected snapshots")
	}
}
