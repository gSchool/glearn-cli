package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [not_sure_yet]",
	Short: "Create something new",
	Long:  `Long description for creating new`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Usage: `learn new` takes one argument")
			os.Exit(1)
		}

		fmt.Println("Called new")
	},
}
