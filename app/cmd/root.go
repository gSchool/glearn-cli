package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"time"

	"github.com/Galvanize-IT/glearn-cli/apis/learn"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd is the base for all our commands. It currently just checks for all the
// necessary credentials and prompts the user to set them if they are not there.
var rootCmd = &cobra.Command{
	Use:   "glearn",
	Short: "glearn is a cli application for Learn",
	Long:  `A longer description of what glearn is`,
	Args: func(cmd *cobra.Command, args []string) error {
		if viper.Get("api_token") == "" || viper.Get("api_token") == nil {
			return errors.New("Please set your API token first with `glearn set --api_token=value`")
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
			initialConfig := []byte(`api_token:`)

			// Write a ~/.glearn-config.yaml file with all the needed credential keys to fill in.
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

	apiToken, ok := viper.Get("api_token").(string)
	if !ok {
		fmt.Println("Please set your api_token in ~/.glearn-config.yaml")
		os.Exit(1)
	}

	client := http.Client{Timeout: 15 * time.Second}
	baseUrl := "https://learn-2.galvanize.com"
	alternateUrl := os.Getenv("LEARN_BASE_URL")
	if alternateUrl != "" {
		baseUrl = alternateUrl
	}
	learn.Api = learn.NewAPI(apiToken, baseUrl, client)

	// Add all the other glearn commands defined in cmd/ directory
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(previewCmd)
	rootCmd.AddCommand(buildCmd)

	// Check for flags set by the user and hyrate their corresponding variables.
	setCmd.Flags().StringVarP(&APIToken, "api_token", "", "", "Your Learn api token")
}

// Execute runs the glearn CLI according to the user's command/subcommand/flags
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
