package main

import (
	"fmt"
	"os"

	"github.com/envoy-cli/envoy/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
