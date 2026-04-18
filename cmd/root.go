package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "envoy",
	Short: "envoy — manage and switch between .env file profiles",
	Long: `envoy is a lightweight CLI for managing and switching between
.env file profiles across projects. Use it to create, list,
switch, and delete named environment profiles.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(useCmd)
	rootCmd.AddCommand(addCmd)

	rootCmd.PersistentFlags().StringP("dir", "d", ".", "project directory to operate in")
}

func projectDir(cmd *cobra.Command) string {
	dir, err := cmd.Flags().GetString("dir")
	if err != nil || dir == "" {
		wd, _ := os.Getwd()
		return wd
	}
	return dir
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}
