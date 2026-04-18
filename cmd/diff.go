package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "diff <profile1> <profile2>",
		Short: "Show differences between two env profiles",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			diffProfiles(args[0], args[1])
		},
	})
}

func parseEnvFile(path string) (map[string]string, []string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	vals := make(map[string]string)
	var keys []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		if _, exists := vals[key]; !exists {
			keys = append(keys, key)
		}
		vals[key] = val
	}
	return vals, keys, scanner.Err()
}

func diffProfiles(a, b string) {
	pathA := profilePath(a)
	pathB := profilePath(b)

	valsA, keysA, err := parseEnvFile(pathA)
	if err != nil {
		fatalf("could not read profile '%s': %v", a, err)
	}
	valsB, keysB, err := parseEnvFile(pathB)
	if err != nil {
		fatalf("could not read profile '%s': %v", b, err)
	}

	seen := make(map[string]bool)
	printed := false

	for _, k := range keysA {
		seen[k] = true
		valA := valsA[k]
		valB, inB := valsB[k]
		if !inB {
			fmt.Printf("- %s=%s\n", k, valA)
			printed = true
		} else if valA != valB {
			fmt.Printf("< %s=%s\n", k, valA)
			fmt.Printf("> %s=%s\n", k, valB)
			printed = true
		}
	}

	for _, k := range keysB {
		if !seen[k] {
			fmt.Printf("+ %s=%s\n", k, valsB[k])
			printed = true
		}
	}

	if !printed {
		fmt.Printf("Profiles '%s' and '%s' are identical.\n", a, b)
	}
}
