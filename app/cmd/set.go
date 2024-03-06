package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setCmd = &cobra.Command{
	Use:   "set --api_token=value",
	Short: "Set your your credentials for ~/.glearn-config.yaml",
	Long: `
In order to use learn resources through our CLI you must set your
credentials inside ~/.glearn-config.yaml
	`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			fmt.Fprintln(os.Stderr, "The set command does not take any arguments. Instead set variables with set --api_token=value")
			os.Exit(1)
		}

		// If the --api_token=some_value flag was given, set it in viper
		if APIToken == "" {
			fmt.Fprintln(os.Stderr, "The set command needs '--api_token' flag.\n\nUse: learn set --api_token=value")
			os.Exit(1)
		} else {
			viper.Set("api_token", APIToken)
		}

		// Write any changes made above to the config
		err := viper.WriteConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "There was an error writing credentials to your config: %v", err)
			os.Exit(1)
			return
		}

		fmt.Println("Successfully added credentials!")
	},
}
