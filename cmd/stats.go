package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "stats [profile]",
		Short: "Show statistics for a profile (key count, comment lines, empty lines)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := profileStats(args[0]); err != nil {
				fatalf("%v", err)
			}
		},
	})
}

type envStats struct {
	Keys     int
	Comments int
	Empty    int
	Total    int
}

func profileStats(name string) error {
	p := filepath.Join(projectDir, ".envoy", name+".env")
	data, err := os.ReadFile(p)
	if err != nil {
		return fmt.Errorf("profile %q not found", name)
	}

	lines := strings.Split(string(data), "\n")
	var s envStats
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			s.Empty++
		} else if strings.HasPrefix(trimmed, "#") {
			s.Comments++
		} else if strings.Contains(trimmed, "=") {
			s.Keys++
		}
		s.Total++
	}

	fmt.Printf("Profile : %s\n", name)
	fmt.Printf("Keys    : %d\n", s.Keys)
	fmt.Printf("Comments: %d\n", s.Comments)
	fmt.Printf("Empty   : %d\n", s.Empty)
	fmt.Printf("Total   : %d\n", s.Total)
	return nil
}
