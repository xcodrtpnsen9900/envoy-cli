package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	var shell string
	exportCmd := &cobra.Command{
		Use:   "export [profile]",
		Short: "Print export statements for a profile's env vars",
		Long:  "Print shell export statements for the given profile (or active profile if omitted). Pipe to eval to load into current shell.",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var name string
			if len(args) == 1 {
				name = args[0]
			} else {
				active, err := activeProfile()
				if err != nil || active == "" {
					fatalf("no active profile; specify a profile name or run 'envoy switch'")
				}
				name = active
			}
			output, err := exportProfile(name, shell)
			if err != nil {
				fatalf("%v", err)
			}
			fmt.Print(output)
		},
	}
	exportCmd.Flags().StringVarP(&shell, "shell", "s", "bash", "Shell format: bash or fish")
	rootCmd.AddCommand(exportCmd)
}

func exportProfile(name, shell string) (string, error) {
	path := profilePath(name)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("profile %q not found", name)
		}
		return "", err
	}

	var sb strings.Builder
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if !strings.Contains(line, "=") {
			continue
		}
		switch strings.ToLower(shell) {
		case "fish":
			parts := strings.SplitN(line, "=", 2)
			sb.WriteString(fmt.Sprintf("set -x %s %s;\n", parts[0], parts[1]))
		default:
			sb.WriteString(fmt.Sprintf("export %s\n", line))
		}
	}
	return sb.String(), nil
}

func profilePath(name string) string {
	return filepath.Join(projectDir(), ".envoy", "profiles", name+".env")
}
