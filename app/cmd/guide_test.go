package cmd

import "testing"

const withWalkthroughFixture = "../../fixtures/test-folder-with-walkthrough"
const withNoWalkthroughFixture = "../../fixtures/test-folder-without-walkthrough"

func Test_DoesDirHaveWalkthrough(t *testing.T) {
	hasWalkthrough := doesCurrentDirHaveWalkthrough(withWalkthroughFixture)
	if hasWalkthrough == false {
		t.Errorf("Should of found a walkthrough folder in directory")
	}
	hasNoWalkthrough := doesCurrentDirHaveWalkthrough(withNoWalkthroughFixture)
	if hasWalkthrough == true {
		t.Errorf("Should not of found a walkthrough folder in directory")
	}
}
