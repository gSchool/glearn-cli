package cmd

import (
	"fmt"
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

		setupLearnAPI()

		if len(args) != 0 {
			fmt.Println("Usage: `learn publish` takes no arguments, merely pushing latest master and releasing a version to Learn. Use the command from inside a block repository.")
			os.Exit(1)
		}

		// Start benchmarking the total time spent in publish cmd
		startOfCmd := time.Now()

		repoPieces, err := remotePieces()
		if err != nil {
			fmt.Printf("Cannot run git remote detection with command: %s\n%s\n", pushRemoteCommand, err)
			os.Exit(1)
		}
		if repoPieces.RepoName == "" {
			fmt.Println("no fetch remote detected")
			os.Exit(1)
		}

		block, err := learn.API.GetBlockByRepoName(repoPieces)
		if err != nil {
			fmt.Printf("Error fetching block from learn: %s\n", err)
			os.Exit(1)
		}
		if !block.Exists() {
			block, err = learn.API.CreateBlockByRepoName(repoPieces)
			if err != nil {
				fmt.Printf("Error creating block from learn: %s\n", err)
				os.Exit(1)
			}
		}

		branch, err := currentBranch()
		if err != nil {
			fmt.Println("Cannot run git branch detection with bash:", err)
			os.Exit(1)
		}

		if branch != "master" {
			fmt.Printf("Branch publishing is cohort-specific. To continue publishing from branch '%s', go to https://learn-2.galvanize.com/cohorts/<cohortID>/setup and click the 'recycle' button for this repo.\n", branch)
			os.Exit(1)
		}

		// Detect config file
		path, _ := os.Getwd()
		createdConfig, err := publishFindOrCreateConfigDir(path + "/")
		if err != nil {
			fmt.Printf(fmt.Sprintf("Failed to find or create a config file for repo: (%s). Err: %v", branch, err))
			os.Exit(1)
		}
		fmt.Printf("Publishing block with repo name %s\n", repoPieces.RepoName)

		if createdConfig {
			fmt.Println("Committing autoconfig.yaml to", branch)
			err = addAutoConfigAndCommit()

			if err != nil && !strings.Contains(err.Error(), "Your branch is up to date with 'origin/master'.") {
				fmt.Printf("Error committing the autoconfig.yaml to origin remote on branch, run 'git rm autoconfig.yaml' to remove it from reference then add a new commit: %s", err)
				os.Exit(1)
			}
		}

		fmt.Println("Pushing work to remote origin", branch)

		// TODO what happens when they do not have work in remote and push fails?
		err = pushToRemote(branch)
		if err != nil {
			fmt.Printf("\nError pushing to origin remote on branch:\n\n%s", err)
			os.Exit(1)
		}

		// Start benchmark for creating master release & building on learn
		startOfMasterReleaseAndBuild := time.Now()

		// Start a processing spinner that runs until Learn is finsihed building the preview
		fmt.Println("\nBuilding release...")
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

		var attempts uint8 = 30
		p, err := learn.API.PollForBuildResponse(releaseID, &attempts)
		if err != nil {
			s.Stop()

			block, err := learn.API.GetBlockByRepoName(repoPieces)
			if err != nil {
				fmt.Printf("Error fetching block from learn: %s\n", err)
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

func remotePieces() (learn.RepoPieces, error) {
	var repoPieces learn.RepoPieces
	s, err := runBashCommand(pushRemoteCommand)
	if err != nil {
		return repoPieces, err
	}
	parts := strings.Split(s, ".git")
	if len(parts) < 1 { // There should only be 1
		return repoPieces, fmt.Errorf("Error parsing git remote from %s", s)
	}
	parts = strings.Split(parts[0], "/")

	// does it start with https
	if parts[0] == "https:" || parts[0] == "ssh:" {
		repoPieces.Origin = strings.ReplaceAll(parts[2], "git@", "")
		repoPieces.Org = parts[3]
		repoPieces.RepoName = parts[4]

		return repoPieces, nil
	}

	repoPieces.RepoName = parts[1]
	parts = strings.Split(parts[0], ":")
	repoPieces.Org = parts[1]
	parts = strings.Split(parts[0], "@")
	repoPieces.Origin = parts[1]

	return repoPieces, nil
}

func pushToRemote(branch string) error {
	out, err := exec.Command("bash", "-c", fmt.Sprintf("git push origin %s", branch)).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", out)
	}

	return nil
}

func addAutoConfigAndCommit() error {
	top, err := GitTopLevelDir()
	addCmd := "git add " + strings.TrimSpace(top) + "/autoconfig.yaml"
	out, err := exec.Command("bash", "-c", addCmd).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", out)
	}
	out, err = exec.Command("bash", "-c", "git commit -m \"learn cli tool publish command: adding autoconfig.yaml\"").CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", out)
	}

	return nil
}

func runBashCommand(command string) (string, error) {
	out, err := exec.Command("bash", "-c", command).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s- %s", out, err)
	}

	return strings.TrimSpace(string(out)), nil
}
