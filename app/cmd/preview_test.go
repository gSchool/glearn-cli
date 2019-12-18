package cmd

import (
	"fmt"
	"os"
	"testing"
)

func Test_createNewTarget(t *testing.T) {
	result, err := createNewTarget("../../fixtures/test-links/nested/test.md", []string{"./mrsmall-invert.png", "../mrsmall.png", "mkdown.md", "../image/nested-small.png", "deeper/deep-small.png"})
	fmt.Println(result)
	if err != nil {
		t.Errorf("Attempting to createNewTarget errored: %s\n", err)
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

	if _, err = os.Stat(fmt.Sprintf("single-file-upload/%s", "mkdown.png")); !os.IsNotExist(err) {
		t.Errorf("linked markdown files should not copied")
	}
}
