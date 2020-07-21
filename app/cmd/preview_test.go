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

func Test_createNewTarget(t *testing.T) {
	result, err := createNewTarget("../../fixtures/test-links/nested/test.md", []string{"./mrsmall-invert.png", "../mrsmall.png", "../image/nested-small.png", "deeper/deep-small.png"})
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

func Test_createNewTargetSingleFileThatIsSQL(t *testing.T) {
	output := captureOutput(func() {
		_, err := createNewTarget("test.md", []string{"/data/some.sql", "image/nested-small.png"})
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

	err := os.RemoveAll("single-file-upload")
	if err != nil {
		t.Errorf("could not remove single-file-upload directory: %s\n", err)
	}
}

func Test_createNewTargetSingleFileSQLWithImage(t *testing.T) {
	output := captureOutput(func() {
		_, err := createNewTarget("test.md", []string{"/data/some.sql"})
		_, err = os.Stat(fmt.Sprintf("single-file-upload/%s", "data/some.sql"))
		if err == nil {
			t.Errorf("data/some.sql should have been copied over and it was not")
		}
	})

	if strings.Contains(output, "Link not found with path") {
		t.Errorf("output should not print 'Link not found with path', output was:\n%s\n", output)
	}
}

func Test_createNewTargetSingleFile(t *testing.T) {
	output := captureOutput(func() {
		result, err := createNewTarget("test.md", []string{"./image/nested-small.png", "image/nested-small.png", "../nested-small.png"})
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
}

func captureOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}
