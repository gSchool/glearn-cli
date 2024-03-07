package templates

import (
	"regexp"
	"strings"
	"testing"
)

func Test_StaticTemplateJustDoesStatic(t *testing.T) {
	temp := NewStaticTemplate("Name", "Content")

	if temp.GetName() != "Name" {
		t.Error("Did not get the expected name")
	}
	if temp.Render() != "Content" {
		t.Error("Render modified content")
	}
	if temp.GetUnrenderedContent() != "Content" {
		t.Error("UnrenderedContent did not return content")
	}
}

func Test_IdTemplateReplacesFormatHolderWithUuid(t *testing.T) {
	temp := NewIdTemplate("Name", "Content %s")
	re := regexp.MustCompile(`Content [[:alnum:]]{8}-[[:alnum:]]{4}-[[:alnum:]]{4}-[[:alnum:]]{4}-[[:alnum:]]{12}`)

	if temp.GetName() != "Name" {
		t.Error("Did not get the expected name")
	}
	if !re.MatchString(temp.Render()) {
		t.Error("Render did not return the correct string: " + temp.Render())
	}
	if temp.GetUnrenderedContent() != "Content %s" {
		t.Error("UnrenderedContent did not return content")
	}
}

func Test_AttributeTemplateReplacesFormatHolderWithUuid(t *testing.T) {
	temp := NewAttributeTemplate("Name", "Content %s", true, true, true, 10)
	re := regexp.MustCompile(`Content [[:alnum:]]{8}-[[:alnum:]]{4}-[[:alnum:]]{4}-[[:alnum:]]{4}-[[:alnum:]]{12}`)

	if temp.GetName() != "Name" {
		t.Error("Did not get the expected name")
	}
	if !re.MatchString(temp.Render()) {
		t.Error("Render did not return the correct string: " + temp.Render())
	}
	if temp.GetUnrenderedContent() != "Content %s" {
		t.Error("UnrenderedContent did not return content")
	}
}

func Test_AttributeTemplateWithGeneratesBlockHeader(t *testing.T) {
	testCases := []struct {
		min bool
		we  bool
		wr  bool
		wh  int
	}{
		{false, false, false, 0},
		{false, false, false, 1},
		{false, false, true, 0},
		{false, false, true, 1},
		{false, true, false, 0},
		{false, true, false, 1},
		{false, true, true, 0},
		{false, true, true, 1},
		{true, false, false, 0},
		{true, false, false, 1},
		{true, false, true, 0},
		{true, false, true, 1},
		{true, true, false, 0},
		{true, true, false, 1},
		{true, true, true, 1},
		{true, true, true, 0},
	}
	for _, tc := range testCases {
		temp := NewAttributeTemplate("Name", "Content %s\n<optional-attributes>", tc.min, tc.we, tc.wr, tc.wh)
		content := temp.Render()

		if !tc.min && !(tc.we && tc.wr && tc.wh != 0) {
			if !strings.Contains(content, blockHeader) {
				t.Errorf("Did not contain block header\n%s", content)
			}
		} else {
			if strings.Contains(content, blockHeader) {
				t.Errorf("Contained block header\n%s", content)
			}
		}

	}
}

func Test_AttributeTemplateWithGeneratesExceptionAttribute(t *testing.T) {
	testCases := []struct {
		min bool
		we  bool
		wr  bool
		wh  int
	}{
		{false, false, false, 0},
		{false, false, false, 1},
		{false, false, true, 0},
		{false, false, true, 1},
		{false, true, false, 0},
		{false, true, false, 1},
		{false, true, true, 0},
		{false, true, true, 1},
		{true, false, false, 0},
		{true, false, false, 1},
		{true, false, true, 0},
		{true, false, true, 1},
		{true, true, false, 0},
		{true, true, false, 1},
		{true, true, true, 1},
		{true, true, true, 0},
	}
	for _, tc := range testCases {
		temp := NewAttributeTemplate("Name", "Content %s\n<optional-attributes>", tc.min, tc.we, tc.wr, tc.wh)
		content := temp.Render()

		if tc.min && tc.we {
			if !strings.Contains(content, explanationTemplateMin) {
				t.Errorf("Did not contain min exception\n%s", content)
			}
		} else if !tc.min && tc.we {
			if !strings.Contains(content, explanationTemplate) {
				t.Errorf("Did not contain exception\n%s", content)
			}
		} else if !tc.min {
			if !strings.Contains(content, explanationTemplateSilent) {
				t.Errorf("Did not contain exception comment")
			}
		}

	}
}

func Test_AttributeTemplateWithGeneratesRubricAttribute(t *testing.T) {
	testCases := []struct {
		min bool
		we  bool
		wr  bool
		wh  int
	}{
		{false, false, false, 0},
		{false, false, false, 1},
		{false, false, true, 0},
		{false, false, true, 1},
		{false, true, false, 0},
		{false, true, false, 1},
		{false, true, true, 0},
		{false, true, true, 1},
		{true, false, false, 0},
		{true, false, false, 1},
		{true, false, true, 0},
		{true, false, true, 1},
		{true, true, false, 0},
		{true, true, false, 1},
		{true, true, true, 1},
		{true, true, true, 0},
	}
	for _, tc := range testCases {
		temp := NewAttributeTemplate("Name", "Content %s\n<optional-attributes>", tc.min, tc.we, tc.wr, tc.wh)
		content := temp.Render()

		if tc.min && tc.wr {
			if !strings.Contains(content, rubricTemplateMin) {
				t.Errorf("Did not contain min rubric\n%s", content)
			}
		} else if !tc.min && tc.wr {
			if !strings.Contains(content, rubricTemplate) {
				t.Errorf("Did not contain rubric\n%s", content)
			}
		} else if !tc.min {
			if !strings.Contains(content, rubricTemplateSilent) {
				t.Errorf("Did not contain rubric comment")
			}
		}

	}
}

func Test_AttributeTemplateWithGeneratesHintAttribute(t *testing.T) {
	testCases := []struct {
		min bool
		we  bool
		wr  bool
		wh  int
	}{
		{false, false, false, 0},
		{false, false, false, 1},
		{false, false, true, 0},
		{false, false, true, 1},
		{false, true, false, 0},
		{false, true, false, 1},
		{false, true, true, 0},
		{false, true, true, 1},
		{true, false, false, 0},
		{true, false, false, 1},
		{true, false, true, 0},
		{true, false, true, 1},
		{true, true, false, 0},
		{true, true, false, 1},
		{true, true, true, 1},
		{true, true, true, 0},
	}
	for _, tc := range testCases {
		temp := NewAttributeTemplate("Name", "Content %s\n<optional-attributes>", tc.min, tc.we, tc.wr, tc.wh)
		content := temp.Render()

		if tc.min && tc.wh > 0 {
			if !strings.Contains(content, hintTemplateMin) {
				t.Errorf("Did not contain min hint\n%s", content)
			}
		} else if !tc.min && tc.wh > 0 {
			if !strings.Contains(content, hintTemplate) {
				t.Errorf("Did not contain hint\n%s", content)
			}
		} else if !tc.min {
			if !strings.Contains(content, hintTemplateSilent) {
				t.Errorf("Did not contain hint comment")
			}
		}

	}
}
