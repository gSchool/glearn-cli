package cmd

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const guideDir = "learn-curriculum-walkthrough"

//go:embed embeds/walkthrough/README.md
var readme []byte

//go:embed embeds/walkthrough/01-example-unit/00-hello-world.md
var helloWorldMd []byte

//go:embed embeds/walkthrough/01-example-unit/01-configuration.md
var configurationMd []byte

//go:embed embeds/walkthrough/01-example-unit/02-publishing.md
var publishingMd []byte

//go:embed embeds/walkthrough/01-example-unit/03-markdown-examples.md
var markdownExamplesMd []byte

//go:embed embeds/walkthrough/01-example-unit/04-challenges.md
var challengesMd []byte

//go:embed embeds/walkthrough/01-example-unit/05-checkpoint.md
var checkpointMd []byte

//go:embed embeds/walkthrough/01-example-unit/description.yaml
var descriptionYml []byte

//go:embed embeds/walkthrough/01-example-unit/images/github.jpg
var githubJpg []byte

//go:embed embeds/walkthrough/01-example-unit/images/kmeans.png
var kmeansPng []byte

//go:embed embeds/walkthrough/01-example-unit/images/react.png
var reactPng []byte

//go:embed embeds/walkthrough/01-example-unit/sql-files/foodtruck.sql
var foodtruckSql []byte

//go:embed embeds/walkthrough/01-example-unit/custom-snippets/hello-world/Dockerfile
var Dockerfile []byte

//go:embed embeds/walkthrough/01-example-unit/custom-snippets/hello-world/submission.txt
var submissionTxt []byte

//go:embed embeds/walkthrough/01-example-unit/custom-snippets/hello-world/test.sh
var testSh []byte

type guideFile struct {
	path    string
	content []byte
}

var guideCmd = &cobra.Command{
	Use:     "walkthrough",
	Aliases: []string{"guide"},
	Short:   "Generate examples for use in the walkthrough",
	Long:    "Generate examples for use in the walkthrough",
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println("Could not detect a working directory")
			os.Exit(1)
		}

		// Does that directory have a config file
		hasConfig, _ := doesCurrentDirHaveConfig(currentDir)
		if hasConfig {
			fmt.Println("WARNING: configuration file detected and cannot continue with `learn walkthrough` command.")
			os.Exit(1)
		}
		_, dirExists := os.Stat("/" + guideDir)
		if dirExists == nil {
			fmt.Printf("A directory already exists by the name '%s', rename or move it.\n", guideDir)
			os.Exit(1)
		}

		fmt.Printf("\nWriting '%s' directory and contents...\n", guideDir)

		// Create contents in the directory
		err = generateGuide(currentDir)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(`
Success!
========

Open kj`)
		fmt.Printf("\nTo get started run 'cd %s && learn preview 01-example-unit/00-hello-world.md' and follow the instructions to find your content.\n\n", guideDir)
	},
}

func doesCurrentDirHaveConfig(currentDir string) (bool, bool) {
	configExist := false
	autoConfigExist := false

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
		autoConfigExist = true
	}

	return configExist, autoConfigExist
}

func generateGuide(currentDir string) error {
	guideFiles := []guideFile{
		{"README.md", readme},
		{"01-example-unit/00-hello-world.md", helloWorldMd},
		{"01-example-unit/01-configuration.md", configurationMd},
		{"01-example-unit/02-publishng.md", publishingMd},
		{"01-example-unit/03-markdown-examples.md", markdownExamplesMd},
		{"01-example-unit/04-challenges.md", challengesMd},
		{"01-example-unit/05-checkpoint.md", checkpointMd},
		{"01-example-unit/description.yaml", descriptionYml},
		{"01-example-unit/images/github.jpg", githubJpg},
		{"01-example-unit/images/kmeans.png", kmeansPng},
		{"01-example-unit/images/react.png", reactPng},
		{"01-example-unit/sql-files/foodtruck.sql", foodtruckSql},
		{"01-example-unit/custom-snippets/hello-world/Dockerfile", Dockerfile},
		{"01-example-unit/custom-snippets/hello-world/submission.txt", submissionTxt},
		{"01-example-unit/custom-snippets/hello-world/test.sh", testSh},
	}

	os.MkdirAll(guideDir, os.FileMode(0777))
	os.MkdirAll(guideDir+"/01-example-unit", os.FileMode(0777))
	os.MkdirAll(guideDir+"/01-example-unit/images", os.FileMode(0777))
	os.MkdirAll(guideDir+"/01-example-unit/sql-files", os.FileMode(0777))
	os.MkdirAll(guideDir+"/01-example-unit/custom-snippets", os.FileMode(0777))
	os.MkdirAll(guideDir+"/01-example-unit/custom-snippets/hello-world", os.FileMode(0777))

	for _, file := range guideFiles {
		location := fmt.Sprintf("./%s/%s", guideDir, file.path)
		err := os.WriteFile(location, file.content, 0677)
		if err != nil {
			return fmt.Errorf("Error writing guide contents '%s': %v\n", file.path, err)
		}
	}

	return nil
}
