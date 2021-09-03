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
		parser.ParseResources()
		result := parser.DockerDirectoryPaths
		if strings.Join(result, "") != strings.Join(v, "") {
			t.Errorf("DockerDirectoryPaths %s expected %s but got %v", k, v, result)
		}
	}
}

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
