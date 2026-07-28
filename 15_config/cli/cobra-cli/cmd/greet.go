package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var name string
var times int

var greetCmd = &cobra.Command{
	Use:   "greet",
	Short: "Greet someone",
	Long:  `Prints a greeting message. You can specify the name and number of times to greet.`,
	Run: func(cmd *cobra.Command, args []string) {
		for i := 0; i < times; i++ {
			fmt.Printf("Hello, %s!\n", name)
		}
	},
}

func init() {
	rootCmd.AddCommand(greetCmd)

	// Флаги для команды greet
	greetCmd.Flags().StringVarP(&name, "name", "n", "World", "Name to greet")
	greetCmd.Flags().IntVarP(&times, "times", "t", 1, "Number of times to greet")
}
