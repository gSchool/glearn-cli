package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Publish master for your curriculum repository",
	Long:  `The Learn system recognizes blocks of content held in GitHub respositories. This command publishes the latest commit on master as a new release.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Usage: `learn build` takes one argument")
			os.Exit(1)
		}

		fmt.Println("Called build")
	},
}
