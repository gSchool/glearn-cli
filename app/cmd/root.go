package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "glearn",
	Short: "glearn is a cli application for Learn",
	Long:  `A longer description of what glearn is`,
	Args: func(cmd *cobra.Command, args []string) error {
		if viper.Get("api_token") == "" || viper.Get("api_token") == nil {
			return errors.New("Please set your API token first with `glearn set --api_token=value`")
		}

		if viper.Get("aws_access_key_id") == "" || viper.Get("aws_access_key_id") == nil {
			return errors.New(
				"Please set your AWS access key ID first with `glearn set --access_key_id=value or by editing your ~/.glearn-config.yaml`",
			)
		}

		if viper.Get("aws_secret_access_key") == "" || viper.Get("aws_secret_access_key") == nil {
			return errors.New(
				"Please set your AWS secret access key first with `glearn set --secret_access_key=value or by editing your ~/.glearn-config.yaml`",
			)
		}

		if viper.Get("aws_s3_bucket") == "" || viper.Get("aws_s3_bucket") == nil {
			return errors.New(
				"Please set your AWS s3 bucket first with `glearn set --s3_bucket=value or by editing your ~/.glearn-config.yaml`",
			)
		}

		if viper.Get("aws_s3_key_prefix") == "" || viper.Get("aws_s3_key_prefix") == nil {
			return errors.New(
				"Please set your AWS s3 key prefix first with `glearn set --s3_prefix=value or by editing your ~/.glearn-config.yaml`",
			)
		}

		if len(args) < 1 {
			return errors.New("Requires at least 1 argument")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ran main command")
	},
}

// APIToken is an initialized string used for holding it's flag value
var APIToken string

// AwsAccessKeyID is an initialized string used for holding it's flag value
var AwsAccessKeyID string

// AwsSecretAccessKey is an initialized string used for holding it's flag value
var AwsSecretAccessKey string

// AwsS3Bucket is an initialized string used for holding it's flag value
var AwsS3Bucket string

// AwsS3KeyPrefix is an initialized string used for holding it's flag value
var AwsS3KeyPrefix string

func init() {
	u, err := user.Current()
	if err != nil {
		fmt.Println("Error retrieving your user path information")
		os.Exit(1)
		return
	}

	viper.AddConfigPath(u.HomeDir)
	viper.SetConfigName(".glearn-config")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found. Either user's first time using CLI or they deleted it
			configPath := fmt.Sprintf("%s/.glearn-config.yaml", u.HomeDir)
			initialConfig := []byte(
				`api_token:
aws_access_key_id:
aws_secret_access_key:
aws_s3_bucket:
aws_s3_key_prefix:`,
			)

			err = ioutil.WriteFile(configPath, initialConfig, 0666)
			if err != nil {
				fmt.Println("Error writing your glearn config file")
				os.Exit(1)
				return
			}
		} else {
			// Config file was found but another error was produced
			fmt.Printf("Error: %s", err)
			os.Exit(1)
			return
		}
	}

	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(previewCmd)
	rootCmd.AddCommand(buildCmd)

	setCmd.Flags().StringVarP(&APIToken, "api_token", "", "", "Your Learn api token")
	setCmd.Flags().StringVarP(&AwsAccessKeyID, "access_key_id", "", "", "Access key ID for glearn-cli")
	setCmd.Flags().StringVarP(&AwsSecretAccessKey, "secret_access_key", "", "", "Secret access key for glearn-cli")
	setCmd.Flags().StringVarP(&AwsS3Bucket, "s3_bucket", "", "", "S3 bucket name for glearn-cli")
	setCmd.Flags().StringVarP(&AwsS3KeyPrefix, "s3_prefix", "", "", "S3 bucket key prefix for glearn-cli")
}

// Execute runs the glearn CLI according to the user's command/subcommand/flags
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
