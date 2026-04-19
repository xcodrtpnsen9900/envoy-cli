package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Show audit log of profile changes",
	Run: func(cmd *cobra.Command, args []string) {
		entries, err := readAuditLog(projectDir)
		if err != nil {
			fatalf("could not read audit log: %v", err)
		}
		if len(entries) == 0 {
			fmt.Println("No audit log entries found.")
			return
		}
		for _, e := range entries {
			fmt.Println(e)
		}
	},
}

func init() {
	rootCmd.AddCommand(auditCmd)
}

func auditLogPath(dir string) string {
	return filepath.Join(dir, ".envoy", "audit.log")
}

func writeAuditEntry(dir, action, profile string) {
	path := auditLogPath(dir)
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	timestamp := time.Now().Format(time.RFC3339)
	fmt.Fprintf(f, "[%s] %s profile=%s\n", timestamp, action, profile)
}

func readAuditLog(dir string) ([]string, error) {
	path := auditLogPath(dir)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	lines := splitLines(string(data))
	var entries []string
	for _, l := range lines {
		if l != "" {
			entries = append(entries, l)
		}
	}
	return entries, nil
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
