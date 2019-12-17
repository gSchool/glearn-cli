package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

const withConfigFixture = "../../fixtures/test-block-with-config"

func Test_PreviewDetectsConfig(t *testing.T) {
	createdConfig, _ := doesConfigExistOrCreate(withConfigFixture, "")
	if createdConfig {
		t.Errorf("Created a config when one existed")
	}
}

const withNoConfigFixture = "../../fixtures/test-block-no-config"

func Test_PreviewBuildsAutoConfig(t *testing.T) {
	createdConfig, _ := doesConfigExistOrCreate(withNoConfigFixture, "")
	if createdConfig == false {
		t.Errorf("Should of created a config file")
	}

	b, err := ioutil.ReadFile(withNoConfigFixture + "/autoconfig.yaml")
	if err != nil {
		fmt.Print(err)
	}

	config := string(b)

	if !strings.Contains(config, "Title: Unit 1") {
		t.Errorf("Autoconfig should have a unit title of Unit 1")
	}

	if !strings.Contains(config, "Path: /units/test.md") {
		t.Errorf("Autoconfig should have a lesson with a path of /units/test.md")
	}
}

const withNoUnitsDirFixture = "../../fixtures/test-block-no-units-dir"

func Test_PreviewBuildsAutoConfigDeclaredUnitsDir(t *testing.T) {
	createdConfig, _ := doesConfigExistOrCreate(withNoUnitsDirFixture, "foo")
	if createdConfig == false {
		t.Errorf("Should of created a config file")
	}

	b, err := ioutil.ReadFile(withNoUnitsDirFixture + "/autoconfig.yaml")
	if err != nil {
		fmt.Print(err)
	}

	config := string(b)

	if !strings.Contains(config, "Title: Foo") {
		t.Errorf("Autoconfig should have a unit title of Foo")
	}

	if !strings.Contains(config, "Path: /foo/test.md") {
		t.Errorf("Autoconfig should have a lesson with a path of /foo/test.md")
	}
}
