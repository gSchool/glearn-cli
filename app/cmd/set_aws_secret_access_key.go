package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setAwsSecretAccessKey = &cobra.Command{
	Use:   "setawssecretaccesskey [secret_access_key]",
	Short: "Set your AWS secret access key",
	Long: `
		In order to use learn resources through our CLI you
		must set your AWS secret access key
	`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Please provide only one argument (your aws secret access key)")
			os.Exit(1)
		}

		viper.Set("aws_secret_access_key", args[0])
		viper.WriteConfig()

		fmt.Println("Successfully set AWS secret access key!")
	},
}
