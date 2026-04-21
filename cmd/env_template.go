package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var templateVarRe = regexp.MustCompile(`\{\{([A-Z0-9_]+)\}\}`)

func init() {
	var outputProfile string

	cmd := &cobra.Command{
		Use:   "template <template-file> <source-profile>",
		Short: "Render a template file using values from a profile",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			templateFile := args[0]
			sourceProfile := args[1]
			if err := renderTemplate(templateFile, sourceProfile, outputProfile); err != nil {
				fatalf("template: %v", err)
			}
		},
	}

	cmd.Flags().StringVarP(&outputProfile, "output", "o", "", "Write rendered output to a file (default: stdout)")
	rootCmd.AddCommand(cmd)
}

func renderTemplate(templateFile, sourceProfile, outputFile string) error {
	envMap, err := readEnvMap(profilePath(sourceProfile))
	if err != nil {
		return fmt.Errorf("reading profile %q: %w", sourceProfile, err)
	}

	tmplBytes, err := os.ReadFile(templateFile)
	if err != nil {
		return fmt.Errorf("reading template file %q: %w", templateFile, err)
	}

	rendered := expandTemplate(string(tmplBytes), envMap)

	if outputFile == "" {
		fmt.Print(rendered)
		return nil
	}

	return os.WriteFile(outputFile, []byte(rendered), 0644)
}

func expandTemplate(tmpl string, envMap map[string]string) string {
	return templateVarRe.ReplaceAllStringFunc(tmpl, func(match string) string {
		key := strings.TrimSpace(match[2 : len(match)-2])
		if val, ok := envMap[key]; ok {
			return val
		}
		return match
	})
}
