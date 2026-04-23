package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	var overwrite bool

	cloneGroupCmd := &cobra.Command{
		Use:   "clone-group [group] [suffix]",
		Short: "Clone all profiles in a group with a new suffix",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			group := args[0]
			suffix := args[1]
			cloned, err := cloneGroup(projectDir(), group, suffix, overwrite)
			if err != nil {
				fatalf("clone-group: %v", err)
			}
			for _, name := range cloned {
				fmt.Printf("cloned → %s\n", name)
			}
			if len(cloned) == 0 {
				fmt.Println("no profiles cloned")
			}
		},
	}

	cloneGroupCmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "overwrite existing profiles")
	rootCmd.AddCommand(cloneGroupCmd)
}

// cloneGroup duplicates every profile in the named group, appending suffix to
// each profile name (e.g. "dev" + "-backup" → "dev-backup").
func cloneGroup(root, group, suffix string, overwrite bool) ([]string, error) {
	groups, err := loadGroups(root)
	if err != nil {
		return nil, fmt.Errorf("loading groups: %w", err)
	}

	members, ok := groups[group]
	if !ok {
		return nil, fmt.Errorf("group %q not found", group)
	}
	if len(members) == 0 {
		return nil, nil
	}

	envoyDir := filepath.Join(root, ".envoy")
	var cloned []string

	for _, profile := range members {
		src := filepath.Join(envoyDir, profile+".env")
		dst := filepath.Join(envoyDir, profile+suffix+".env")

		if _, err := os.Stat(src); os.IsNotExist(err) {
			return nil, fmt.Errorf("source profile %q does not exist", profile)
		}

		if !overwrite {
			if _, err := os.Stat(dst); err == nil {
				return nil, fmt.Errorf("profile %q already exists (use --overwrite to replace)", profile+suffix)
			}
		}

		if err := copyFile(src, dst); err != nil {
			return nil, fmt.Errorf("copying %q: %w", profile, err)
		}
		cloned = append(cloned, profile+suffix)
	}

	return cloned, nil
}
