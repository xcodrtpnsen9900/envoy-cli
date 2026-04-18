package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	var showCmd = &cobra.Command{
		Use:   "show [profile]",
		Short: "Display the contents of a profile",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var name string
			if len(args) == 0 {
				active, err := activeProfile()
				if err != nil || active == "" {
					fatalf("no active profile and no profile specified")
				}
				name = active
			} else {
				name = args[0]
			}
			if err := showProfile(name); err != nil {
				fatalf("%v", err)
			}
		},
	}
	rootCmd.AddCommand(showCmd)
}

func showProfile(name string) error {
	path := filepath.Join(projectDir(), ".envoy", "profiles", name+".env")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("profile %q does not exist", name)
		}
		return err
	}
	fmt.Printf("# Profile: %s\n", name)
	fmt.Print(string(data))
	return nil
}
