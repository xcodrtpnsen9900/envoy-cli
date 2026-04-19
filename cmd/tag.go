package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	tagCmd := &cobra.Command{
		Use:   "tag [profile] [tags...]",
		Short: "Tag a profile with labels",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			profile := args[0]
			tags := args[1:]
			if err := tagProfile(projectDir, profile, tags); err != nil {
				fatalf("%v", err)
			}
		},
	}

	tagListCmd := &cobra.Command{
		Use:   "tag-list [profile]",
		Short: "List tags for a profile",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			tags, err := getTags(projectDir, args[0])
			if err != nil {
				fatalf("%v", err)
			}
			if len(tags) == 0 {
				fmt.Println("no tags")
				return
			}
			fmt.Println(strings.Join(tags, ", "))
		},
	}

	rootCmd.AddCommand(tagCmd)
	rootCmd.AddCommand(tagListCmd)
}

func tagFilePath(dir, profile string) string {
	return filepath.Join(dir, ".envoy", profile+".tags")
}

func tagProfile(dir, profile string, tags []string) error {
	envoyDir := filepath.Join(dir, ".envoy")
	if _, err := os.Stat(envoyDir); os.IsNotExist(err) {
		return fmt.Errorf("envoy not initialized in %s", dir)
	}
	pPath := filepath.Join(envoyDir, profile+".env")
	if _, err := os.Stat(pPath); os.IsNotExist(err) {
		return fmt.Errorf("profile %q does not exist", profile)
	}
	existing, _ := getTags(dir, profile)
	merged := mergeTags(existing, tags)
	content := strings.Join(merged, "\n") + "\n"
	return os.WriteFile(tagFilePath(dir, profile), []byte(content), 0644)
}

func getTags(dir, profile string) ([]string, error) {
	path := tagFilePath(dir, profile)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var tags []string
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if t := strings.TrimSpace(line); t != "" {
			tags = append(tags, t)
		}
	}
	return tags, nil
}

func mergeTags(existing, newTags []string) []string {
	seen := map[string]bool{}
	var result []string
	for _, t := range existing {
		if !seen[t] {
			seen[t] = true
			result = append(result, t)
		}
	}
	for _, t := range newTags {
		if !seen[t] {
			seen[t] = true
			result = append(result, t)
		}
	}
	return result
}
