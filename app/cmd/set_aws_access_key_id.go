package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setAwsAccessKeyId = &cobra.Command{
	Use:   "setawsaccesskeyid [access_key_id]",
	Short: "Set your AWS access key id",
	Long: `
		In order to use learn resources through our CLI you
		must set your AWS access key id
	`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Please provide only one argument (your aws access key id)")
			os.Exit(1)
		}

		viper.Set("aws_access_key_id", args[0])
		viper.WriteConfig()

		fmt.Println("Successfully set AWS access key ID!")
	},
}
