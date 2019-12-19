package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/gSchool/glearn-cli/api/learn"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	branchCommand     = `git branch | grep \* | cut -d ' ' -f2`
	pushRemoteCommand = `git remote get-url --push origin`
)

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish master for your curriculum repository",
	Long: `
The Learn system recognizes blocks of content held in GitHub respositories. This
command pushes the latest commit for the remote origin master (which should be
GitHub), then attempts the release of a new Learn block version at the HEAD of
master. If the block doesn't exist, running the publish command will create a
new block. If the block already exists, it will update the existing block.
	`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if viper.Get("api_token") == "" || viper.Get("api_token") == nil {
			fmt.Println(setAPITokenMessage)
			os.Exit(1)
		}

		if len(args) != 0 {
			fmt.Println("Usage: `learn publish` takes no arguments, merely pushing latest master and releasing a version to Learn. Use the command from inside a block repository.")
			os.Exit(1)
		}

		// Start benchmarking the total time spent in publish cmd
		startOfCmd := time.Now()

		remote, err := remoteName()
		if err != nil {
			log.Printf("Cannot run git remote detection with command: %s\n%s\n", pushRemoteCommand, err)
			os.Exit(1)
		}
		if remote == "" {
			log.Println("no fetch remote detected")
			os.Exit(1)
		}
		fmt.Printf("Publishing block with repo name %s\n", remote)

		block, err := learn.API.GetBlockByRepoName(remote)
		if err != nil {
			log.Printf("Error fetchng block from learn: %s\n", err)
			os.Exit(1)
		}
		if !block.Exists() {
			block, err = learn.API.CreateBlockByRepoName(remote)
			if err != nil {
				log.Printf("Error creating block from learn: %s\n", err)
				os.Exit(1)
			}
		}

		branch, err := currentBranch()
		if err != nil {
			log.Println("Cannot run git branch detection with bash:", err)
			os.Exit(1)
		}

		// Detect config file
		path, _ := os.Getwd()
		createdConfig, err := doesConfigExistOrCreate(path+"/", UnitsDirectory)
		if err != nil {
			log.Printf(fmt.Sprintf("Failed to find or create a config file for repo: (%s). Err: %v", branch, err))
			os.Exit(1)
		}

		if branch != "master" {
			fmt.Println("You are currently not on branch 'master'- the `learn publish` command must be on master branch to push all currently committed work to your 'origin master' remote.")
			os.Exit(1)
		}

		if createdConfig {
			fmt.Println("We are going to attempt to commit the autoconfig.yaml to", branch)
			err = addAutoConfigAndCommit()

			if err != nil {
				fmt.Printf("Error committing the autoconfig.yaml to origin remote on branch: %s", err)
				os.Exit(1)
			}
		}

		fmt.Println("Pushing work to remote origin", branch)

		// TODO what happens when they do not have work in remote and push fails?
		err = pushToRemote(branch)
		if err != nil {
			fmt.Printf("Error pushing to origin remote on branch: %s\n", err)
			os.Exit(1)
		}

		// Start benchmark for creating master release & building on learn
		startOfMasterReleaseAndBuild := time.Now()

		// Start a processing spinner that runs until Learn is finsihed building the preview
		fmt.Println("\nPlease wait while Learn builds your release...")
		s := spinner.New(spinner.CharSets[32], 100*time.Millisecond)
		s.Color("green")
		s.FinalMSG = fmt.Sprintf("Block %d released!\n", block.ID)
		s.Start()

		// Create a release on learn, notify user
		releaseID, err := learn.API.CreateMasterRelease(block.ID)
		if err != nil || releaseID == 0 {
			fmt.Printf("error creating master release for releaseID: %d. Error: %s\n", releaseID, err)
			os.Exit(1)
		}

		var attempts uint8 = 20
		p, err := learn.API.PollForBuildResponse(releaseID, &attempts)
		if err != nil {
			s.Stop()

			block, err := learn.API.GetBlockByRepoName(remote)
			if err != nil {
				log.Printf("Error fetchng block from learn: %s\n", err)
				os.Exit(1)
			}
			fmt.Println("Errors on block:")
			for _, e := range block.SyncErrors {
				fmt.Println(e)
			}
			os.Exit(1)
		}

		// Add benchmark in milliseconds for compressDirectory
		bench := &learn.CLIBenchmark{
			MasterReleaseAndBuild: time.Since(startOfMasterReleaseAndBuild).Milliseconds(),
			TotalCmdTime:          time.Since(startOfCmd).Milliseconds(),
			CmdName:               "publish",
		}

		s.Stop()

		if len(p.SyncWarnings) > 0 {
			fmt.Println("Warnings on new release:")
			for _, warning := range p.SyncWarnings {
				fmt.Println(warning)
			}
		}

		err = learn.API.SendMetadataToLearn(&learn.CLIBenchmarkPayload{
			CLIBenchmark: bench,
		})
		if err != nil {
			learn.API.NotifySlack(err)
			os.Exit(1)
		}
	},
}

func currentBranch() (string, error) {
	return runBashCommand(branchCommand)
}

func remoteName() (string, error) {
	s, err := runBashCommand(pushRemoteCommand)
	if err != nil {
		return "", err
	}
	parts := strings.Split(s, ".git") // isolate url from .git
	if len(parts) < 1 {
		return "", fmt.Errorf("Error parsing remote name from %s", s)
	}
	parts = strings.Split(parts[0], "/") // select final part of the url
	if len(parts) < 1 {
		return "", fmt.Errorf("Error parsing remote name from %s", s)
	}
	return parts[len(parts)-1], nil
}

func pushToRemote(branch string) error {
	_, err := exec.Command("bash", "-c", fmt.Sprintf("git push origin %s", branch)).CombinedOutput()
	if err != nil {
		return err
	}

	return nil
}

func addAutoConfigAndCommit() error {
	_, err := exec.Command("bash", "-c", "git add autoconfig.yaml").Output()
	if err != nil {
		return err
	}
	_, err = exec.Command("bash", "-c", "git commit -m \"learn cli tool publish command: adding autoconfig.yaml\"").Output()
	if err != nil {
		return err
	}

	return nil
}

func runBashCommand(command string) (string, error) {
	out, err := exec.Command("bash", "-c", command).Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}
