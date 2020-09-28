package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"
)

const withConfigFixture = "../../fixtures/test-block-with-config"

func Test_PreviewDetectsConfig(t *testing.T) {
	createdConfig, _ := doesConfigExistOrCreate(withConfigFixture, false, []string{})
	if createdConfig {
		t.Errorf("Created a config when one existed")
	}
}

const withNoConfigFixture = "../../fixtures/test-block-no-config"

func Test_PreviewBuildsAutoConfig(t *testing.T) {
	createdConfig, _ := doesConfigExistOrCreate(withNoConfigFixture, false, []string{})
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
	UnitsDirectory = "foo"
	createdConfig, _ := doesConfigExistOrCreate(withNoUnitsDirFixture, false, []string{})
	if createdConfig == false {
		t.Errorf("Should of created a config file")
	}
	UnitsDirectory = ""

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

func Test_PreviewBuildFailsWhenPreviewingSingleUnit(t *testing.T) {
	createdConfig, err := doesConfigExistOrCreate(withNoUnitsDirFixture+"/single_unit", false, []string{})

	if createdConfig == true {
		t.Errorf("Should not of created a config file")
	}

	if err == nil {
		t.Errorf("Should of alerted user that no units where found and single unit preview is not supported")
	}
}

func Test_AutoConfigAddsInFileTypesOrVisibility(t *testing.T) {
	createdConfig, _ := doesConfigExistOrCreate(withNoConfigFixture, false, []string{})
	if createdConfig == false {
		t.Errorf("Should of created a config file")
	}

	b, err := ioutil.ReadFile(withNoConfigFixture + "/autoconfig.yaml")
	if err != nil {
		fmt.Print(err)
	}

	config := string(b)

	if !strings.Contains(config, "Type: Checkpoint") {
		t.Errorf("Autoconfig should have a content path of checkpoint but the type should not of changed")
	}

	if !strings.Contains(config, "Type: Instructor") {
		t.Errorf("Autoconfig should have a content file of type Instructor")
	}

	if !strings.Contains(config, "Type: Resource") {
		t.Errorf("Autoconfig should have a content path of resource but the type should not of changed")
	}

	if !strings.Contains(config, "DefaultVisibility: hidden") {
		t.Errorf("Autoconfig should have a content file of with a DefaultVisibility of hidden")
	}
}

func Test_IgnoresFilesAndUnitsThatStartWithTwoUnderscores(t *testing.T) {
	createdConfig, _ := doesConfigExistOrCreate(withNoConfigFixture, false, []string{})
	if createdConfig == false {
		t.Errorf("Should of created a config file")
	}

	b, err := ioutil.ReadFile(withNoConfigFixture + "/autoconfig.yaml")
	if err != nil {
		fmt.Print(err)
	}

	config := string(b)

	if strings.Contains(config, "__skip") {
		t.Errorf("Autoconfig have units that start with __")
	}

	if strings.Contains(config, "__skipthis.md") {
		t.Errorf("Autoconfig have contentfiles that start with __")
	}
}

func Test_IgnoresExcludedFiles(t *testing.T) {
	createdConfig, _ := doesConfigExistOrCreate(withNoConfigFixture, false, []string{"/units"})
	if createdConfig == false {
		t.Errorf("Should of created a config file")
	}

	b, err := ioutil.ReadFile(withNoConfigFixture + "/autoconfig.yaml")
	if err != nil {
		fmt.Print(err)
	}

	config := string(b)

	if strings.Contains(config, "Title: Unit 1") {
		t.Errorf("Autoconfig should have excluded a unit titled Unit 1")
	}

	if strings.Contains(config, "Path: /units/test.md") {
		t.Errorf("Autoconfig should have excluded a lesson with a path of /units/test.md")
	}
}
