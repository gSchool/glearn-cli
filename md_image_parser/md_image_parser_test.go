package mdimageparser

import (
	"strings"
	"testing"
)

func Test_ParseImage(t *testing.T) {
	tableTest := map[string][]string{
		"[example](linkresult)":   []string{"linkresult"},
		"[example]()":             []string{""},
		"[](has-no-link-text)":    []string{"has-no-link-text"},
		"[more](than)[one](link)": []string{"than", "link"},
		`[more](than)
[one](line)
[with](links)
		`: []string{"than", "line", "links"},
		"[example](linkresult[contains](valid-link))": []string{"linkresult[contains](valid-link"}, // not actually supported, checks terminating link character
	}

	for k, v := range tableTest {
		parser := New(k)
		parser.ParseImages()
		result := parser.Images
		if strings.Join(result, "") != strings.Join(v, "") {
			t.Errorf("ParseImages %s expected %v but got %v", k, v, result)
		}
	}
}
