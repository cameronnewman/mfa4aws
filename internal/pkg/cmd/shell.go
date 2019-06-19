package cmd

import (
	"fmt"
	"mfa4aws/internal/pkg/aws"
	"mfa4aws/internal/pkg/shell"
	"os"

	"github.com/spf13/cobra"
)

var (
	awsProfile string
	mfaToken   string
)

func init() {
	rootCmd.AddCommand(shellCmd)

	persistentFlags := shellCmd.PersistentFlags()
    persistentFlags.StringVarP(&awsProfile, "profile", "p", "default", "AWS Profile name in $HOME/.aws/credentials")
    persistentFlags.StringVarP(&mfaToken, "token", "t", "", "Current MFA value to use for STS generation")
    cobra.MarkFlagRequired(persistentFlags, "token")
}

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Generates AWS STS access keys for use on the shell by wrapping the result in eval",
	Run: func(cmd *cobra.Command, args []string) {
		creds, err := aws.GenerateSTSCredentials(awsProfile, mfaToken)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		shell.PrintVars(os.Stdout, shell.BuildEnvVars(creds))
	},
}
