package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type SchemaKey struct {
	Required    bool   `json:"required"`
	Description string `json:"description,omitempty"`
	Default     string `json:"default,omitempty"`
}

type EnvSchema struct {
	Keys map[string]SchemaKey `json:"keys"`
}

func schemaFilePath(dir string) string {
	return filepath.Join(dir, ".envoy", "schema.json")
}

func init() {
	var schemaCmd = &cobra.Command{
		Use:   "schema",
		Short: "Manage the env schema for a project",
	}

	var defineCmd = &cobra.Command{
		Use:   "define <key> [--required] [--desc <text>] [--default <val>]",
		Short: "Define or update a key in the schema",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			key := strings.TrimSpace(args[0])
			required, _ := cmd.Flags().GetBool("required")
			desc, _ := cmd.Flags().GetString("desc")
			def, _ := cmd.Flags().GetString("default")
			dir := projectDir()
			if err := defineSchemaKey(dir, key, required, desc, def); err != nil {
				fatalf("schema define: %v", err)
			}
			fmt.Printf("Schema key '%s' defined.\n", key)
		},
	}
	defineCmd.Flags().Bool("required", false, "Mark key as required")
	defineCmd.Flags().String("desc", "", "Description of the key")
	defineCmd.Flags().String("default", "", "Default value for the key")

	var validateSchemaCmd = &cobra.Command{
		Use:   "validate <profile>",
		Short: "Validate a profile against the schema",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dir := projectDir()
			errs := validateAgainstSchema(dir, args[0])
			if len(errs) == 0 {
				fmt.Println("Profile is valid against schema.")
				return
			}
			for _, e := range errs {
				fmt.Println(" -", e)
			}
			os.Exit(1)
		},
	}

	var showSchemaCmd = &cobra.Command{
		Use:   "show",
		Short: "Show the current schema",
		Run: func(cmd *cobra.Command, args []string) {
			dir := projectDir()
			schema, err := loadSchema(dir)
			if err != nil {
				fatalf("schema show: %v", err)
			}
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			_ = enc.Encode(schema)
		},
	}

	schemaCmd.AddCommand(defineCmd, validateSchemaCmd, showSchemaCmd)
	rootCmd.AddCommand(schemaCmd)
}
