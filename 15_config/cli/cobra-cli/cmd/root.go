package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "my-cli",
	Short: "A simple CLI app built with Cobra",
	Long: `my-cli is a demonstration of how to build a CLI application
using the Cobra library in Go.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to my-cli! Use --help to see available commands.")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
