package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/Galvanize-IT/glearn-cli/apis/learn"
	"github.com/spf13/cobra"
)

var branchCommand = `git branch | grep \* | cut -d ' ' -f2`
var remoteNameCommand = `git remote -v | grep push | cut -f2- -d/ | sed 's/[.].*$//'`

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Publish master for your curriculum repository",
	Long:  `The Learn system recognizes blocks of content held in GitHub respositories. This command pushes the latest commit for the remote origin master (which should be GitHub), then attemptes the release of a new Learn block version at the HEAD of master.`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 0 {
			fmt.Println("Usage: `learn build` takes no arguments, merely pushing latest master and releasing a version to Learn")
			os.Exit(1)
		}

		remote, err := remoteName()
		if err != nil {
			log.Println("Cannot run git remote detection with command: git remote -v | grep push | cut -f2- -d/ | sed 's/[.].*$//'\n", err)
			os.Exit(1)
		}
		if remote == "" {
			log.Println("no fetch remote detected")
			os.Exit(1)
		}
		// TODO refactor learn api block := learn.GetBlockByRepoName(remote)
		block, err := learn.Api.GetBlockByRepoName(remote)
		if err != nil {
			log.Println("Error fetchng block from learn", err)
			os.Exit(1)
		}
		// TODO if block does not exist, create one
		if block.Exists() {
			// block, err = CreateBlockByRepoName(remote)
			// if err != nil {
			// 	log.Println("Error creating block from learn", err)
			// 	os.Exit(1)
			// }
		}

		branch, err := currentBranch()
		if err != nil {
			log.Println("Cannot run git branch detection with bash:", err)
			os.Exit(1)
		}
		if branch != "publish#169468994" { // TODO change to master before merging
			fmt.Println("You are currently not on branch 'master'- the `learn build` command must be on master branch to push all currently committed work to your 'origin master' remote.")
			os.Exit(1)
		}
		fmt.Println("Pushing work to remote origin", branch)
		err = pushToRemote(branch)
		if err != nil {
			fmt.Printf("Error pushing to origin remote on branch %s: %s", err)
			os.Exit(1)
		}
		// create a release on learn, notify user
		// TODO resp, err := learn.CreateMasterRelease(remote)
		// if err != nil {
		// 	fmt.Printf("", err)
		// 	os.Exit(1)
		// }
	},
}

func currentBranch() (string, error) {
	return runBashCommand(branchCommand)
}

func remoteName() (string, error) {
	return runBashCommand(remoteNameCommand)
}

func runBashCommand(command string) (string, error) {
	out, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

func pushToRemote(branch string) error {
	out, err := exec.Command("bash", "-c", fmt.Sprintf("git push origin %s", branch)).CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}
