package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setBucketName = &cobra.Command{
	Use:   "setbucket [bucket_name]",
	Short: "Set your AWS bucket name",
	Long: `
		In order to use learn resources through our CLI you
		must set your AWS bucket name
	`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Please provide only one argument (the s3 bucket name)")
			os.Exit(1)
		}

		viper.Set("aws_s3_key_prefix", args[0])
		viper.WriteConfig()

		fmt.Println("Successfully set AWS key prefix!")
	},
}
