package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var branchCommand = `git branch | grep \* | cut -d ' ' -f2`

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
		block, err := GetBlockByRepoName(remote)
		if err != nil {
			log.Println("Error fetchng block from learn", err)
			os.Exit(1)
		}
		// TODO if block does not exist, create one
		if block.id == 0 {
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
	return runBashCommand("git remote -v | grep push | cut -f2- -d/ | sed 's/[.].*$//'")
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

// TODO move this into a learn api package
type block struct {
	id            int
	repo_name     string
	sync_errors   []string
	title         string
	cohorts_using []int
}

type blockResponse struct {
	blocks []block `json:"blocks"`
}

func GetBlockByRepoName(repoName string) (block, error) {
	apiToken, ok := viper.Get("api_token").(string)
	if !ok {
		return block{}, errors.New("Please set your api_token in ~/.glearn-config.yaml")
	}

	u, err := url.Parse("http://localhost:3000/api/v1/blocks")
	if err != nil {
		return block{}, errors.New("unable to parse Learn remote")
	}
	v := url.Values{}
	v.Set("repo_name", repoName)
	u.RawQuery = v.Encode()

	client := &http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s", u), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error: response status: %d", res.StatusCode)
	}

	var blockResp blockResponse
	json.NewDecoder(res.Body).Decode(responseBody)

	if len(responseBody.blocks) == 1 {
		return responseBody.blocks[0]
	}
	return &block, nil
}

func CreateBlockByRepoName(repoName string) (block, error) {
	apiToken, ok := viper.Get("api_token").(string)
	if !ok {
		return nil, errors.New("Please set your api_token in ~/.glearn-config.yaml")
	}

	payload := map[string]string{
		"bucket_key_name": bucketKey,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: time.Second * 10}

	req, err := http.NewRequest("POST", "https://httpbin.org/post", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiToken))

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error: response status: %d", res.StatusCode)
	}

}
