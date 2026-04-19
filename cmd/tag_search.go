package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	tagSearchCmd := &cobra.Command{
		Use:   "tag-search [tag]",
		Short: "Find all profiles with a given tag",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			profiles, err := profilesByTag(projectDir, args[0])
			if err != nil {
				fatalf("%v", err)
			}
			if len(profiles) == 0 {
				fmt.Printf("no profiles tagged %q\n", args[0])
				return
			}
			for _, p := range profiles {
				fmt.Println(p)
			}
		},
	}
	rootCmd.AddCommand(tagSearchCmd)
}

func profilesByTag(dir, tag string) ([]string, error) {
	envoyDir := filepath.Join(dir, ".envoy")
	entries, err := os.ReadDir(envoyDir)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("envoy not initialized in %s", dir)
	}
	if err != nil {
		return nil, err
	}
	var matched []string
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".tags") {
			continue
		}
		profile := strings.TrimSuffix(e.Name(), ".tags")
		tags, err := getTags(dir, profile)
		if err != nil {
			continue
		}
		for _, t := range tags {
			if strings.EqualFold(t, tag) {
				matched = append(matched, profile)
				break
			}
		}
	}
	return matched, nil
}
