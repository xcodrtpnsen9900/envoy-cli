package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"
)

func aliasFilePath(dir string) string {
	return filepath.Join(dir, ".envoy", "aliases.json")
}

func loadAliases(dir string) (map[string]string, error) {
	path := aliasFilePath(dir)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return map[string]string{}, nil
	}
	if err != nil {
		return nil, err
	}
	var aliases map[string]string
	if err := json.Unmarshal(data, &aliases); err != nil {
		return nil, err
	}
	return aliases, nil
}

func saveAliases(dir string, aliases map[string]string) error {
	path := aliasFilePath(dir)
	data, err := json.MarshalIndent(aliases, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func init() {
	aliasCmd := &cobra.Command{
		Use:   "alias",
		Short: "Manage profile aliases",
	}

	setAliasCmd := &cobra.Command{
		Use:   "set <alias> <profile>",
		Short: "Create an alias for a profile",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			alias, profile := args[0], args[1]
			dir := projectDir()
			if _, err := os.Stat(profilePath(dir, profile)); os.IsNotExist(err) {
				return fmt.Errorf("profile %q does not exist", profile)
			}
			aliases, err := loadAliases(dir)
			if err != nil {
				return err
			}
			aliases[alias] = profile
			if err := saveAliases(dir, aliases); err != nil {
				return err
			}
			fmt.Printf("Alias %q -> %q set.\n", alias, profile)
			return nil
		},
	}

	removeAliasCmd := &cobra.Command{
		Use:   "remove <alias>",
		Short: "Remove an alias",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			alias := args[0]
			dir := projectDir()
			aliases, err := loadAliases(dir)
			if err != nil {
				return err
			}
			if _, ok := aliases[alias]; !ok {
				return fmt.Errorf("alias %q not found", alias)
			}
			delete(aliases, alias)
			if err := saveAliases(dir, aliases); err != nil {
				return err
			}
			fmt.Printf("Alias %q removed.\n", alias)
			return nil
		},
	}

	listAliasCmd := &cobra.Command{
		Use:   "list",
		Short: "List all aliases",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := projectDir()
			aliases, err := loadAliases(dir)
			if err != nil {
				return err
			}
			if len(aliases) == 0 {
				fmt.Println("No aliases defined.")
				return nil
			}
			keys := make([]string, 0, len(aliases))
			for k := range aliases {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				fmt.Printf("%-20s -> %s\n", k, aliases[k])
			}
			return nil
		},
	}

	aliasCmd.AddCommand(setAliasCmd, removeAliasCmd, listAliasCmd)
	rootCmd.AddCommand(aliasCmd)
}
