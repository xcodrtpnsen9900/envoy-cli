package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

type ProfileGroup struct {
	Name     string   `json:"name"`
	Profiles []string `json:"profiles"`
}

func groupFilePath(dir string) string {
	return filepath.Join(dir, ".envoy", "groups.json")
}

func init() {
	groupCmd := &cobra.Command{
		Use:   "group",
		Short: "Manage profile groups",
	}

	addGroupCmd := &cobra.Command{
		Use:   "add <group> <profile...>",
		Short: "Add profiles to a group",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			groupName := args[0]
			profiles := args[1:]
			if err := addProfilesToGroup(projectDir, groupName, profiles); err != nil {
				fatalf("group add: %v", err)
			}
			fmt.Printf("Added %d profile(s) to group %q\n", len(profiles), groupName)
		},
	}

	listGroupCmd := &cobra.Command{
		Use:   "list",
		Short: "List all groups",
		Run: func(cmd *cobra.Command, args []string) {
			groups, err := loadGroups(projectDir)
			if err != nil {
				fatalf("group list: %v", err)
			}
			if len(groups) == 0 {
				fmt.Println("No groups defined.")
				return
			}
			names := make([]string, 0, len(groups))
			for k := range groups {
				names = append(names, k)
			}
			sort.Strings(names)
			for _, name := range names {
				fmt.Printf("%s: %s\n", name, strings.Join(groups[name].Profiles, ", "))
			}
		},
	}

	removeGroupCmd := &cobra.Command{
		Use:   "remove <group>",
		Short: "Remove a group",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := deleteGroup(projectDir, args[0]); err != nil {
				fatalf("group remove: %v", err)
			}
			fmt.Printf("Group %q removed\n", args[0])
		},
	}

	groupCmd.AddCommand(addGroupCmd, listGroupCmd, removeGroupCmd)
	rootCmd.AddCommand(groupCmd)
}

func loadGroups(dir string) (map[string]*ProfileGroup, error) {
	path := groupFilePath(dir)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return map[string]*ProfileGroup{}, nil
	}
	if err != nil {
		return nil, err
	}
	var groups map[string]*ProfileGroup
	if err := json.Unmarshal(data, &groups); err != nil {
		return nil, fmt.Errorf("parse groups: %w", err)
	}
	return groups, nil
}

func saveGroups(dir string, groups map[string]*ProfileGroup) error {
	path := groupFilePath(dir)
	data, err := json.MarshalIndent(groups, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func addProfilesToGroup(dir, groupName string, profiles []string) error {
	groups, err := loadGroups(dir)
	if err != nil {
		return err
	}
	g, ok := groups[groupName]
	if !ok {
		g = &ProfileGroup{Name: groupName}
	}
	existing := make(map[string]bool)
	for _, p := range g.Profiles {
		existing[p] = true
	}
	for _, p := range profiles {
		if !existing[p] {
			g.Profiles = append(g.Profiles, p)
			existing[p] = true
		}
	}
	sort.Strings(g.Profiles)
	groups[groupName] = g
	return saveGroups(dir, groups)
}

func deleteGroup(dir, groupName string) error {
	groups, err := loadGroups(dir)
	if err != nil {
		return err
	}
	if _, ok := groups[groupName]; !ok {
		return fmt.Errorf("group %q not found", groupName)
	}
	delete(groups, groupName)
	return saveGroups(dir, groups)
}
