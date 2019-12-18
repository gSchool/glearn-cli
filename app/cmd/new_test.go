package cmd

import "testing"

func Test_DoesDirHaveConfig(t *testing.T) {
	hasConfig, _ := doesCurrentDirHaveConfig(withConfigFixture)
	if hasConfig == false {
		t.Errorf("Should of found a config file in directory")
	}

	_, hasAutoConfig := doesCurrentDirHaveConfig(withNoConfigFixture)
	if hasAutoConfig == false {
		t.Errorf("Should of found an auto config file in directory")
	}
}
