package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const historyFileName = "switch_history.log"
const maxHistoryEntries = 50

func historyFilePath(dir string) string {
	return filepath.Join(dir, ".envoy", historyFileName)
}

func recordSwitchHistory(dir, profile string) error {
	path := historyFilePath(dir)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	existing, _ := readHistoryEntries(path)
	entry := fmt.Sprintf("%s %s", time.Now().Format(time.RFC3339), profile)
	entries := append([]string{entry}, existing...)
	if len(entries) > maxHistoryEntries {
		entries = entries[:maxHistoryEntries]
	}
	return os.WriteFile(path, []byte(strings.Join(entries, "\n")+"\n"), 0644)
}

func readHistoryEntries(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var entries []string
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if line != "" {
			entries = append(entries, line)
		}
	}
	return entries, nil
}

func showHistory(dir string, limit int) ([]string, error) {
	entries, err := readHistoryEntries(historyFilePath(dir))
	if err != nil {
		return nil, err
	}
	if limit > 0 && limit < len(entries) {
		entries = entries[:limit]
	}
	return entries, nil
}

func init() {
	var limit int
	var clear bool

	cmd := &cobra.Command{
		Use:   "history",
		Short: "Show profile switch history",
		Run: func(cmd *cobra.Command, args []string) {
			dir := projectDir()
			if clear {
				path := historyFilePath(dir)
				if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
					fatalf("failed to clear history: %v", err)
				}
				fmt.Println("Switch history cleared.")
				return
			}
			entries, err := showHistory(dir, limit)
			if err != nil {
				fatalf("failed to read history: %v", err)
			}
			if len(entries) == 0 {
				fmt.Println("No switch history found.")
				return
			}
			for _, e := range entries {
				fmt.Println(e)
			}
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "n", 10, "Number of recent entries to show")
	cmd.Flags().BoolVar(&clear, "clear", false, "Clear the switch history")
	rootCmd.AddCommand(cmd)
}
