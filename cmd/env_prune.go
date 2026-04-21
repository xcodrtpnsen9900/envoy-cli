package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	var dryRun bool

	pruneCmd := &cobra.Command{
		Use:   "prune <profile>",
		Short: "Remove duplicate and empty keys from a profile",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			pruneProfile(args[0], dryRun)
		},
	}

	pruneCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "Show what would be removed without modifying the file")
	rootCmd.AddCommand(pruneCmd)
}

func pruneProfile(name string, dryRun bool) {
	path := profilePath(name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fatalf("profile %q does not exist", name)
	}

	f, err := os.Open(path)
	if err != nil {
		fatalf("failed to open profile: %v", err)
	}
	defer f.Close()

	var kept []string
	seen := make(map[string]bool)
	removed := 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Always keep comments and blank lines
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			kept = append(kept, line)
			continue
		}

		eqIdx := strings.Index(trimmed, "=")
		if eqIdx < 0 {
			// Not a valid key=value line; keep as-is
			kept = append(kept, line)
			continue
		}

		key := strings.TrimSpace(trimmed[:eqIdx])
		if key == "" {
			fmt.Printf("  [prune] removing empty key line: %q\n", line)
			removed++
			continue
		}

		if seen[key] {
			fmt.Printf("  [prune] removing duplicate key %q\n", key)
			removed++
			continue
		}

		seen[key] = true
		kept = append(kept, line)
	}

	if err := scanner.Err(); err != nil {
		fatalf("error reading profile: %v", err)
	}

	if removed == 0 {
		fmt.Printf("profile %q is already clean (no duplicates or empty keys)\n", name)
		return
	}

	if dryRun {
		fmt.Printf("dry-run: would remove %d line(s) from profile %q\n", removed, name)
		return
	}

	out := strings.Join(kept, "\n")
	if !strings.HasSuffix(out, "\n") {
		out += "\n"
	}

	if err := os.WriteFile(path, []byte(out), 0644); err != nil {
		fatalf("failed to write pruned profile: %v", err)
	}

	fmt.Printf("pruned %d line(s) from profile %q\n", removed, name)
}
