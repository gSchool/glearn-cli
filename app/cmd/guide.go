package cmd

import (
	"fmt"
	"os"
	"os/exec"
	// "path/filepath"
	// "strings"

	"github.com/spf13/cobra"
)

var guideCmd = &cobra.Command{
	Use:     "walkthrough",
	Aliases: []string{"guide"},
	Short:   "Download examples for use in the walkthrough",
	Long:    "Download examples for use in the walkthrough",
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		// Get the current directory
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println("Could not detect a working directory")
			os.Exit(1)
		}

		// Does that directory have a walkthrough
		hasWalkthrough := doesCurrentDirHaveWalkthrough(currentDir)

		if hasWalkthrough {
			fmt.Println("'learn-walkthrough' folder already exists, cannot continue with command.")
			os.Exit(1)
		}

		// Clone the template from github
		fmt.Println("Copying curriculum template from Github")
		fmt.Println("=======================================")
		fmt.Println("\nCloning into 'learn-walkthrough'...")
		err = cloneTemplate()
		if err != nil {
			fmt.Println("We had trouble cloning into learn-walkthrough, please check that you have the correct github credentials")
			os.Exit(1)
		}

		// remove git folder
		err = removeGit(currentDir)
		if err != nil {
			fmt.Println("Could not remove git folder")
			os.Exit(1)
		}

		fmt.Println(`
Success!
========

A small example curriculum has been added in ./learn-walkthrough.`)
	},
}

func doesCurrentDirHaveWalkthrough(currentDir string) (bool) {
	walkthroughExist := false

	walkthrough := currentDir + "/learn-walkthrough/"
	_, walkthroughExists := os.Stat(walkthrough)
	if walkthroughExists == nil {
		walkthroughExist = true
	}

	return walkthroughExist
}

func cloneTemplate() error {
	_, err := exec.Command("bash", "-c", "git clone git@github.com:gSchool/learn-walkthrough.git").CombinedOutput()
	if err != nil {
		_, errr := exec.Command("bash", "-c", "git clone https://github.com/gSchool/learn-walkthrough.git").CombinedOutput()
		if errr != nil {
			return errr
		}
	}

	return nil
}

func removeGit(currentDir string) error {
	initDir := "/learn-walkthrough"
	os.RemoveAll(currentDir + initDir + "/.git/")
	return nil
}
