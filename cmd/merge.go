package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	mergeCmd := &cobra.Command{
		Use:   "merge <base> <overlay> <output>",
		Short: "Merge two profiles into a new profile",
		Long:  "Merge two env profiles, with overlay values taking precedence over base values.",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			if err := mergeProfiles(args[0], args[1], args[2]); err != nil {
				fatalf("merge failed: %v", err)
			}
		},
	}
	rootCmd.AddCommand(mergeCmd)
}

func mergeProfiles(base, overlay, output string) error {
	dir := projectDir()

	baseVars, err := readEnvMap(profilePath(dir, base))
	if err != nil {
		return fmt.Errorf("cannot read base profile %q: %w", base, err)
	}

	overlayVars, err := readEnvMap(profilePath(dir, overlay))
	if err != nil {
		return fmt.Errorf("cannot read overlay profile %q: %w", overlay, err)
	}

	outPath := profilePath(dir, output)
	if _, err := os.Stat(outPath); err == nil {
		return fmt.Errorf("profile %q already exists", output)
	}

	for k, v := range overlayVars {
		baseVars[k] = v
	}

	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("cannot create output profile: %w", err)
	}
	defer f.Close()

	for k, v := range baseVars {
		fmt.Fprintf(f, "%s=%s\n", k, v)
	}

	fmt.Printf("Merged %q + %q -> %q\n", base, overlay, output)
	return nil
}

func readEnvMap(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result, scanner.Err()
}
