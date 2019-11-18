package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var setCmd = &cobra.Command{
	Use:   "set [...flags]",
	Short: "Set your your credentials for ~/.glearn-config.yaml",
	Long: `
		In order to use learn resources through our CLI you
		must set your credentials inside ~/.glearn-config.yaml
	`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			fmt.Println("The set command does not take any arguments. Instead set variables with set --credentialFlag=value")
			os.Exit(1)
		}

		if APIToken != "" {
			viper.Set("api_token", APIToken)
		}

		if AwsAccessKeyID != "" {
			viper.Set("aws_access_key_id", AwsAccessKeyID)
		}

		if AwsSecretAccessKey != "" {
			viper.Set("aws_secret_access_key", AwsSecretAccessKey)
		}

		if AwsS3Bucket != "" {
			viper.Set("aws_s3_bucket", AwsS3Bucket)
		}

		if AwsS3KeyPrefix != "" {
			viper.Set("aws_s3_key_prefix", AwsS3KeyPrefix)
		}

		err := viper.WriteConfig()
		if err != nil {
			fmt.Printf("There was an error writing credentials to your config: %v", err)
			os.Exit(1)
			return
		}

		fmt.Println("Successfully wrote credentials!")
	},
}
