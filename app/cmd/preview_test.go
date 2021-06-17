package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

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

const dockerIgnore = `ignore_me.jpg
*.txt
`

func Test_createNewTarget(t *testing.T) {
	result, err := createNewTarget("../../fixtures/test-links/nested/test.md", []string{"./mrsmall-invert.png", "../mrsmall.png", "../image/nested-small.png", "deeper/deep-small.png"}, []string{})
	if err != nil {
		t.Errorf("Attempting to createNewTarget errored: %s\n", err)
	}
	if result != "single-file-upload" {
		t.Errorf("result should be the temp directory with the target markdown, '%s'", result)
	}

	if _, err := os.Stat(fmt.Sprintf("single-file-upload/%s", "test.md")); os.IsNotExist(err) {
		t.Errorf("test.md should have been created")
	}
	if _, err = os.Stat(fmt.Sprintf("single-file-upload/%s", "mrsmall-invert.png")); os.IsNotExist(err) {
		t.Errorf("mrsmall-invert should have been created, was not")
	}
	if _, err = os.Stat(fmt.Sprintf("single-file-upload/deeper/%s", "deep-small.png")); os.IsNotExist(err) {
		t.Errorf("deeper/deep-small.png should have been created in its directory, was not")
	}

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

func Test_createNewTargetSingleFileSQLWithImage(t *testing.T) {
	err := createTestMD(testMDContent)
	err = os.MkdirAll("image", os.FileMode(0777))
	_, err = os.Create("image/nested-small.png")
	if err != nil {
		t.Errorf("Error generating test image: %s\n", err)
	}

	output := captureOutput(func() {
		_, err := createNewTarget("test.md", []string{"/data/some.sql", "image/nested-small.png"}, []string{})
		_, err = os.Stat(fmt.Sprintf("single-file-upload/%s", "data/some.sql"))
		if err == nil {
			t.Errorf("data/some.sql should have been copied over and it was not")
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

func Test_createNewTargetSingleFileThatIsSQL(t *testing.T) {
	err := createTestMD(testMDContent)
	if err != nil {
		t.Errorf("Error creating test.md: %s\n", err)
	}
	output := captureOutput(func() {
		_, err := createNewTarget("test.md", []string{"/data/some.sql"}, []string{})
		_, err = os.Stat(fmt.Sprintf("single-file-upload/%s", "data/some.sql"))
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
	err := createTestMD(testMDContent)
	err = os.MkdirAll("image", os.FileMode(0777))
	_, err = os.Create("image/nested-small.png")
	_, err = os.Create("../nested-small.png")
	if err != nil {
		t.Errorf("Error generating test fixtures: %s\n", err)
	}

	output := captureOutput(func() {
		result, err := createNewTarget("test.md", []string{"./image/nested-small.png", "image/nested-small.png", "../nested-small.png"}, []string{})
		if err != nil {
			t.Errorf("Attempting to createNewTarget errored: %s\n", err)
		}
		if result != "single-file-upload" {
			t.Errorf("result should be the temp directory with the target markdown, '%s'", result)
		}

		if _, err := os.Stat(fmt.Sprintf("single-file-upload/%s", "test.md")); os.IsNotExist(err) {
			t.Errorf("test.md should have been created")
		}

		if _, err = os.Stat(fmt.Sprintf("single-file-upload/image/%s", "nested-small.png")); os.IsNotExist(err) {
			t.Errorf("nested-small.png should have been created and it's image dir moved to the root of the single file directory, was not")
		}

		if _, err = os.Stat(fmt.Sprintf("single-file-upload/%s", "nested-small.png")); os.IsNotExist(err) {
			t.Errorf("nested-small.png should have been created and it's image dir moved to the root of the single file directory, was not")
		}

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
	err := createTestMD(dockerMDContent)
	err = os.MkdirAll("path/to/dir/child", os.FileMode(0777))
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

	output := captureOutput(func() {
		result, err := createNewTarget("test.md", []string{}, []string{"/path/to/dir"})
		if err != nil {
			t.Errorf("Attempting to createNewTarget errored: %s\n", err)
		}

		if result != "single-file-upload" {
			t.Errorf("result should be the temp directory with the target markdown, '%s'", result)
		}

		if _, err := os.Stat(fmt.Sprintf("single-file-upload/%s", "test.md")); os.IsNotExist(err) {
			t.Errorf("test.md should have been created")
		}

		if _, err := os.Stat(fmt.Sprintf("single-file-upload/path/to/dir/child/%s", "dont_agnore_me.png")); os.IsNotExist(err) {
			t.Errorf("dont_agnore_me.png should have been created")
		}

		if _, err = os.Stat(fmt.Sprintf("single-file-upload/%s", "path/to/dir/ignore_me.jpg")); !os.IsNotExist(err) {
			t.Errorf("path/to/dir/ignore_me.jpg should NOT have been created because its in the docker ignore file")
		}

		if _, err = os.Stat(fmt.Sprintf("single-file-upload/%s", "path/to/dir/child/ignore_me.jpg")); !os.IsNotExist(err) {
			t.Errorf("path/to/dir/child/ignore_me.jpg should NOT have been created because its in the docker ignore file")
		}

		if _, err = os.Stat(fmt.Sprintf("single-file-upload/%s", "path/to/dir/child/nest.png")); os.IsNotExist(err) {
			t.Errorf("path/to/dir/child/nest.png should have been created and it's image dir moved to the root of the single file directory, was not")
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
	err := createTestMD(dockerMDContent)
	err = os.MkdirAll("../path/to/dir/child", os.FileMode(0777))
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
		result, err := createNewTarget("test.md", []string{}, []string{"/path/to/dir"})
		if err != nil {
			t.Errorf("Attempting to createNewTarget errored: %s\n", err)
		}

		if result != "single-file-upload" {
			t.Errorf("result should be the temp directory with the target markdown, '%s'", result)
		}

		if _, err := os.Stat(fmt.Sprintf("single-file-upload/%s", "test.md")); os.IsNotExist(err) {
			t.Errorf("test.md should have been created")
		}

		if _, err = os.Stat(fmt.Sprintf("single-file-upload/%s", "path/to/dir/root.png")); os.IsNotExist(err) {
			t.Errorf("path/to/dir/root.png should have been created and it's image dir moved to the root of the single file directory, was not")
		}

		if _, err = os.Stat(fmt.Sprintf("single-file-upload/%s", "path/to/dir/child/nest.png")); os.IsNotExist(err) {
			t.Errorf("path/to/dir/child/nest.png should have been created and it's image dir moved to the root of the single file directory, was not")
		}

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
	err := createTestMD(dockerMDContent)
	err = os.MkdirAll("../../path/to/dir/child", os.FileMode(0777))
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
		result, err := createNewTarget("test.md", []string{}, []string{"/path/to/dir"})
		if err != nil {
			t.Errorf("Attempting to createNewTarget errored: %s\n", err)
		}

		if result != "single-file-upload" {
			t.Errorf("result should be the temp directory with the target markdown, '%s'", result)
		}

		if _, err := os.Stat(fmt.Sprintf("single-file-upload/%s", "test.md")); os.IsNotExist(err) {
			t.Errorf("test.md should have been created")
		}

		if _, err = os.Stat(fmt.Sprintf("single-file-upload/%s", "path/to/dir/root.png")); os.IsNotExist(err) {
			t.Errorf("path/to/dir/root.png should have been created and it's image dir moved to the root of the single file directory, was not")
		}

		if _, err = os.Stat(fmt.Sprintf("single-file-upload/%s", "path/to/dir/child/nest.png")); os.IsNotExist(err) {
			t.Errorf("path/to/dir/child/nest.png should have been created and it's image dir moved to the root of the single file directory, was not")
		}

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

func createTestMD(content string) error {
	f, err := os.OpenFile("test.md", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
