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

		branch, err := currentBranch()
		if err != nil {
			log.Fatal("Cannot run git branch detection with bash:", err)
		}

		if branch != "publish#169468994" { // TODO change to master before merging
			fmt.Println("You are currently not on branch 'master'- the `learn build` command must be on master branch to push all currently committed work to your 'origin master' remote.")
			os.Exit(1)
		}
		fmt.Println("Pushing work to origin remote for ", branch)
		err = pushToRemote(branch)
		if err != nil {
			fmt.Printf("Error pushing to origin remote on branch %s: %s", err)
			os.Exit(1)
		}
	},
}

func currentBranch() (string, error) {
	out, err := exec.Command("bash", "-c", branchCommand).Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

func pushToRemote(branch string) error {
	out, err := exec.Command("bash", "-c", fmt.Sprintf("git push origin %s", branch)).Output()
	if err != nil {
		return "", err
	}

	fmt.Println("push out:", strings.TrimSpace(string(out)))
	return nil
}

// func pushToRemote(branch string) error {
// 	cmd := exec.Command("git", "push", "origin", branch)
// 	stdout, err := cmd.StdoutPipe()
// 	if err != nil {
// 		return err
// 	}
// 	defer stdout.Close()
// 	if err := cmd.Start(); err != nil {
// 		return err
// 	}
// 	resp, err := ioutil.ReadAll(stdout)
// 	if err != nil {
// 		return err
// 	}
// 	if err := cmd.Wait(); err != nil {
// 		return err
// 	}
// 	fmt.Println("stdout:", string(resp))
// 	return nil
// }
