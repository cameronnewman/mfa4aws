package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	awsProfile string
	mfaToken   string

	releaseVersion string
)

var rootCmd = &cobra.Command{Use: "shell"}

// Execute is the entry point for the MFA command
func Execute(version string) {
	rootCmd.PersistentFlags().StringVarP(&awsProfile, "profile", "p", "default", "AWS Profile name in $HOME/.aws/credentials")
	rootCmd.PersistentFlags().StringVarP(&mfaToken, "token", "t", "", "Current MFA value to use for STS generation")

	releaseVersion = version

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
