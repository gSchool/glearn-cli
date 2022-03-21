package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"time"

	"github.com/gSchool/glearn-cli/api/github"
	"github.com/gSchool/glearn-cli/api/learn"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const setAPITokenMessage = `
Please set your API token with this command: learn set --api_token=<your_api_token>
You can get your api token at https://learn-2.galvanize.com/api_token
`

// currentReleaseVersion is used to print the version the user currently has downloaded
const currentReleaseVersion = "v0.10.1"

// rootCmd is the base for all our commands. It currently just checks for all the
// necessary credentials and prompts the user to set them if they are not there.
var rootCmd = &cobra.Command{
	Use:   "learn [command]",
	Short: "learn is a CLI tool for communicating with Learn",
	Long: `learn is a CLI tool for communicating with Learn

Edit existing curriculum:
  1. Clone and edit curriculum
  2. Preview your changes. Run:
      learn preview -o <directory|file>
  3. Git add / commit / push changes to the master branch
  4. Publish changes for any cohort in Learn. Run:
      learn publish

Learn more and build new curriculum:
  walkthrough at https://galvanize-learn.zendesk.com/hc/en-us/articles/1500000930401-Introduction`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Requires at least 1 argument")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Unknown command. Try `learn help` for more information")
	},
}

// APIToken is an initialized string used for holding it's flag value
var APIToken string

// UnitsDirectory is a flag for preview command that denotes a location for the units
var UnitsDirectory string

// FileOnly is the flag boolean which will force a single file upload to ignore any images
// and only upload the markdown file
var FileOnly bool

// OpenPreview is the flag boolean which will open the preview in browser
var OpenPreview bool

// Ignore local changes and publish remote only
var IgnoreLocal bool

// Running in a CI environment and should not try to push changes
var CiCdEnvironment bool

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
			err = ioutil.WriteFile(configPath, initialConfig, 0600)
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

	// Add all the other learn commands defined in cmd/ directory
	rootCmd.AddCommand(markdownCmd)
	rootCmd.AddCommand(previewCmd)
	rootCmd.AddCommand(publishCmd)
	rootCmd.AddCommand(guideCmd)
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(versionCmd)

	// Check for flags set by the user and hydrate their corresponding variables.
	setCmd.Flags().StringVarP(&APIToken, "api_token", "", "", "Your Learn api token")
	previewCmd.Flags().StringVarP(&UnitsDirectory, "units", "u", "", "The directory where your units exist")
	previewCmd.Flags().BoolVarP(&OpenPreview, "open", "o", false, "Open the preview in the browser")
	previewCmd.Flags().BoolVarP(&FileOnly, "fileonly", "x", false, "Excludes images when previewing a single file, defaults false")
	publishCmd.Flags().StringVarP(&UnitsDirectory, "units", "u", "", "The directory where your units exist")
	publishCmd.Flags().BoolVarP(&IgnoreLocal, "ignore-local", "", false, "Ignore local changes and publish remote only")
	publishCmd.Flags().BoolVarP(&CiCdEnvironment, "ci-cd", "", false, "Running in a CI/CD environment (cannot use with autoconfig feature)")
	markdownCmd.Flags().BoolVarP(&PrintTemplate, "out", "o", false, "Prints the template to stdout")
	markdownCmd.Flags().BoolVarP(&Minimal, "min", "m", false, "Uses a terse, minimal version of the template")
}

// Execute runs the learn CLI according to the user's command/subcommand/flags
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func setupLearnAPI() {
	client := http.Client{Timeout: 15 * time.Second}
	baseURL := "https://learn-2.galvanize.com"
	alternateURL := os.Getenv("LEARN_BASE_URL")
	if alternateURL != "" {
		baseURL = alternateURL
	}

	api, err := learn.NewAPI(baseURL, &client)
	if err != nil {
		fmt.Printf("Error creating API client. Err: %v", err)
		os.Exit(1)
		return
	}

	githubAPI := github.NewAPI(&client)
	version, err := githubAPI.GetLatestVersion()
	if err != nil {
		fmt.Printf("Something went wrong when fetching latest CLI version: %s\n", err)
	} else if version != currentReleaseVersion {
		fmt.Printf("\nWARNING: There is newer version of the learn tool available.\nLatest: %s\nCurrent: %s\nTo avoid issues, upgrade by following the instructions at this link:\nhttps://github.com/gSchool/glearn-cli/blob/master/upgrade_instructions.md\n\n", version, currentReleaseVersion)
	}

	learn.API = api
}
