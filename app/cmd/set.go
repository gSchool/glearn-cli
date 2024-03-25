package cmd

import (
	"fmt"
	"os"

	appConfig "github.com/gSchool/glearn-cli/app/config"
	"github.com/spf13/cobra"
)

// UpgradeConfig is an initialized Boolean indicating that the user wants to
// upgrade their config
var UpgradeConfig bool

// APIToken is an initialized string used for holding it's flag value
var APIToken string

func NewSetCommand() *cobra.Command {
	setCmd.Flags().StringVarP(&APIToken, "api_token", "", "", "Your Learn api token")
	setCmd.Flags().BoolVarP(&UpgradeConfig, "upgrade", "", false, "Upgrade your CLI config file")
	return setCmd
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: fmt.Sprintf("Set your your credentials in %s", appConfig.ConfigPath()),
	Long: fmt.Sprintf(`In order to use learn resources through our CLI you must set your
credentials inside %s
	`, appConfig.ConfigPath()),
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			fmt.Fprintln(os.Stderr, "The set command does not take any arguments. Instead set variables with set --api_token=value")
			os.Exit(1)
		} else if !UpgradeConfig && APIToken == "" {
			cmd.Usage()
			os.Exit(1)
		}

		if UpgradeConfig {
			if updated, err := appConfig.Upgrade(); err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			} else if updated {
				fmt.Println("Configuration file upgraded")
			}
		}

		// If the --api_token=some_value flag was given, set it in app config
		if APIToken != "" {
			appConfig.Set("api_token", APIToken)
			// Write any changes made above to the config
			err := appConfig.Write()
			if err != nil {
				fmt.Fprintf(os.Stderr, "There was an error writing credentials to your config: %v", err)
				os.Exit(1)
				return
			}

			fmt.Println("Successfully added credentials!")
		}
	},
}
