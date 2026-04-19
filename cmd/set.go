package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	var setCmd = &cobra.Command{
		Use:   "set [profile] [KEY=VALUE...]",
		Short: "Set one or more key-value pairs in a profile",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			profile := args[0]
			pairs := args[1:]
			if err := setProfileKeys(projectDir, profile, pairs); err != nil {
				fatalf("set: %v", err)
			}
		},
	}
	rootCmd.AddCommand(setCmd)
}

func setProfileKeys(root, profile string, pairs []string) error {
	path := profilePath(root, profile)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("profile %q does not exist", profile)
	}

	updates := make(map[string]string)
	for _, pair := range pairs {
		idx := strings.IndexByte(pair, '=')
		if idx < 1 {
			return fmt.Errorf("invalid key=value pair: %q", pair)
		}
		updates[pair[:idx]] = pair[idx+1:]
	}

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	f.Close()

	seen := make(map[string]bool)
	for i, line := range lines {
		if strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		key := strings.SplitN(line, "=", 2)[0]
		if val, ok := updates[key]; ok {
			lines[i] = key + "=" + val
			seen[key] = true
		}
	}
	for k, v := range updates {
		if !seen[k] {
			lines = append(lines, k+"="+v)
		}
	}

	out := strings.Join(lines, "\n")
	if len(lines) > 0 {
		out += "\n"
	}
	return os.WriteFile(path, []byte(out), 0644)
}
