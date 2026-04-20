package cmd

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	var inPlace bool

	sortCmd := &cobra.Command{
		Use:   "sort <profile>",
		Short: "Sort keys in a profile alphabetically",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			path := profilePath(name)
			lines, comments, err := sortProfile(path)
			if err != nil {
				fatalf("sort: %v", err)
			}
			if inPlace {
				if err := writeSortedProfile(path, comments, lines); err != nil {
					fatalf("sort: %v", err)
				}
				fmt.Printf("Profile '%s' sorted in place.\n", name)
			} else {
				for _, c := range comments {
					fmt.Println(c)
				}
				for _, l := range lines {
					fmt.Println(l)
				}
			}
		},
	}

	sortCmd.Flags().BoolVarP(&inPlace, "in-place", "i", false, "Overwrite the profile file with sorted output")
	rootCmd.AddCommand(sortCmd)
}

func sortProfile(path string) (kvLines []string, commentLines []string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot open profile: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			commentLines = append(commentLines, line)
		} else {
			kvLines = append(kvLines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}
	sort.Slice(kvLines, func(i, j int) bool {
		return strings.ToLower(kvLines[i]) < strings.ToLower(kvLines[j])
	})
	return kvLines, commentLines, nil
}

func writeSortedProfile(path string, comments, kvLines []string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, c := range comments {
		fmt.Fprintln(w, c)
	}
	for _, l := range kvLines {
		fmt.Fprintln(w, l)
	}
	return w.Flush()
}
