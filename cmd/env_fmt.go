package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	fmtCmd := &cobra.Command{
		Use:   "fmt <profile>",
		Short: "Format a .env profile: trim whitespace, normalize spacing around '='",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			if err := fmtProfile(args[0], dryRun); err != nil {
				fatalf("%v", err)
			}
		},
	}
	fmtCmd.Flags().Bool("dry-run", false, "Print formatted output without writing to file")
	rootCmd.AddCommand(fmtCmd)
}

func fmtProfile(name string, dryRun bool) error {
	path := profilePath(name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("profile %q does not exist", name)
	}

	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("could not open profile: %w", err)
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, formatEnvLine(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading profile: %w", err)
	}

	formatted := strings.Join(lines, "\n") + "\n"

	if dryRun {
		fmt.Print(formatted)
		return nil
	}

	if err := os.WriteFile(path, []byte(formatted), 0644); err != nil {
		return fmt.Errorf("could not write formatted profile: %w", err)
	}
	fmt.Printf("formatted profile %q\n", name)
	return nil
}

func formatEnvLine(line string) string {
	trimmed := strings.TrimSpace(line)
	// Preserve blank lines and comments as-is
	if trimmed == "" || strings.HasPrefix(trimmed, "#") {
		return trimmed
	}
	idx := strings.Index(trimmed, "=")
	if idx < 0 {
		return trimmed
	}
	key := strings.TrimSpace(trimmed[:idx])
	value := strings.TrimSpace(trimmed[idx+1:])
	return key + "=" + value
}
