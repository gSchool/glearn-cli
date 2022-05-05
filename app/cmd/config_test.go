package cmd

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	yaml "gopkg.in/yaml.v2"
)

const withConfigFixture = "../../fixtures/test-block-with-config"

func Test_PreviewDetectsConfig(t *testing.T) {
	createdConfig, _ := previewFindOrCreateConfig(withConfigFixture, false, []string{})

	if createdConfig {
		t.Errorf("Created a config when one existed")
	}
}

const withNoConfigFixture = "../../fixtures/test-block-no-config"

func Test_PublishBuildsAutoConfig(t *testing.T) {
	gitTopLevelCmd = "echo ../../fixtures/test-block-no-config"
	createdConfig, err := publishFindOrCreateConfig(withNoConfigFixture)
	if err != nil {
		t.Errorf("Should not have errored but got error: '%s'\n", err)
	}
	if createdConfig == false {
		t.Errorf("Should of created a config file")
	}

	b, err := ioutil.ReadFile(withNoConfigFixture + "/autoconfig.yaml")
	if err != nil {
		fmt.Print(err)
	}

	var configMade ConfigYaml
	err = yaml.Unmarshal(b, &configMade)
	if err != nil {
		t.Errorf("File created could not unmarshal into ConfigYaml struct: %s", err)
	}
	if len(configMade.Standards) != 4 {
		t.Errorf("test-block-with-config fixture should have made 4 standards but made %d", len(configMade.Standards))
	}
	// standard 1
	standardOne := configMade.Standards[0]
	if standardOne.Title != "Checkpoint" {
		t.Errorf("test-block-with-config fixture first standard should have title Checkpoint, but had '%s'", standardOne.Title)
	}
	if len(standardOne.ContentFiles) != 1 {
		t.Errorf("test-block-with-config fixture first standard should have 1 content file but had %d", len(standardOne.ContentFiles))
	}
	if standardOne.ContentFiles[0].Type != "Checkpoint" {
		t.Errorf("test-block-with-config fixture first standard first content file should be of type Checkpoint but was '%s'", standardOne.ContentFiles[0].Type)
	}
	if standardOne.ContentFiles[0].Path != "/units/01-checkpoint/another-checkpoint.md" {
		t.Errorf("test-block-with-config fixture first standard first content file path should be '/units/01-checkpoint/another-checkpoint.md' but was '%s'", standardOne.ContentFiles[0].Path)
	}

	// standard 2
	standardTwo := configMade.Standards[1]
	if standardTwo.Title != "Resource" {
		t.Errorf("test-block-with-config fixture second standard should have title Checkpoint, but had '%s'", standardTwo.Title)
	}
	if len(standardTwo.ContentFiles) != 1 {
		t.Errorf("test-block-with-config fixture second standard should have 1 content file but had %d", len(standardTwo.ContentFiles))
	}
	if standardTwo.ContentFiles[0].Type != "Resource" {
		t.Errorf("test-block-with-config fixture second standard first content file should be of type Resource but was '%s'", standardTwo.ContentFiles[0].Type)
	}
	if standardTwo.ContentFiles[0].Path != "/units/03.resource/name.resource.md" {
		t.Errorf("test-block-with-config fixture second standard first content file path should be '/units/03.resource/name.resource.md' but was '%s'", standardTwo.ContentFiles[0].Path)
	}

	// standard 3
	standardThree := configMade.Standards[2]
	if standardThree.Title != "Title From Yaml" {
		t.Errorf("test-block-with-config fixture third standard should have yaml supplied title 'Title From Yaml', but had '%s'", standardThree.Title)
	}
	if standardThree.Description != "My description here" {
		t.Errorf("test-block-with-config fixture third standard should have yaml supplied description 'My description here', but had '%s'", standardThree.Description)
	}
	if standardThree.UID != "02210548f12da09aa7a0bd1f1308c423" {
		t.Errorf("test-block-with-config fixture third standard should have yaml supplied UID '02210548f12da09aa7a0bd1f1308c423', but had '%s'", standardThree.UID)
	}
	if len(standardThree.SuccessCriteria) != 2 {
		t.Errorf("test-block-with-config fixture third standard should have 2 SuccessCriteria but had %d", len(standardThree.SuccessCriteria))
	}
	if standardThree.SuccessCriteria[0] != "success" {
		t.Errorf("test-block-with-config fixture third standard first success criteria should be 'success' but was %s", standardThree.SuccessCriteria[0])
	}
	if standardThree.SuccessCriteria[1] != "criteria" {
		t.Errorf("test-block-with-config fixture third standard first success criteria should be 'criteria' but was %s", standardThree.SuccessCriteria[1])
	}

	// standard 4
	standardFour := configMade.Standards[3]
	if standardFour.Title != "Unit 1" {
		t.Errorf("test-block-with-config fixture fourth standard should have title Unit 1, but had '%s'", standardFour.Title)
	}
	if len(standardFour.ContentFiles) != 5 {
		t.Errorf("test-block-with-config fixture fourth standard should have 5 content files but had %d", len(standardFour.ContentFiles))
	}
	// standard 4 file 1
	if standardFour.ContentFiles[0].Type != "Lesson" {
		t.Errorf("test-block-with-config fixture fourth standard first content file should be of type Lesson but was '%s'", standardFour.ContentFiles[0].Type)
	}
	if standardFour.ContentFiles[0].Path != "/units/file.hidden.file.md" {
		t.Errorf("test-block-with-config fixture fourth standard first content file path should be '/units/file.hidden.file.md' but was '%s'", standardFour.ContentFiles[0].Path)
	}
	if standardFour.ContentFiles[0].DefaultVisibility != "hidden" {
		t.Errorf("test-block-with-config fixture fourth standard first content file DefaultVisibility should be 'hidden' but was '%s'", standardFour.ContentFiles[0].DefaultVisibility)
	}
	// standard 4 file 2
	if standardFour.ContentFiles[1].Type != "Survey" {
		t.Errorf("test-block-with-config fixture fourth standard second content file should be of type Survey but was '%s'", standardFour.ContentFiles[1].Type)
	}
	if standardFour.ContentFiles[1].Path != "/units/file.survey.md" {
		t.Errorf("test-block-with-config fixture fourth standard second content file path should be '/units/file.survey.md' but was '%s'", standardFour.ContentFiles[1].Path)
	}
	// standard 4 file 3
	if standardFour.ContentFiles[2].Type != "Resource" {
		t.Errorf("test-block-with-config fixture fourth standard third content file should be of type Resource but was '%s'", standardFour.ContentFiles[2].Type)
	}
	if standardFour.ContentFiles[2].Path != "/units/hidden.resource.md" {
		t.Errorf("test-block-with-config fixture fourth standard third content file path should be '/units/hidden.resource.md' but was '%s'", standardFour.ContentFiles[2].Path)
	}
	if standardFour.ContentFiles[2].DefaultVisibility != "hidden" {
		t.Errorf("test-block-with-config fixture fourth standard third content file DefaultVisibility should be 'hidden' but was '%s'", standardFour.ContentFiles[2].DefaultVisibility)
	}
	// standard 4 file 4
	if standardFour.ContentFiles[3].Type != "Instructor" {
		t.Errorf("test-block-with-config fixture fourth standard fourth content file should be of type Instructor but was '%s'", standardFour.ContentFiles[3].Type)
	}
	if standardFour.ContentFiles[3].Path != "/units/teacher-instructor.md" {
		t.Errorf("test-block-with-config fixture fourth standard fourth content file path should be '/units/teacher-instructor.md' but was '%s'", standardFour.ContentFiles[3].Path)
	}
	// standard 4 file 5 values set from header
	if standardFour.ContentFiles[4].Type != "Checkpoint" {
		t.Errorf("test-block-with-config fixture fourth standard fifth content file should be of type Instructor but was '%s'", standardFour.ContentFiles[4].Type)
	}
	if standardFour.ContentFiles[4].Path != "/units/test.md" {
		t.Errorf("test-block-with-config fixture fourth standard fifth content file path should be '/units/test.md' but was '%s'", standardFour.ContentFiles[4].Path)
	}
	if standardFour.ContentFiles[4].DefaultVisibility != "hidden" {
		t.Errorf("test-block-with-config fixture fourth standard fifth content file DefaultVisibility should be 'hidden' but was '%s'", standardFour.ContentFiles[4].DefaultVisibility)
	}
	if standardFour.ContentFiles[4].UID != "abc123" {
		t.Errorf("test-block-with-config fixture fourth standard fifth content file UID should be 'abc123' but was '%s'", standardFour.ContentFiles[4].UID)
	}
	if standardFour.ContentFiles[4].MaxCheckpointSubmissions != 1 {
		t.Errorf("test-block-with-config fixture fourth standard fifth content file MaxCheckpointSubmissions should be 1 but was %d", standardFour.ContentFiles[4].MaxCheckpointSubmissions)
	}
	if !standardFour.ContentFiles[4].EmailOnCompletion {
		t.Errorf("test-block-with-config fixture fourth standard fifth content file EmailOnCompletion should be true but was false")
	}
	if standardFour.ContentFiles[4].TimeLimit != 45 {
		t.Errorf("test-block-with-config fixture fourth standard fifth content file TimeLimit should be 45 but was %d", standardFour.ContentFiles[4].TimeLimit)
	}
	if !standardFour.ContentFiles[4].Autoscore {
		t.Errorf("test-block-with-config fixture fourth standard fifth content file Autscore should be true but was false")
	}
}

