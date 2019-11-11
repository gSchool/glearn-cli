package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setApiTokenCmd = &cobra.Command{
	Use:   "setapitoken [token]",
	Short: "Set your API token",
	Long: `
		In order to use learn resources through our CLI you
		must create and set an API token for yourself
	`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Please provide only one argument (your api token)")
			os.Exit(1)
		}

		viper.Set("api_token", args[0])
		viper.WriteConfig()

		fmt.Println("Successfully set API token!")
	},
}
