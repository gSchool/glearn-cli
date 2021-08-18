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
The Learn system recognizes blocks of content held in GitHub repositories. This
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

		if IgnoreLocal == false {
			notCurrentWithRemote := notCurrentWithRemote(branch)
			if notCurrentWithRemote {
				fmt.Println("\nWARNING:")
				fmt.Println("You have local changes that are not on remote, run `git status` for details.")
				fmt.Println("\nPublishing from current remote")
			}
		}

		// Detect config file
		path, _ := os.Getwd()
		createdConfig, err := publishFindOrCreateConfigDir(path + "/")
		if err != nil {
			fmt.Printf(fmt.Sprintf("Failed to find or create a config file for repo: (%s). Err: %v", branch, err))
			os.Exit(1)
		}
		fmt.Printf("Publishing block with repo name %s from branch %s\n", repoPieces.RepoName, branch)

		if createdConfig {
			fmt.Println("Committing autoconfig.yaml to", branch)
			err = addAutoConfigAndCommit()

			if err != nil && !strings.Contains(err.Error(), fmt.Sprintf("Your branch is up to date with 'origin/%s'.", branch)) {
				fmt.Printf("Error committing the autoconfig.yaml to origin remote on branch, run 'git rm autoconfig.yaml' to remove it from reference then add a new commit: %s", err)
				os.Exit(1)
			}
		}

		fmt.Println("Pushing work to remote origin", branch)

		err = pushToRemote(branch)
		if err != nil {
			fmt.Printf("\nError pushing to origin remote on branch:\n\n%s", err)
			os.Exit(1)
		}

		// Start benchmark for creating master release & building on learn
		startOfMasterReleaseAndBuild := time.Now()

		// Start a processing spinner that runs until Learn is finished building the preview
		fmt.Println("\nBuilding release...")
		s := spinner.New(spinner.CharSets[32], 100*time.Millisecond)
		s.Color("green")
		s.Start()

		// Create a release on learn, notify user
		releaseID, err := learn.API.CreateBranchRelease(block.ID, branch)
		if err != nil || releaseID == 0 {
			fmt.Printf("Release failed. releaseID: %d. Error: %s\n", releaseID, err)
			os.Exit(1)
		}

		var attempts uint8 = 30
		p, err := learn.API.PollForBuildResponse(releaseID, false, "", &attempts)
		if err != nil {
			s.Stop()

			if p != nil && p.Errors != "" {
				fmt.Printf("Release failed: %s\n", p.Errors)
				os.Exit(1)
			}

			if p != nil && len(p.SyncWarnings) > 0 {
				fmt.Printf("Release warnings:")

				for _, sw := range p.SyncWarnings {
					fmt.Println(sw)
				}
			}

			block, err := learn.API.GetBlockByRepoName(repoPieces)
			if err != nil {
				fmt.Printf("Release failed. Error fetching block from learn: %s\n", err)
				os.Exit(1)
			}
			if len(block.SyncErrors) > 0 {
				fmt.Println("Release failed. Errors on block:")
				for _, e := range block.SyncErrors {
					fmt.Println(e)
				}
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

		fmt.Printf("Block released! %s/blocks/%d?branch_name=%s\n", learn.API.BaseURL(), block.ID, branch)

		if len(p.SyncWarnings) > 0 {
			fmt.Println("\nWarnings on new release:")
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

func notCurrentWithRemote(branch string) bool {
	out, err := runBashCommand("git status")
	if err != nil {
		return false
	}

	if strings.Contains(out, "Changes not staged for commit:") || strings.Contains(out, "Changes to be committed:") {
		return true
	}
	// look up the remote branches and their push state
	remoteOut, err := runBashCommand("git remote show origin")
	if err != nil {
		return false
	}
	// Get to the section which defines local refs configured for git push
	// read the lines until we find one which starts with the branch name
	// If it contains (up to date) then we would be in the clear to publish
	var afterPushRefs bool
	for _, line := range strings.Split(remoteOut, "\n") {
		if afterPushRefs {
			trimLine := strings.TrimSpace(line)
			// Lines we are concerned with look like this:
			//   main        pushes to main        (up to date)
			if strings.HasPrefix(trimLine, branch+" ") { // branch names can't have whitespace, and the name now starts and ends in whitespace
				if strings.Contains(trimLine, "(up to date)") {
					return false
				} else {
					return true
				}
			}
		}
		if strings.Contains(line, "configured for 'git push'") && !afterPushRefs {
			afterPushRefs = true
		}
	}

	return true
}
