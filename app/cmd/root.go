package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Masterminds/semver"
	"github.com/gSchool/glearn-cli/api/github"
	"github.com/gSchool/glearn-cli/api/learn"
	"github.com/gSchool/glearn-cli/app/cmd/markdown"
	"github.com/spf13/cobra"
)

const setAPITokenMessage = `
Please set your API token with this command: learn set --api_token=<your_api_token>
You can get your api token at https://learn-2.galvanize.com/api_token`

// currentReleaseVersion is used to print the version the user currently has downloaded
const currentReleaseVersion = "v0.10.13"

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

Learn more by running 'learn walkthrough' to create sample materials, or visit
  https://learn-2.galvanize.com/cohorts/667/blocks/13/content_files/walkthrough/01-overview.md`,
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
	// Add all the other learn commands defined in cmd/ directory
	rootCmd.AddCommand(markdown.NewMarkdownCommand())
	rootCmd.AddCommand(NewSetCommand())
	rootCmd.AddCommand(previewCmd)
	rootCmd.AddCommand(publishCmd)
	rootCmd.AddCommand(guideCmd)
	rootCmd.AddCommand(versionCmd)

	// Check for flags set by the user and hydrate their corresponding variables.
	previewCmd.Flags().StringVarP(&UnitsDirectory, "units", "u", "", "The directory where your units exist")
	previewCmd.Flags().BoolVarP(&OpenPreview, "open", "o", false, "Open the preview in the browser")
	previewCmd.Flags().BoolVarP(&FileOnly, "fileonly", "x", false, "Excludes images when previewing a single file, defaults false")
	publishCmd.Flags().StringVarP(&UnitsDirectory, "units", "u", "", "The directory where your units exist")
	publishCmd.Flags().BoolVarP(&IgnoreLocal, "ignore-local", "", false, "Ignore local changes and publish remote only")
	publishCmd.Flags().BoolVarP(&CiCdEnvironment, "ci-cd", "", false, "Running in a CI/CD environment (cannot use with autoconfig feature)")
}

// Execute runs the learn CLI according to the user's command/subcommand/flags
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func setupLearnAPI(getPresignedPostUrl bool) {
	client := http.Client{Timeout: 15 * time.Second}
	baseURL := "https://learn-2.galvanize.com"
	alternateURL := os.Getenv("LEARN_BASE_URL")
	if alternateURL != "" {
		baseURL = alternateURL
	}

	api, err := learn.NewAPI(baseURL, &client, getPresignedPostUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating API client. Err: %v", err)
		os.Exit(1)
		return
	}

	githubAPI := github.NewAPI(&client)
	version, err := githubAPI.GetLatestVersion()
	if err != nil {
		fmt.Printf("Something went wrong when fetching latest CLI version: %s\n", err)
	} else if version != currentReleaseVersion {
		versionRemote, versionRemoteErr := semver.NewVersion(version)
		versionInstalled, versionInstalledErr := semver.NewVersion(currentReleaseVersion)
		if versionRemoteErr != nil {
			fmt.Printf("Failed to parse the CLI's current version. Err: %v", err)
		} else if versionInstalledErr != nil {
			fmt.Printf("Failed to parse the latest CLI release version. Err: %v", err)
		} else if versionInstalled.LessThan(versionRemote) {
			fmt.Printf("\nWARNING: There is newer version of the learn tool available.\nLatest: %s\nCurrent: %s\nTo avoid issues, upgrade by following the instructions at this link:\nhttps://github.com/gSchool/glearn-cli/blob/master/upgrade_instructions.md\n\n", version, currentReleaseVersion)
		}
	}

	learn.API = api
}
