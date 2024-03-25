package cmd

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

const ignoredMDContent = `## Ignored

This file should be ignored.`

const testMDContent = `## Test links

![alt](./image/nested-small.png)
![alt](image/nested-small.png)
![alt](../nested-small.png)`

const dockerMDContent = `### !challenge

* type: custom-snippet
* language: text
* id: 8c406f4f-6428-498b-be24-6bd0a6c9096b
* title: Title
* docker_directory_path: /path/to/dir

##### !question

Question

##### !end-question

##### !placeholder

##### !end-placeholder

### !end-challenge
`

func Test_ParseConfigFileForPaths(t *testing.T) {
	previewFindOrCreateConfig(withNoConfigFixture, false, []string{})
	p := previewBuilder{target: withNoConfigFixture}
	err := p.parseConfigAndGatherPaths()

	if err != nil || len(p.configYamlPaths) == 0 {
		t.Errorf("Should of parse the yaml and gathered some content file paths")
	}
}

func Test_compressDirectory(t *testing.T) {
	source := "../../fixtures/test-block-auto-config"
	p := previewBuilder{target: source}
	err := p.parseConfigAndGatherPaths()
	if err != nil {
		t.Errorf("Attempting to parseConfigAndGatherLinkedPaths errored: %s\n", err)
	}
	if len(p.configYamlPaths) < 1 {
		t.Errorf("There should be paths parsed from the target")
	}

	err = createTestFile("../../fixtures/test-block-auto-config/ignored.md", ignoredMDContent)
	if err != nil {
		t.Error("Could not make ignored.md")
	}

	err = os.MkdirAll("../../fixtures/test-block-auto-config/ignored/", os.FileMode(0777))
	if err != nil {
		t.Error("Could not make ignored/")
	}

	err = createTestFile("../../fixtures/test-block-auto-config/ignored/ignored.md", ignoredMDContent)
	if err != nil {
		t.Error("Could not make ignored/ignored.md")
	}

	tmpZipFile := "../../fixtures/test-block-auto-config/preview-curriculum.zip"

	var challengePaths []string
	challengePaths = append(challengePaths, "test-block-auto-config/docker/text.text")
	challengePaths = append(challengePaths, "test-block-auto-config/sql/database.sql")

	previewer := previewBuilder{
		target:          source,
		challengePaths:  challengePaths,
		configYamlPaths: p.configYamlPaths,
	}
	err = previewer.compressDirectory(tmpZipFile)
	if err != nil {
		t.Errorf("compressDirectory failed to do its job: %s\n", err)
	}

	read, _ := zip.OpenReader(tmpZipFile)
	defer read.Close()

	var paths = make(map[string]bool)
	for _, file := range read.File {
		if strings.HasSuffix(file.Name, "/") == false {
			paths[file.Name] = false
		}
	}

	for path := range paths {
		for _, includedPath := range p.configYamlPaths {
			if strings.Contains(includedPath, path) {
				paths[path] = true
			}
		}
		for _, includedPath := range challengePaths {
			if strings.Contains(includedPath, path) {
				paths[path] = true
			}
		}
		if strings.Contains(path, "autoconfig.yaml") {
			paths[path] = true
		}
	}

	for path, found := range paths {
		if found == false {
			t.Errorf("Should of found: %s In zipped dir", path)
		}
	}

	fmt.Printf("%+v", paths)

	if _, ok := paths["ignored.md"]; ok {
		t.Error("ZIP should not contain ignored.md")
	}

	if _, ok := paths["ignored/ignored.md"]; ok {
		t.Error("ZIP should not contain ignored/ignored.md")
	}

	os.Remove(tmpZipFile)
}

