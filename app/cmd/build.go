package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var branchCommand = `git branch | grep \* | cut -d ' ' -f2`

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Publish master for your curriculum repository",
	Long:  `The Learn system recognizes blocks of content held in GitHub respositories. This command publishes the latest commit on master to your origin remote (which should be GitHub), then releases a new Learn block version at the HEAD of master.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			fmt.Println("Usage: `learn build` takes no arguments, merely pushing latest master and releasing a version to Learn")
			os.Exit(1)
		}

		if currentBranch() != "master" {
			fmt.Println("You are currently not on branch 'master'- the `learn build` command must be on master branch to push all currently committed work to your 'origin master' remote.")
		}
		fmt.Println(currentBranch())
	},
}

func currentBranch() string {
	out, err := exec.Command("bash", "-c", branchCommand).Output()
	if err != nil {
		log.Fatal("Cannot run git branch detection with bash:", err)
	}

	return strings.TrimSpace(string(out))
}
