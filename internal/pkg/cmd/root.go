package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	releaseVersion string
)

var rootCmd = &cobra.Command{Use: "shell"}

// Execute is the entry point for the MFA command
func Execute(version string) {
	releaseVersion = version

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
