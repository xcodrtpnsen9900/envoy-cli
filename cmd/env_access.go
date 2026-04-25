package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type AccessEntry struct {
	Profile   string    `json:"profile"`
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
}

func accessLogPath(root string) string {
	return filepath.Join(root, ".envoy", "access_log.json")
}

func recordAccess(root, profile, action string) error {
	entries, _ := readAccessLog(root)
	entries = append(entries, AccessEntry{
		Profile:   profile,
		Action:    action,
		Timestamp: time.Now().UTC(),
	})
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(accessLogPath(root), data, 0644)
}

func readAccessLog(root string) ([]AccessEntry, error) {
	data, err := os.ReadFile(accessLogPath(root))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var entries []AccessEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func showAccessLog(root, filterProfile string, limit int) {
	entries, err := readAccessLog(root)
	if err != nil {
		fatalf("error reading access log: %v", err)
	}
	if len(entries) == 0 {
		fmt.Println("No access log entries.")
		return
	}
	filtered := entries
	if filterProfile != "" {
		filtered = nil
		for _, e := range entries {
			if strings.EqualFold(e.Profile, filterProfile) {
				filtered = append(filtered, e)
			}
		}
	}
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Timestamp.After(filtered[j].Timestamp)
	})
	if limit > 0 && len(filtered) > limit {
		filtered = filtered[:limit]
	}
	for _, e := range filtered {
		fmt.Printf("[%s] %-10s %s\n", e.Timestamp.Format("2006-01-02 15:04:05"), e.Action, e.Profile)
	}
}

func init() {
	var profile string
	var limit int

	accessCmd := &cobra.Command{
		Use:   "access",
		Short: "Show profile access log",
		Run: func(cmd *cobra.Command, args []string) {
			showAccessLog(projectDir, profile, limit)
		},
	}
	accessCmd.Flags().StringVarP(&profile, "profile", "p", "", "Filter by profile name")
	accessCmd.Flags().IntVarP(&limit, "limit", "n", 20, "Max entries to show")
	rootCmd.AddCommand(accessCmd)
}
