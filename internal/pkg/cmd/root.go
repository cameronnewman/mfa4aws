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
    persistentFlags := rootCmd.PersistentFlags()
    persistentFlags.StringVarP(&awsProfile, "profile", "p", "default", "AWS Profile name in $HOME/.aws/credentials")
    persistentFlags.StringVarP(&mfaToken, "token", "t", "", "Current MFA value to use for STS generation")
    cobra.MarkFlagRequired(persistentFlags, "token")

    releaseVersion = version

    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