func Test_createNewTarget(t *testing.T) {
	result, err := createNewTarget("../../fixtures/test-links/nested/test.md", []string{}, []string{"./mrsmall-invert.png", "../mrsmall.png", "../image/nested-small.png", "deeper/deep-small.png"}, []string{})
	if err != nil {
		t.Errorf("Attempting to createNewTarget errored: %s\n", err)
	}
	if result != "single-file-upload" {
		t.Errorf("result should be the temp directory with the target markdown, '%s'", result)
	}

	testFilesExist(t, []string{"test.md", "mrsmall-invert.png", "deeper/deep-small.png"})

	// test cases for ,./
	if _, err = os.Stat(fmt.Sprintf("single-file-upload/%s", "mrsmall.png")); os.IsNotExist(err) {
		t.Errorf("mrsmall.png should have been created and moved to the root of the single file directory, was not")
	}
	if _, err = os.Stat(fmt.Sprintf("single-file-upload/image/%s", "nested-small.png")); os.IsNotExist(err) {
		t.Errorf("nested-small.png should have been created and it's image dir moved to the root of the single file directory, was not")
	}

	b, err := ioutil.ReadFile("single-file-upload/test.md")
	if err != nil {
		t.Errorf("could not Open test.md")
	}
	if strings.Contains(string(b), "../mrsmall.png") {
		t.Errorf("test.md file should not contain '../mrsmall.png' but does:\n%s\n", string(b))
	}
	if strings.Contains(string(b), "../image/nested-small.png") {
		t.Errorf("test.md file should not contain '../image/nested-small.png' but does:\n%s\n", string(b))
	}
	err = os.RemoveAll("single-file-upload")
	if err != nil {
		t.Errorf("could not remove single-file-upload directory: %s\n", err)
	}
}

func Test_createNewTargetChallengePathsAndLinks(t *testing.T) {
	createTestMD(testMDContent)
	os.MkdirAll("image", os.FileMode(0777))
	_, err := os.Create("image/nested-small.png")
	if err != nil {
		t.Errorf("Error generating test image: %s\n", err)
	}

	output := captureOutput(func() {
		createNewTarget("test.md", []string{"/data/some.sql"}, []string{"image/nested-small.png"}, []string{})
		_, err := os.Stat(fmt.Sprintf("single-file-upload/%s", "data/some.sql"))
		if err == nil {
			t.Errorf("data/some.sql should not have been copied over")
		}

		if _, err = os.Stat(fmt.Sprintf("single-file-upload/%s", "/image/nested-small.png")); os.IsNotExist(err) {
			t.Errorf("mrsmall-invert should have been created, was not")
		}
	})

	if strings.Contains(output, "Link not found with path") {
		t.Errorf("output should not print 'Link not found with path', output was:\n%s\n", output)
	}

	if strings.Contains(output, "Failed build tmp files around single file preview for") {
		t.Errorf("output should not print 'Failed build tmp files around single file preview for', output was:\n%s\n", output)
	}

	err = os.RemoveAll("single-file-upload")
	if err != nil {
		t.Errorf("could not remove single-file-upload directory: %s\n", err)
	}

	err = os.RemoveAll("image")
	if err != nil {
		t.Errorf("could not remove image directory: %s\n", err)
	}

	err = os.Remove("test.md")
	if err != nil {
		t.Errorf("could not remove 'test.md' file: %s\n", err)
	}
}

const allContent = ` ## all content
![alt](image/nested-small.png)
* docker_directory_path: /path/to/dir
* test_file: /tests/test.js
* setup_file: /setup.js
`

