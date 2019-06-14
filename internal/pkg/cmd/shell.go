package cmd

import (
	"mfa4aws/internal/pkg/shell"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(shellCmd)
}

var shellCmd = &cobra.Command{
	Use:   "shell",
	Short: "Generates AWS STS access keys for use on the shell by wrapping the result in eval",
	Run: func(cmd *cobra.Command, args []string) {
		shell.GenerateSTSKeysForExport(awsProfile, mfaToken)
	},
}
