package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func init() {
	var overwrite bool

	importCmd := &cobra.Command{
		Use:   "import [profile] [file]",
		Short: "Import keys from an external .env file into a profile",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			profile := args[0]
			source := args[1]
			importProfile(profile, source, overwrite)
		},
	}

	importCmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "Overwrite existing keys in the profile")
	rootCmd.AddCommand(importCmd)
}

func importProfile(profile, source string, overwrite bool) {
	if _, err := os.Stat(source); os.IsNotExist(err) {
		fatalf("source file not found: %s", source)
	}

	dest := profilePath(profile)
	if _, err := os.Stat(dest); os.IsNotExist(err) {
		fatalf("profile %q does not exist", profile)
	}

	srcMap, err := readImportEnvMap(source)
	if err != nil {
		fatalf("failed to read source file: %v", err)
	}

	destMap, err := readImportEnvMap(dest)
	if err != nil {
		fatalf("failed to read profile: %v", err)
	}

	added, skipped := 0, 0
	for k, v := range srcMap {
		if _, exists := destMap[k]; exists && !overwrite {
			skipped++
			continue
		}
		destMap[k] = v
		added++
	}

	if err := writeImportEnvMap(dest, destMap); err != nil {
		fatalf("failed to write profile: %v", err)
	}

	fmt.Printf("Imported %d key(s) into profile %q", added, profile)
	if skipped > 0 {
		fmt.Printf(" (%d skipped, use --overwrite to replace)", skipped)
	}
	fmt.Println()
	writeAuditEntry(filepath.Join(projectDir(), ".envoy"), "import", fmt.Sprintf("%s from %s", profile, source))
}

func readImportEnvMap(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			m[strings.TrimSpace(parts[0])] = parts[1]
		}
	}
	return m, scanner.Err()
}

func writeImportEnvMap(path string, m map[string]string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for k, v := range m {
		fmt.Fprintf(w, "%s=%s\n", k, v)
	}
	return w.Flush()
}
