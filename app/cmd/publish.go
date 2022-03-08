package cmd

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"regexp"

	"github.com/briandowns/spinner"
	"github.com/gSchool/glearn-cli/api/learn"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	branchCommand     = `git rev-parse --abbrev-ref HEAD`
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

		if !IgnoreLocal {
			notCurrentWithRemote := notCurrentWithRemote(branch)
			if notCurrentWithRemote {
				fmt.Println("\nWARNING:")
				fmt.Println("You have local changes that are not on remote, run `git status` for details.")
				fmt.Println("\nPublishing from current remote")
			}
		}

		// Detect config file
		path, _ := os.Getwd()
		createdConfig, err := publishFindOrCreateConfig(path + "/")
		if err != nil {
			fmt.Printf("%s", fmt.Sprintf("failed to find or create a config file for repo: (%s). Err: %v", branch, err))
			os.Exit(1)
		}
		fmt.Printf("Publishing block with repo name %s from branch %s\n", repoPieces.RepoName, branch)

		if createdConfig && CiCdEnvironment {
			fmt.Println("\nError: You cannot use autoconfig.yaml from a CI/CD environment.")
			fmt.Println("Please create a config.yaml file and commit it.")
			os.Exit(1)
		} else if createdConfig {
			fmt.Println("Committing autoconfig.yaml to", branch)
			err = addAutoConfigAndCommit()

			if err != nil && !strings.Contains(err.Error(), fmt.Sprintf("Your branch is up to date with 'origin/%s'.", branch)) {
				fmt.Printf("Error committing the autoconfig.yaml to origin remote on branch, run 'git rm autoconfig.yaml' to remove it from reference then add a new commit: %s", err)
				os.Exit(1)
			}
		}

		// Do not push if in a CI/CD environment
		if !CiCdEnvironment {
			fmt.Println("Pushing work to remote origin", branch)

			err = pushToRemote(branch)
			if err != nil {
				fmt.Printf("\nError pushing to origin remote on branch:\n\n%s", err)
				os.Exit(1)
			}
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

		blockUrl := fmt.Sprintf("%s/blocks/%d?branch_name=%s", learn.API.BaseURL(), block.ID, url.QueryEscape(branch))
		fmt.Printf("Block released! %s\n", blockUrl)

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

var hostedGitRe = regexp.MustCompile(`^(:?(\w+):\/\/\/?)?(?:(~?\w+)@)?([\w\d\.\-_]+)(:?:([\d]+))?(?::)?\/*(.*)\/([\w\d\.\-_]+)(?:\.git)\/?$`)

func parseHostedGitRegex(originUrl string) (repoParts learn.RepoPieces, err error) {
	var repoPieces learn.RepoPieces
	if m := hostedGitRe.FindStringSubmatch(originUrl); m != nil {
		repoPieces, err = learn.RepoPieces{
			Origin: m[4],
			Org: m[7],
			RepoName: m[8],
		}, nil
	}
	return repoPieces, err
}

func remotePieces() (learn.RepoPieces, error) {
	var repoPieces learn.RepoPieces
	s, err := runBashCommand(pushRemoteCommand)
	if err != nil {
		return repoPieces, err
	}
	return parseHostedGitRegex(s)
}

func pushToRemote(branch string) error {
	out, err := exec.Command("bash", "-c", fmt.Sprintf("git push origin %s", branch)).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", out)
	}

	return nil
}

func addAutoConfigAndCommit() error {
	top, _ := GitTopLevelDir()
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
