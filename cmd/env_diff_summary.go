package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	var jsonOutput bool

	diffSummaryCmd := &cobra.Command{
		Use:   "diff-summary <profile1> <profile2>",
		Short: "Show a summary of differences between two profiles",
		Long: `Compare two .env profiles and display a concise summary:
- Keys only in profile1
- Keys only in profile2
- Keys present in both but with different values
- Total counts per category`,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			diffSummaryProfiles(args[0], args[1], jsonOutput)
		},
	}

	diffSummaryCmd.Flags().BoolVar(&jsonOutput, "json", false, "Output summary as JSON")
	rootCmd.AddCommand(diffSummaryCmd)
}

// diffSummary holds the categorised results of comparing two profiles.
type diffSummary struct {
	OnlyInA    []string          // keys present only in profile A
	OnlyInB    []string          // keys present only in profile B
	Changed    map[string][2]string // keys in both but with differing values [A, B]
	Identical  int               // count of keys with identical values
}

// buildDiffSummary compares two env maps and returns a diffSummary.
func buildDiffSummary(a, b map[string]string) diffSummary {
	summary := diffSummary{
		Changed: make(map[string][2]string),
	}

	for k, va := range a {
		if vb, ok := b[k]; ok {
			if va == vb {
				summary.Identical++
			} else {
				summary.Changed[k] = [2]string{va, vb}
			}
		} else {
			summary.OnlyInA = append(summary.OnlyInA, k)
		}
	}

	for k := range b {
		if _, ok := a[k]; !ok {
			summary.OnlyInB = append(summary.OnlyInB, k)
		}
	}

	sort.Strings(summary.OnlyInA)
	sort.Strings(summary.OnlyInB)

	return summary
}

// diffSummaryProfiles loads two profiles, computes the diff summary, and prints it.
func diffSummaryProfiles(nameA, nameB string, asJSON bool) {
	pathA := profilePath(nameA)
	pathB := profilePath(nameB)

	mapA, err := readEnvMap(pathA)
	if err != nil {
		if os.IsNotExist(err) {
			fatalf("profile %q not found", nameA)
		}
		fatalf("reading profile %q: %v", nameA, err)
	}

	mapB, err := readEnvMap(pathB)
	if err != nil {
		if os.IsNotExist(err) {
			fatalf("profile %q not found", nameB)
		}
		fatalf("reading profile %q: %v", nameB, err)
	}

	s := buildDiffSummary(mapA, mapB)

	if asJSON {
		printDiffSummaryJSON(nameA, nameB, s)
		return
	}
	printDiffSummaryText(nameA, nameB, s)
}

func printDiffSummaryText(nameA, nameB string, s diffSummary) {
	fmt.Printf("Diff summary: %s ↔ %s\n", nameA, nameB)
	fmt.Printf("  Identical keys : %d\n", s.Identical)
	fmt.Printf("  Changed keys   : %d\n", len(s.Changed))
	fmt.Printf("  Only in %-8s: %d\n", nameA, len(s.OnlyInA))
	fmt.Printf("  Only in %-8s: %d\n", nameB, len(s.OnlyInB))

	if len(s.OnlyInA) > 0 {
		fmt.Printf("\nOnly in %s:\n", nameA)
		for _, k := range s.OnlyInA {
			fmt.Printf("  - %s\n", k)
		}
	}

	if len(s.OnlyInB) > 0 {
		fmt.Printf("\nOnly in %s:\n", nameB)
		for _, k := range s.OnlyInB {
			fmt.Printf("  + %s\n", k)
		}
	}

	if len(s.Changed) > 0 {
		// Sort changed keys for deterministic output.
		keys := make([]string, 0, len(s.Changed))
		for k := range s.Changed {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		fmt.Println("\nChanged values:")
		for _, k := range keys {
			pair := s.Changed[k]
			fmt.Printf("  ~ %s\n", k)
			fmt.Printf("      %s: %s\n", nameA, pair[0])
			fmt.Printf("      %s: %s\n", nameB, pair[1])
		}
	}
}

func printDiffSummaryJSON(nameA, nameB string, s diffSummary) {
	// Minimal hand-rolled JSON to avoid importing encoding/json and keep deps light.
	changedKeys := make([]string, 0, len(s.Changed))
	for k := range s.Changed {
		changedKeys = append(changedKeys, k)
	}
	sort.Strings(changedKeys)

	changedParts := make([]string, 0, len(changedKeys))
	for _, k := range changedKeys {
		pair := s.Changed[k]
		changedParts = append(changedParts,
			fmt.Sprintf(`%q:{%q:%q,%q:%q}`, k, nameA, pair[0], nameB, pair[1]))
	}

	fmt.Printf(`{"profile_a":%q,"profile_b":%q,"identical":%d,"only_in_a":[%s],"only_in_b":[%s],"changed":{%s}}\n`,
		nameA, nameB, s.Identical,
		quoteJoin(s.OnlyInA),
		quoteJoin(s.OnlyInB),
		strings.Join(changedParts, ","),
	)
}

// quoteJoin returns a comma-separated, JSON-quoted list of strings.
func quoteJoin(ss []string) string {
	quoted := make([]string, len(ss))
	for i, s := range ss {
		quoted[i] = fmt.Sprintf("%q", s)
	}
	return strings.Join(quoted, ",")
}
