package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	var overwrite bool
	var keys []string

	promoteCmd := &cobra.Command{
		Use:   "promote <source-profile> <target-profile>",
		Short: "Promote keys from one profile to another",
		Long:  `Copy specific keys (or all keys) from a source profile into a target profile, optionally overwriting existing values.`,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			src, dst := args[0], args[1]
			if err := promoteProfile(projectDir(), src, dst, keys, overwrite); err != nil {
				fatalf("%v", err)
			}
		},
	}

	promoteCmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "Overwrite existing keys in target profile")
	promoteCmd.Flags().StringSliceVarP(&keys, "keys", "k", nil, "Comma-separated list of keys to promote (default: all)")

	rootCmd.AddCommand(promoteCmd)
}

func promoteProfile(root, src, dst string, keys []string, overwrite bool) error {
	srcPath := filepath.Join(root, ".envoy", "profiles", src+".env")
	dstPath := filepath.Join(root, ".envoy", "profiles", dst+".env")

	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return fmt.Errorf("source profile %q does not exist", src)
	}
	if _, err := os.Stat(dstPath); os.IsNotExist(err) {
		return fmt.Errorf("target profile %q does not exist", dst)
	}

	srcMap, err := readEnvMap(srcPath)
	if err != nil {
		return fmt.Errorf("reading source profile: %w", err)
	}
	dstMap, err := readEnvMap(dstPath)
	if err != nil {
		return fmt.Errorf("reading target profile: %w", err)
	}

	promoted, skipped := applyPromotion(srcMap, dstMap, keys, overwrite)

	if err := writeEnvMap(dstPath, dstMap); err != nil {
		return fmt.Errorf("writing target profile: %w", err)
	}

	for _, k := range promoted {
		fmt.Printf("promoted: %s\n", k)
	}
	for _, k := range skipped {
		fmt.Printf("skipped (already exists): %s\n", k)
	}
	writeAuditEntry(root, fmt.Sprintf("promote %s -> %s", src, dst))
	return nil
}
