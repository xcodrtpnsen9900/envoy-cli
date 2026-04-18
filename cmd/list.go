package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("verbose", "v", false, "Show file sizes and modification times")
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available profiles",
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		profiles, err := listProfilesDetailed(projectDir, verbose)
		if err != nil {
			fatalf("Error listing profiles: %v", err)
		}
		if len(profiles) == 0 {
			fmt.Println("No profiles found. Use 'envoy add' to create one.")
			return
		}
		active := activeProfile(projectDir)
		for _, p := range profiles {
			marker := "  "
			if p.name == active {
				marker = "* "
			}
			if verbose {
				fmt.Printf("%s%-20s %6d bytes  %s\n", marker, p.name, p.size, p.modTime)
			} else {
				fmt.Printf("%s%s\n", marker, p.name)
			}
		}
	},
}

type profileInfo struct {
	name    string
	size    int64
	modTime string
}

func listProfilesDetailed(dir string, verbose bool) ([]profileInfo, error) {
	envoyDir := filepath.Join(dir, ".envoy", "profiles")
	entries, err := os.ReadDir(envoyDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var result []profileInfo
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		info := profileInfo{name: e.Name()}
		if verbose {
			fi, err := e.Info()
			if err == nil {
				info.size = fi.Size()
				info.modTime = fi.ModTime().Format("2006-01-02 15:04")
			}
		}
		result = append(result, info)
	}
	return result, nil
}
