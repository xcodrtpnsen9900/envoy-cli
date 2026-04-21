package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

func init() {
	watchCmd := &cobra.Command{
		Use:   "watch [profile]",
		Short: "Watch a profile for changes and display diffs",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			interval, _ := cmd.Flags().GetInt("interval")
			watchProfile(args[0], interval)
		},
	}
	watchCmd.Flags().IntP("interval", "i", 2, "Poll interval in seconds")
	rootCmd.AddCommand(watchCmd)
}

func watchProfile(name string, intervalSecs int) {
	path := profilePath(name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fatalf("profile %q does not exist", name)
	}

	fmt.Printf("Watching profile %q (every %ds). Press Ctrl+C to stop.\n", name, intervalSecs)

	lastContent, err := os.ReadFile(path)
	if err != nil {
		fatalf("failed to read profile: %v", err)
	}
	lastMod := fileModTime(path)

	ticker := time.NewTicker(time.Duration(intervalSecs) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		mod := fileModTime(path)
		if mod.Equal(lastMod) {
			continue
		}
		newContent, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading profile: %v\n", err)
			continue
		}
		fmt.Printf("\n[%s] Change detected in %q:\n", time.Now().Format("15:04:05"), filepath.Base(path))
		printWatchDiff(string(lastContent), string(newContent))
		lastContent = newContent
		lastMod = mod
	}
}

func fileModTime(path string) time.Time {
	info, err := os.Stat(path)
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}

func printWatchDiff(oldContent, newContent string) {
	oldLines := splitLines(oldContent)
	newLines := splitLines(newContent)

	oldSet := make(map[string]string)
	for _, line := range oldLines {
		if k, v, ok := parseEnvLine(line); ok {
			oldSet[k] = v
		}
	}
	newSet := make(map[string]string)
	for _, line := range newLines {
		if k, v, ok := parseEnvLine(line); ok {
			newSet[k] = v
		}
	}

	changes := 0
	for k, nv := range newSet {
		if ov, exists := oldSet[k]; !exists {
			fmt.Printf("  + %s=%s\n", k, nv)
			changes++
		} else if ov != nv {
			fmt.Printf("  ~ %s: %s -> %s\n", k, ov, nv)
			changes++
		}
	}
	for k, ov := range oldSet {
		if _, exists := newSet[k]; !exists {
			fmt.Printf("  - %s=%s\n", k, ov)
			changes++
		}
	}
	if changes == 0 {
		fmt.Println("  (no key-level changes detected)")
	}
}

func parseEnvLine(line string) (string, string, bool) {
	if len(line) == 0 || line[0] == '#' {
		return "", "", false
	}
	for i, ch := range line {
		if ch == '=' {
			return line[:i], line[i+1:], true
		}
	}
	return "", "", false
}
