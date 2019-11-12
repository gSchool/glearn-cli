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
			return errors.New("Please set your API token first with `glearn setapitoken [token]`")
		}

		if viper.Get("aws_access_key_id") == "" || viper.Get("aws_access_key_id") == nil {
			return errors.New(
				"Please set your AWS credentials first with `glearn setawsaccesskeyid [access_key_id] or by editing your ~/.glearn-config.yaml`",
			)
		}

		if viper.Get("aws_secret_access_key") == "" || viper.Get("aws_secret_access_key") == nil {
			return errors.New(
				"Please set your AWS credentials first with `glearn setawssecretaccesskey [secret_access_key] or by editing your ~/.glearn-config.yaml`",
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

func init() {
	u, err := user.Current()
	if err != nil {
		fmt.Println("Error retrieving your user path information")
		os.Exit(1)
		return
	}

	viper.AddConfigPath(u.HomeDir)
	viper.SetConfigName(".glearn-config.yaml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found. Either user's first time using CLI or they deleted it
			configPath := fmt.Sprintf("%s/.glearn-config.yaml", u.HomeDir)
			initialConfig := []byte(
				`
api_token:
aws_access_key_id:
aws_secret_access_key:
`,
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

	rootCmd.AddCommand(setApiTokenCmd)
	rootCmd.AddCommand(setAwsAccessKeyId)
	rootCmd.AddCommand(setAwsSecretAccessKey)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(previewCmd)
	rootCmd.AddCommand(buildCmd)
}

// Execute runs the glearn CLI according to the user's command/subcommand/flags
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
