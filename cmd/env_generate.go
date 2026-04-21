package cmd

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const (
	charsetAlpha   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	charsetNumeric = "0123456789"
	charsetSpecial = "!@#$%^&*()-_=+[]{}"
	charsetAll     = charsetAlpha + charsetNumeric + charsetSpecial
)

func init() {
	var length int
	var noSpecial bool
	var noNumeric bool
	var prefix string
	var keys []string
	var overwrite bool

	generateCmd := &cobra.Command{
		Use:   "generate <profile>",
		Short: "Generate random values for keys in a profile",
		Long: `Generate cryptographically random values for specified keys in a profile.
If no keys are provided, generates values for all keys that are empty or missing.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return generateProfileKeys(args[0], keys, length, noSpecial, noNumeric, prefix, overwrite)
		},
	}

	generateCmd.Flags().IntVarP(&length, "length", "l", 32, "Length of generated value")
	generateCmd.Flags().BoolVar(&noSpecial, "no-special", false, "Exclude special characters")
	generateCmd.Flags().BoolVar(&noNumeric, "no-numeric", false, "Exclude numeric characters")
	generateCmd.Flags().StringVar(&prefix, "prefix", "", "Prefix to prepend to generated values")
	generateCmd.Flags().StringSliceVarP(&keys, "keys", "k", nil, "Specific keys to generate values for (comma-separated)")
	generateCmd.Flags().BoolVar(&overwrite, "overwrite", false, "Overwrite existing non-empty values")

	rootCmd.AddCommand(generateCmd)
}

// generateProfileKeys generates random values for keys in the given profile.
func generateProfileKeys(profile string, keys []string, length int, noSpecial, noNumeric bool, prefix string, overwrite bool) error {
	path := filepath.Join(projectDir(), ".envoy", "profiles", profile+".env")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("profile %q does not exist", profile)
	}

	lines, err := readLines(path)
	if err != nil {
		return fmt.Errorf("failed to read profile: %w", err)
	}

	// Build charset based on flags
	charset := charsetAlpha
	if !noNumeric {
		charset += charsetNumeric
	}
	if !noSpecial {
		charset += charsetSpecial
	}

	// Build a set of target keys for quick lookup
	targetKeys := make(map[string]bool, len(keys))
	for _, k := range keys {
		targetKeys[strings.TrimSpace(k)] = true
	}

	updated := make([]string, 0, len(lines))
	generated := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Preserve comments and blank lines
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			updated = append(updated, line)
			continue
		}

		eqIdx := strings.Index(line, "=")
		if eqIdx < 0 {
			updated = append(updated, line)
			continue
		}

		key := strings.TrimSpace(line[:eqIdx])
		value := line[eqIdx+1:]

		// Determine whether this key should be generated
		shouldGenerate := false
		if len(targetKeys) > 0 {
			shouldGenerate = targetKeys[key]
		} else {
			// No specific keys provided: generate for empty values
			shouldGenerate = strings.TrimSpace(value) == ""
		}

		if shouldGenerate && (overwrite || strings.TrimSpace(value) == "") {
			rand, err := randomString(charset, length)
			if err != nil {
				return fmt.Errorf("failed to generate value for key %q: %w", key, err)
			}
			updated = append(updated, key+"="+prefix+rand)
			generated++
			fmt.Printf("generated: %s\n", key)
		} else {
			updated = append(updated, line)
		}
	}

	if generated == 0 {
		fmt.Println("no keys were updated")
		return nil
	}

	content := strings.Join(updated, "\n")
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("failed to write profile: %w", err)
	}

	fmt.Printf("\n%d key(s) updated in profile %q\n", generated, profile)
	return nil
}

// randomString generates a cryptographically secure random string of the given
// length using characters from the provided charset.
func randomString(charset string, length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be greater than 0")
	}
	if len(charset) == 0 {
		return "", fmt.Errorf("charset must not be empty")
	}

	result := make([]byte, length)
	max := big.NewInt(int64(len(charset)))

	for i := range result {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		result[i] = charset[n.Int64()]
	}

	return string(result), nil
}
