package mdresourceparser

import (
	"strings"
	"testing"
)

func Test_ParseLink(t *testing.T) {
	tableTest := map[string][]string{
		"[example](linkresult)":   {"linkresult"},
		"[example](ends-in.md)":   {"ends-in.md"},
		"[](has-no-link-text)":    {"has-no-link-text"},
		"[more](than)[one](link)": {"than", "link"},
		`[more](than)
[one](line)
[with](links)
		`: {"than", "line", "links"},
		"[example](linkresult[contains](valid-link))": {"linkresult[contains](valid-link"}, // not actually supported, checks terminating link character
		"[)":                  {""},
		"[here](./../result)": {"./../result"},
		`var myarr = [];
myarr[0] = (val != otherval);`: {""},
		`var myarr = [(arg) => { console.log(arg) }];
myarr[0]("code-test-case");`: {"\"code-test-case\""},
	}

	for k, v := range tableTest {
		parser := New([]rune(k))
		parser.ParseResources()
		result := parser.Links
		if strings.Join(result, "") != strings.Join(v, "") {
			t.Errorf("Links %s expected %v but got %v", k, v, result)
		}
	}

	emptyParser := New([]rune("[empty]()"))
	emptyParser.ParseResources()
	if len(emptyParser.Links) != 0 {
		t.Errorf("A valid link syntax with an empty link should produce no links")
	}
}

func Test_ParseDockerDirectoryPaths(t *testing.T) {
	tableTest := map[string][]string{
		challengeContent: {"/path/to/dir"},
	}
	for k, v := range tableTest {
		parser := New([]rune(k))
		_, result, _, _ := parser.ParseResources()
		if strings.Join(result, "") != strings.Join(v, "") {
			t.Errorf("DockerDirectoryPaths %s expected %s but got %v", k, v, result)
		}
	}
}

func Test_ParseSeveralChallengeContents(t *testing.T) {
	result := New([]rune(multipleChallengeContent))
	// note that testFilePaths refer to challenge tests returned, not a variable to test against
	dataPaths, dockerDirectoryPaths, testFilePaths, setupFilePaths := result.ParseResources()
	if len(dockerDirectoryPaths) != 1 {
		t.Fatalf("length DockerDirectoryPaths expected 1,  got %d", len(dockerDirectoryPaths))
	}
	if len(testFilePaths) != 2 {
		t.Fatalf("length TestFilePaths expected 1, got %d", len(testFilePaths))
	}
	if len(setupFilePaths) != 1 {
		t.Fatalf("length SetupFilePaths expected 1, got %d", len(setupFilePaths))
	}
	if len(dataPaths) != 1 {
		t.Fatalf("length dataPaths expected 1, got %d", len(dataPaths))
	}

	expectedDockerDirectoryPath := "/path/to/dir"
	expectedFirstTestFilePath := "/tests/title.js"
	expectedSecondTestFilePath := "/tests/sql.sql"
	expectedSetupFilePath := "/setups/title.js"
	expectedDataPath := "/data/sql.dump"

	if expectedDockerDirectoryPath != dockerDirectoryPaths[0] {
		t.Errorf("Expected DockerDirectoryPaths '%s', got '%s'", expectedDockerDirectoryPath, dockerDirectoryPaths[0])
	}
	if expectedFirstTestFilePath != testFilePaths[0] {
		t.Errorf("Expected first TestFilePaths '%s', got '%s'", expectedFirstTestFilePath, testFilePaths[0])
	}
	if expectedSecondTestFilePath != testFilePaths[1] {
		t.Errorf("Expected second TestFilePaths '%s', got '%s'", expectedSecondTestFilePath, testFilePaths[0])
	}
	if expectedSetupFilePath != setupFilePaths[0] {
		t.Errorf("Expected SetupFilePaths '%s', got '%s'", expectedSetupFilePath, setupFilePaths[0])
	}
	if expectedDataPath != dataPaths[0] {
		t.Errorf("Expected SetupFilePaths '%s', got '%s'", expectedSetupFilePath, dataPaths[0])
	}
}

func Test_hasPathBullet(t *testing.T) {
	tableTest := map[string]bool{
		"* content_after": true,
		"- content_after": true,
		"_ not_valid":     false,
	}

	for k, b := range tableTest {
		parser := New([]rune(k))
		if parser.hasPathBullet() != b {
			t.Errorf("Expected %s hasPathBullet to be %v, was not", k, b)
		}
	}
}

const minimalBullets = `
* docker_directory_path: /path/to/dir
* test_file: /tests/title.js
* setup_file: /setup/title.js
`

const multipleChallengeContent = `### !challenge

- type: custom-snippet
- language: text
- id: 8c406f4f-6428-498b-be24-6bd0a6c9096a
- title: Title
- docker_directory_path: /path/to/dir

##### !question

Question

##### !end-question

##### !placeholder

##### !end-placeholder

### !end-challenge

### !challenge

* type: code-snippet
* language: javascript18
* id: 8c406f4f-6428-498b-be24-6bd0a6c9096b
* title: title
* test_file: /tests/title.js
* setup_file: /setups/title.js

##### !question

Question

##### !end-question

##### !placeholder

##### !end-placeholder

### !end-challenge

### !challenge

* type: code-snippet
* language: sql
* id: 8c406f4f-6428-498b-be24-6bd0a6c9096c
* title: sql
* test_file: /tests/sql.sql
* data_path: /data/sql.dump

##### !question

Question

##### !end-question

##### !placeholder

##### !end-placeholder

### !end-challenge

`

const challengeContent = `### !challenge

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
