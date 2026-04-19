package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "compare [profile1] [profile2]",
		Short: "Compare keys between two profiles showing missing/extra keys",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			if err := compareProfiles(args[0], args[1]); err != nil {
				fatalf("%v", err)
			}
		},
	})
}

func compareProfiles(profileA, profileB string) error {
	pathA := profilePath(profileA)
	pathB := profilePath(profileB)

	mapA, err := readEnvMap(pathA)
	if err != nil {
		return fmt.Errorf("cannot read profile %q: %w", profileA, err)
	}
	mapB, err := readEnvMap(pathB)
	if err != nil {
		return fmt.Errorf("cannot read profile %q: %w", profileB, err)
	}

	onlyInA := keysOnlyIn(mapA, mapB)
	onlyInB := keysOnlyIn(mapB, mapA)
	inBoth := sharedKeys(mapA, mapB)

	fmt.Printf("Keys only in %s (%d):\n", profileA, len(onlyInA))
	for _, k := range onlyInA {
		fmt.Printf("  - %s\n", k)
	}

	fmt.Printf("Keys only in %s (%d):\n", profileB, len(onlyInB))
	for _, k := range onlyInB {
		fmt.Printf("  + %s\n", k)
	}

	fmt.Printf("Shared keys: %d\n", len(inBoth))
	_ = strings.Join(inBoth, ",")
	_ = os.Stdout
	return nil
}

func keysOnlyIn(a, b map[string]string) []string {
	var result []string
	for k := range a {
		if _, ok := b[k]; !ok {
			result = append(result, k)
		}
	}
	sort.Strings(result)
	return result
}

func sharedKeys(a, b map[string]string) []string {
	var result []string
	for k := range a {
		if _, ok := b[k]; ok {
			result = append(result, k)
		}
	}
	sort.Strings(result)
	return result
}