const withNoUnitsDirFixture = "../../fixtures/test-block-no-units-dir"

func Test_PreviewBuildsAutoConfigDeclaredUnitsDir(t *testing.T) {
	UnitsDirectory = "foo"
	createdConfig, _ := previewFindOrCreateConfig(withNoUnitsDirFixture, false, []string{})
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
	gitTopLevelCmd = "echo ../../fixtures/test-block-no-units-dir"
	createdConfig, err := previewFindOrCreateConfig(withNoUnitsDirFixture+"/single_unit", false, []string{})

	if createdConfig == true {
		t.Errorf("Should not of created a config file")
	}

	if err == nil {
		t.Errorf("Should of alerted user that no units were found and single unit preview is not supported")
	}
}

func Test_AutoConfigAddsInFileTypesOrVisibility(t *testing.T) {
	gitTopLevelCmd = "echo ../../fixtures/test-block-no-config"
	createdConfig, _ := previewFindOrCreateConfig(withNoConfigFixture, false, []string{})
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

	if !strings.Contains(config, "Type: Survey") {
		t.Errorf("Autoconfig should have a content path of survey but the type should not of changed")
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
	createdConfig, _ := previewFindOrCreateConfig(withNoConfigFixture, false, []string{})
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
	createdConfig, _ := previewFindOrCreateConfig(withNoConfigFixture, false, []string{"/units"})
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

func Test_findConfigMethodReturnsProperConfig(t *testing.T) {
	previewFindOrCreateConfig(withNoConfigFixture, false, []string{})

	configString, _ := findConfig(withNoConfigFixture)

	if configString == "" {
		t.Errorf("Should of found a config or autoconig file")
	}
}

func Test_ParseConfigFileForPaths(t *testing.T) {
	previewFindOrCreateConfig(withNoConfigFixture, false, []string{})
	paths, err := parseConfigAndGatherLinkedPaths(withNoConfigFixture)

	if err != nil || len(paths) == 0 {
		t.Errorf("Should of parse the yaml and gathered some content file paths")
	}
}