func Test_createNewTargetAllAssets(t *testing.T) {
	createTestMD(allContent)
	// link
	os.MkdirAll("image", os.FileMode(0777))
	_, err := os.Create("image/nested-small.png")
	if err != nil {
		t.Errorf("Error generating test image: %s\n", err)
	}

	// docker
	os.MkdirAll("path/to/dir", os.FileMode(0777))
	_, err = os.Create("path/to/dir/Dockerfile")
	if err != nil {
		t.Errorf("Error generating test fixtures: %s\n", err)
	}
	_, err = os.Create("path/to/dir/test.sh")
	if err != nil {
		t.Errorf("Error generating test fixtures: %s\n", err)
	}

	// challenge
	os.MkdirAll("tests", os.FileMode(0777))
	_, err = os.Create("tests/test.js")
	if err != nil {
		t.Errorf("Error generating test file tests/test.js: %s\n", err)
	}
	_, err = os.Create("setup.js")
	if err != nil {
		t.Errorf("Error generating test file setup.js: %s\n", err)
	}

	output := captureOutput(func() {
		createNewTarget("test.md", []string{"/tests/test.js", "/setup.js"}, []string{"image/nested-small.png"}, []string{"/path/to/dir"})
		testFilesExist(t, []string{"/tests/test.js", "/setup.js", "image/nested-small.png", "/path/to/dir"})
	})

	if strings.Contains(output, "Link not found with path") {
		t.Errorf("output should not print 'Link not found with path', output was:\n%s\n", output)
	}

	if strings.Contains(output, "Failed build tmp files around single file preview for") {
		t.Errorf("output should not print 'Failed build tmp files around single file preview for', output was:\n%s\n", output)
	}

	err = os.RemoveAll("single-file-upload")
	if err != nil {
		t.Errorf("could not remove single-file-upload directory: %s\n", err)
	}

	err = os.RemoveAll("image")
	if err != nil {
		t.Errorf("could not remove 'image' directory: %s\n", err)
	}

	err = os.RemoveAll("tests")
	if err != nil {
		t.Errorf("could not remove 'tests' directory: %s\n", err)
	}

	err = os.RemoveAll("path")
	if err != nil {
		t.Errorf("could not remove docker directory 'path': %s\n", err)
	}

	err = os.Remove("setup.js")
	if err != nil {
		t.Errorf("could not remove 'setup.js.md' file: %s\n", err)
	}

	err = os.Remove("test.md")
	if err != nil {
		t.Errorf("could not remove 'test.md' file: %s\n", err)
	}
}

func Test_createNewTargetSingleFileThatIsSQL(t *testing.T) {
	err := createTestMD(testMDContent)
	if err != nil {
		t.Errorf("Error creating test.md: %s\n", err)
	}
	output := captureOutput(func() {
		createNewTarget("test.md", []string{"/data/some.sql"}, []string{}, []string{})
		_, err := os.Stat(fmt.Sprintf("single-file-upload/%s", "data/some.sql"))
		if err == nil {
			t.Errorf("data/some.sql should have been copied over and it was not")
		}
	})

	if strings.Contains(output, "Link not found with path") {
		t.Errorf("output should not print 'Link not found with path', output was:\n%s\n", output)
	}

	err = os.Remove("test.md")
	if err != nil {
		t.Errorf("could not remove 'test.md' file: %s\n", err)
	}
}

func Test_createNewTargetSingleFile(t *testing.T) {
	createTestMD(testMDContent)
	os.MkdirAll("image", os.FileMode(0777))
	os.Create("image/nested-small.png")
	_, err := os.Create("../nested-small.png")
	if err != nil {
		t.Errorf("Error generating test fixtures: %s\n", err)
	}

	output := captureOutput(func() {
		result, err := createNewTarget("test.md", []string{}, []string{"./image/nested-small.png", "image/nested-small.png", "../nested-small.png"}, []string{})
		if err != nil {
			t.Errorf("Attempting to createNewTarget errored: %s\n", err)
		}
		if result != "single-file-upload" {
			t.Errorf("result should be the temp directory with the target markdown, '%s'", result)
		}

		testFilesExist(t, []string{"test.md", "image/nested-small.png", "nested-small.png"})

		err = os.RemoveAll("single-file-upload")
		if err != nil {
			t.Errorf("could not remove single-file-upload directory: %s\n", err)
		}
	})

	if strings.Contains(output, "Link not found with path") {
		t.Errorf("output should not print 'Link not found with path', output was:\n%s\n", output)
	}

	err = os.RemoveAll("image")
	if err != nil {
		t.Errorf("could not remove 'image' directory: %s\n", err)
	}

	err = os.Remove("../nested-small.png")
	if err != nil {
		t.Errorf("could not remove '../nested-small.png' file: %s\n", err)
	}

	err = os.Remove("test.md")
	if err != nil {
		t.Errorf("could not remove 'test.md' file: %s\n", err)
	}
}

