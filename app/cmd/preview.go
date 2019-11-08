package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var previewCmd = &cobra.Command{
	Use:   "preview [not_sure_yet]",
	Short: "Preview your content",
	Long:  `Long description for previewing`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Usage: `learn preview` takes one argument")
			os.Exit(1)
		}

		fmt.Println("Called preview")
	},
}
