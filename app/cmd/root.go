package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"time"

	"github.com/gSchool/glearn-cli/api/learn"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const setAPITokenMessage = `
Please set your API token with this command: learn set --api_token=your_api_token 
You can get your api token at https://learn-2.galvanize.com/api_token
`

// currentReleaseVersion is used to print the version the user currently has downloaded
const currentReleaseVersion = "v0.6.5"

// rootCmd is the base for all our commands. It currently just checks for all the
// necessary credentials and prompts the user to set them if they are not there.
var rootCmd = &cobra.Command{
	Use:   "learn [command]",
	Short: "learn is a CLI tool for communicating with Learn",
	Long:  "learn is a CLI tool for communicating with Learn",
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

var fileExtWhitelist = map[string]struct{}{
	".yaml":  struct{}{},
	".yml":   struct{}{},
	".md":    struct{}{},
	".pdf":   struct{}{},
	".ipynb": struct{}{},
	".jpg":   struct{}{},
	".jpeg":  struct{}{},
	".jpe":   struct{}{},
	".jif":   struct{}{},
	".jfif":  struct{}{},
	".jfi":   struct{}{},
	".png":   struct{}{},
	".gif":   struct{}{},
	".tiff":  struct{}{},
	".tif":   struct{}{},
	".bmp":   struct{}{},
	".svg":   struct{}{},
	".svgz":  struct{}{},
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

	// Add all the other learn commands defined in cmd/ directory
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(previewCmd)
	rootCmd.AddCommand(publishCmd)
	rootCmd.AddCommand(versionCmd)

	// Check for flags set by the user and hyrate their corresponding variables.
	setCmd.Flags().StringVarP(&APIToken, "api_token", "", "", "Your Learn api token")
	previewCmd.Flags().StringVarP(&UnitsDirectory, "units", "u", "", "The directory where your units exist")
	previewCmd.Flags().BoolVarP(&OpenPreview, "open", "o", false, "Open the preview in the browser")
	previewCmd.Flags().BoolVarP(&FileOnly, "fileonly", "x", false, "E(x)cludes images when previewing a single file, defaults false")
	publishCmd.Flags().StringVarP(&UnitsDirectory, "units", "u", "", "The directory where your units exist")
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

	learn.API = api
}
