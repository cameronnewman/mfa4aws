package cmd

import (
	"fmt"
	"mfa4aws/internal/pkg/awssts"
	"mfa4aws/internal/pkg/shell"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(shellCmd)
}

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Generates AWS STS access keys for use on the shell by wrapping the result in eval",
	Run: func(cmd *cobra.Command, args []string) {
		creds, err := awssts.GenerateSTSCredentials(awsProfile, mfaToken)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		shell.PrintVars(os.Stdout, shell.BuildEnvVars(creds))
	},
}
