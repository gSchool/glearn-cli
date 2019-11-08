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
			return errors.New("Please set your API token first with `glearn settoken [token]`")
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

var setToken = &cobra.Command{
	Use:   "settoken [token]",
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

var new = &cobra.Command{
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

var preview = &cobra.Command{
	Use:   "preview [not_sure_yet]",
	Short: "Preview your content",
	Long:  `Long description for previewing`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Usage: `learn preview` takes one argument")
			os.Exit(1)
		}

		fmt.Println("Called preview")
	},
}

var build = &cobra.Command{
	Use:   "build [not_sure_yet]",
	Short: "Build your content",
	Long:  `Long description for building`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Usage: `learn build` takes one argument")
			os.Exit(1)
		}

		fmt.Println("Called build")
	},
}

func init() {
	u, err := user.Current()
	if err != nil {
		fmt.Println("Error retrieving your user path information")
		os.Exit(1)
	}

	viper.AddConfigPath(u.HomeDir)
	viper.SetConfigName(".glearn-config.yaml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found. Either user's first time using CLI or they deleted it
			configPath := fmt.Sprintf("%s/.glearn-config.yaml", u.HomeDir)
			initialConfig := []byte(`api_token:`)

			err = ioutil.WriteFile(configPath, initialConfig, 0666)
			if err != nil {
				fmt.Println("Error writing your glearn config file")
				os.Exit(1)
			}
		} else {
			// Config file was found but another error was produced
			fmt.Printf("Error: %s", err)
			os.Exit(1)
		}
	}

	rootCmd.AddCommand(setToken)
	rootCmd.AddCommand(new)
	rootCmd.AddCommand(preview)
	rootCmd.AddCommand(build)
}

// Execute runs the glearn CLI according to the user's command/subcommand/flags
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
