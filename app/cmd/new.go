package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new curriculum repository from a template",
	Long:  "Create a new curriculum repository from a template",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		hasConfig := doesCurrentDirHaveConfig()

		if hasConfig {
			fmt.Println("WARNING: configuration file detected and cannot continue with `learn new` command.")
			os.Exit(1)
		}

		fmt.Println("Copying curriculum template from Github")
		fmt.Println("=======================================")
		fmt.Println("\nCloning into 'learn-curriculum-init'...")
		err := cloneTemplate()
		if err != nil {
			fmt.Println("We had trouble cloning into learn-curriculum-init, please check that you have the correct github credentials")
			os.Exit(1)
		}

		fmt.Println("Copying curriculum")
		err = moveClonedMaterials()
		if err != nil {
			fmt.Println("Could not move template into working repository")
			os.Exit(1)
		}
		fmt.Println("Removing cloned repo")

		fmt.Println(`
Success!
========

See repo structure guidelines in readme.md

Run 'learn preview .' to preview curriculum without pushing
Run 'learn publish' to push to Github and publish to Learn
		`)
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

func cloneTemplate() error {
	_, err := exec.Command("bash", "-c", "git clone git@github.com:gSchool/learn-curriculum-init.git").CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

func moveClonedMaterials() error {
	currentDir, _ := os.Getwd()
	initDir := "/learn-curriculum-init"
	os.RemoveAll(currentDir + initDir + "/.git/")
	err := filepath.Walk(currentDir+initDir, func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, initDir) && !strings.Contains(path, ".git/") {
			oldLocation := path
			newLocation := strings.Replace(path, initDir, "", 1)
			os.Rename(oldLocation, newLocation)
		}
		return nil
	})
	if err != nil {
		return err
	}
	os.RemoveAll(currentDir + initDir)
	return nil
}
