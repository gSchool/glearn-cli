package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build [not_sure_yet]",
	Short: "Build your content",
	Long:  `Long description for building`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Usage: `learn build` takes one argument")
			os.Exit(1)
		}

		fmt.Println("Called build")
	},
}
