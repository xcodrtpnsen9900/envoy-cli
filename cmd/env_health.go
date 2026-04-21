package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	var jsonOutput bool

	healthCmd := &cobra.Command{
		Use:   "health [profile]",
		Short: "Check health of a profile's env file",
		Long:  "Reports issues such as empty values, duplicate keys, and missing keys in a profile.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			report, err := healthCheck(projectDir(), args[0])
			if err != nil {
				fatalf("%v", err)
			}
			if jsonOutput {
				printHealthJSON(report)
			} else {
				printHealthText(report)
			}
		},
	}

	healthCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output results as JSON")
	rootCmd.AddCommand(healthCmd)
}

type healthReport struct {
	Profile    string
	TotalKeys  int
	EmptyVals  []string
	Duplicates []string
	OK         bool
}

func healthCheck(root, profile string) (*healthReport, error) {
	path := profilePath(root, profile)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("profile %q not found", profile)
	}

	seen := map[string]int{}
	var emptyVals []string
	var duplicates []string
	totalKeys := 0

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		totalKeys++
		seen[key]++
		if val == "" {
			emptyVals = append(emptyVals, key)
		}
	}

	for k, count := range seen {
		if count > 1 {
			duplicates = append(duplicates, k)
		}
	}

	ok := len(emptyVals) == 0 && len(duplicates) == 0
	return &healthReport{
		Profile:    profile,
		TotalKeys:  totalKeys,
		EmptyVals:  emptyVals,
		Duplicates: duplicates,
		OK:         ok,
	}, nil
}

func printHealthText(r *healthReport) {
	fmt.Printf("Profile: %s\n", r.Profile)
	fmt.Printf("Total keys: %d\n", r.TotalKeys)
	if r.OK {
		fmt.Println("Status: ✓ healthy")
		return
	}
	fmt.Println("Status: ✗ issues found")
	if len(r.EmptyVals) > 0 {
		fmt.Printf("  Empty values (%d): %s\n", len(r.EmptyVals), strings.Join(r.EmptyVals, ", "))
	}
	if len(r.Duplicates) > 0 {
		fmt.Printf("  Duplicate keys (%d): %s\n", len(r.Duplicates), strings.Join(r.Duplicates, ", "))
	}
}

func printHealthJSON(r *healthReport) {
	empty := "[]"
	if len(r.EmptyVals) > 0 {
		empty = `["` + strings.Join(r.EmptyVals, `","`) + `"]`
	}
	dupes := "[]"
	if len(r.Duplicates) > 0 {
		dupes = `["` + strings.Join(r.Duplicates, `","`) + `"]`
	}
	fmt.Printf(`{"profile":%q,"total_keys":%d,"empty_values":%s,"duplicates":%s,"ok":%v}`+"\n",
		r.Profile, r.TotalKeys, empty, dupes, r.OK)
}