func Test_createNewTarget_DockerDirectoryIgnore(t *testing.T) {
	createTestMD(dockerMDContent)
	err := os.MkdirAll("path/to/dir/child", os.FileMode(0777))
	if err != nil {
		t.Errorf("Error generating test fixtures: %s\n", err)
	}
	_, err = os.Create("path/to/dir/root.png")
	if err != nil {
		t.Errorf("Error generating test fixtures: %s\n", err)
	}
	_, err = os.Create("path/to/dir/child/nest.png")
	if err != nil {
		t.Errorf("Error generating test fixtures: %s\n", err)
	}
	_, err = os.Create("path/to/dir/ignore_me.png")
	if err != nil {
		t.Errorf("Error generating test fixtures: %s\n", err)
	}
	// Always allow files
	_, err = os.Create("path/to/dir/Dockerfile")
	if err != nil {
		t.Errorf("Error generating test fixtures: %s\n", err)
	}
	_, err = os.Create("path/to/dir/test.sh")
	if err != nil {
		t.Errorf("Error generating test fixtures: %s\n", err)
	}
	_, err = os.Create("path/to/dir/docker-compose.yaml")
	if err != nil {
		t.Errorf("Error generating test fixtures: %s\n", err)
	}
	_, err = os.Create("path/to/dir/docker-compose.yml")
	if err != nil {
		t.Errorf("Error generating test fixtures: %s\n", err)
	}

	_, err = os.Create("path/to/dir/child/ignore_me.png")
	if err != nil {
		t.Errorf("Error generating test fixtures: %s\n", err)
	}
	_, err = os.Create("path/to/dir/child/dont_agnore_me.png")
	if err != nil {
		t.Errorf("Error generating test fixtures: %s\n", err)
	}

	// DockerIgnore stuff
	ignoreFile, err := os.Create("path/to/dir/.dockerignore")
	if err != nil {
		t.Errorf("Error generating test fixtures: %s\n", err)
	}
	defer ignoreFile.Close()
	ignoreFile.Write([]byte("ignore_me.png"))
	ignoreFile.Write([]byte("Dockerfile"))
	ignoreFile.Write([]byte("test.sh"))
	ignoreFile.Write([]byte("docker-compose.yaml"))
	ignoreFile.Write([]byte("docker-compose.yml"))

	output := captureOutput(func() {
		result, err := createNewTarget("test.md", []string{}, []string{}, []string{"/path/to/dir"})
		if err != nil {
			t.Errorf("Attempting to createNewTarget errored: %s\n", err)
		}

		if result != "single-file-upload" {
			t.Errorf("result should be the temp directory with the target markdown, '%s'", result)
		}

		testFilesExist(t, []string{
			"test.md",
			"path/to/dir/Dockerfile",
			"path/to/dir/test.sh",
			"path/to/dir/docker-compose.yaml",
			"path/to/dir/docker-compose.yml",
			"path/to/dir/child/dont_agnore_me.png",
			"path/to/dir/child/nest.png",
		})

		if _, err = os.Stat(fmt.Sprintf("single-file-upload/%s", "path/to/dir/ignore_me.jpg")); !os.IsNotExist(err) {
			t.Errorf("path/to/dir/ignore_me.jpg should NOT have been created because its in the docker ignore file")
		}

		if _, err = os.Stat(fmt.Sprintf("single-file-upload/%s", "path/to/dir/child/ignore_me.jpg")); !os.IsNotExist(err) {
			t.Errorf("path/to/dir/child/ignore_me.jpg should NOT have been created because its in the docker ignore file")
		}

		err = os.RemoveAll("single-file-upload")
		if err != nil {
			t.Errorf("could not remove single-file-upload directory: %s\n", err)
		}
	})

	err = os.RemoveAll("path")
	if err != nil {
		t.Errorf("could not remove 'path' directory: %s\n", err)
	}

	err = os.Remove("test.md")
	if err != nil {
		t.Errorf("could not remove 'test.md' file: %s\n", err)
	}

	if output != "" {
		t.Error("expected stdout exist")
	}
}

