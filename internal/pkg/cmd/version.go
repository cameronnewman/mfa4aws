package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "display release version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v" + releaseVersion)
	},
}
