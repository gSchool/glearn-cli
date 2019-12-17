package cmd

import "github.com/spf13/cobra"

import "fmt"

import "os"

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new curriculum repository from a template",
	Long:  "Create a new curriculum repository from a template",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		hasConfig := doesCurrentDirHaveConfig()

		if hasConfig {
			fmt.Println("WARNING: configuration file detected and cannot continue with `learn new` command.")
		}
	},
}

func doesCurrentDirHaveConfig() bool {
	configExist := false
	currentDir, err := os.Getwd()
	if err != nil {
		return configExist
	}

	configYaml := currentDir + "/config.yaml"
	_, ymlExists := os.Stat(configYaml)
	if ymlExists == nil {
		configExist = true
	}

	configYml := currentDir + "/config.yml"
	_, ymlExists = os.Stat(configYml)
	if ymlExists == nil {
		configExist = true
	}

	autoConfigYaml := currentDir + "/autoconfig.yaml"
	_, ymlExists = os.Stat(autoConfigYaml)
	if ymlExists == nil {
		configExist = true
	}

	return configExist
}
