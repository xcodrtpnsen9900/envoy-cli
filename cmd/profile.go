package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const profilesDir = ".envoy"
const activeLink = ".env"

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available env profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := projectDir(cmd)
		profiles, err := listProfiles(dir)
		if err != nil {
			return err
		}
		if len(profiles) == 0 {
			fmt.Println("No profiles found. Use 'envoy add <name>' to create one.")
			return nil
		}
		for _, p := range profiles {
			fmt.Println(" -", p)
		}
		return nil
	},
}

var useCmd = &cobra.Command{
	Use:   "use <profile>",
	Short: "Switch the active .env to the given profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := projectDir(cmd)
		return switchProfile(dir, args[0])
	},
}

var addCmd = &cobra.Command{
	Use:   "add <profile>",
	Short: "Add a new env profile (copies current .env if it exists)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := projectDir(cmd)
		return addProfile(dir, args[0])
	},
}

func profilePath(dir, name string) string {
	return filepath.Join(dir, profilesDir, name+".env")
}

func listProfiles(dir string) ([]string, error) {
	pDir := filepath.Join(dir, profilesDir)
	entries, err := os.ReadDir(pDir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".env" {
			names = append(names, e.Name()[:len(e.Name())-4])
		}
	}
	return names, nil
}

func switchProfile(dir, name string) error {
	src := profilePath(dir, name)
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return fmt.Errorf("profile %q not found", name)
	}
	dst := filepath.Join(dir, activeLink)
	_ = os.Remove(dst)
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.WriteFile(dst, data, 0644); err != nil {
		return err
	}
	fmt.Printf("Switched to profile %q\n", name)
	return nil
}

func addProfile(dir, name string) error {
	pDir := filepath.Join(dir, profilesDir)
	if err := os.MkdirAll(pDir, 0755); err != nil {
		return err
	}
	dst := profilePath(dir, name)
	if _, err := os.Stat(dst); err == nil {
		return fmt.Errorf("profile %q already exists", name)
	}
	var data []byte
	src := filepath.Join(dir, activeLink)
	if d, err := os.ReadFile(src); err == nil {
		data = d
	}
	if err := os.WriteFile(dst, data, 0644); err != nil {
		return err
	}
	fmt.Printf("Created profile %q\n", name)
	return nil
}
