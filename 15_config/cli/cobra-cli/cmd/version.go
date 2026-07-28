package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const version = "v1.0.0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Prints the current version of my-cli.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("my-cli %s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