func Test_createNewTarget_DockerDirectoryNestedMd(t *testing.T) {
	createTestMD(dockerMDContent)
	err := os.MkdirAll("../path/to/dir/child", os.FileMode(0777))
	if err != nil {
		t.Errorf("Error generating test fixtures: %s\n", err)
	}
	_, err = os.Create("../path/to/dir/root.png")
	if err != nil {
		t.Errorf("Error generating test fixtures: %s\n", err)
	}
	_, err = os.Create("../path/to/dir/child/nest.png")
	if err != nil {
		t.Errorf("Error generating test fixtures: %s\n", err)
	}

	output := captureOutput(func() {
		result, err := createNewTarget("test.md", []string{}, []string{}, []string{"/path/to/dir"})
		if err != nil {
			t.Errorf("Attempting to createNewTarget errored: %s\n", err)
		}

		if result != "single-file-upload" {
			t.Errorf("result should be the temp directory with the target markdown, '%s'", result)
		}

		testFilesExist(t, []string{"test.md", "path/to/dir/root.png", "path/to/dir/child/nest.png"})

		err = os.RemoveAll("single-file-upload")
		if err != nil {
			t.Errorf("could not remove single-file-upload directory: %s\n", err)
		}
	})

	err = os.RemoveAll("../path")
	if err != nil {
		t.Errorf("could not remove 'path' directory: %s\n", err)
	}

	err = os.Remove("test.md")
	if err != nil {
		t.Errorf("could not remove 'test.md' file: %s\n", err)
	}

	if output != "" {
		t.Error("expected stdout exist")
	}
}

func Test_createNewTarget_DockerDirectoryDoubleNestedMd(t *testing.T) {
	createTestMD(dockerMDContent)
	err := os.MkdirAll("../../path/to/dir/child", os.FileMode(0777))
	if err != nil {
		t.Errorf("Error generating test fixtures: %s\n", err)
	}
	_, err = os.Create("../../path/to/dir/root.png")
	if err != nil {
		t.Errorf("Error generating test fixtures: %s\n", err)
	}
	_, err = os.Create("../../path/to/dir/child/nest.png")
	if err != nil {
		t.Errorf("Error generating test fixtures: %s\n", err)
	}

	output := captureOutput(func() {
		result, err := createNewTarget("test.md", []string{}, []string{}, []string{"/path/to/dir"})
		if err != nil {
			t.Errorf("Attempting to createNewTarget errored: %s\n", err)
		}

		if result != "single-file-upload" {
			t.Errorf("result should be the temp directory with the target markdown, '%s'", result)
		}

		testFilesExist(t, []string{"test.md", "path/to/dir/root.png", "path/to/dir/child/nest.png"})

		err = os.RemoveAll("single-file-upload")
		if err != nil {
			t.Errorf("could not remove single-file-upload directory: %s\n", err)
		}
	})

	err = os.RemoveAll("../../path")
	if err != nil {
		t.Errorf("could not remove 'path' directory: %s\n", err)
	}

	err = os.Remove("test.md")
	if err != nil {
		t.Errorf("could not remove 'test.md' file: %s\n", err)
	}

	if output != "" {
		t.Error("expected stdout exist")
	}
}

func testFilesExist(t *testing.T, paths []string) {
	for _, file := range paths {
		if _, err := os.Stat(fmt.Sprintf("single-file-upload/%s", file)); os.IsNotExist(err) {
			t.Errorf("%s should have been created\n", file)
		}
	}
}

func createTestMD(content string) error {
	return createTestFile("test.md", content)
}

func createTestFile(fileName string, content string) error {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if _, err := f.Write([]byte(content)); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

func captureOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}
