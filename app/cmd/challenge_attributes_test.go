package cmd

import (
	"strings"
	"testing"
)

var questionTypes = []string{
	"mc",
	"cb",
	"tl",
	"sa",
	"nb",
	"pg",
	"or",
	"js",
	"ja",
	"py",
	"sq",
	"rb",
	"up",
	"cs",
	"pr",
	"tpr",
}

func Test_AttributesCommentsOccurInMaximalWithNoAttributes(t *testing.T) {
	PrintTemplate = false
	Minimal = false
	WithExplanation = false
	WithRubric = false
	WithHints = 0

	for _, command := range questionTypes {
		t.Run(command, func(t2 *testing.T) {
			template, _ := getTemp(command)
			rendered := template.renderTemplate()

			if !strings.Contains(rendered, blockHeader) {
				t2.Errorf("%s did not have block header", command)
			}
			if !strings.Contains(rendered, hintTemplateSilent) {
				t2.Errorf("%s did not have hint comment", command)
			}
			if !strings.Contains(rendered, rubricTemplateSilent) {
				t2.Errorf("%s did not have rubric comment", command)
			}
			if !strings.Contains(rendered, explanationTemplateSilent) {
				t2.Errorf("%s did not have explanation comment", command)
			}
		})
	}
}

func Test_ExplanationBlockAppearsWhenTrueInMaximal(t *testing.T) {
	PrintTemplate = false
	Minimal = false
	WithExplanation = true
	WithRubric = false
	WithHints = 0

	for _, command := range questionTypes {
		t.Run(command, func(t2 *testing.T) {
			template, _ := getTemp(command)
			rendered := template.renderTemplate()

			if !strings.Contains(rendered, explanationTemplate) {
				t2.Errorf("%s did not have explanation block", command)
			}
		})
	}
}

func Test_RubricBlockAppearsWhenTrueInMaximal(t *testing.T) {
	PrintTemplate = false
	Minimal = false
	WithExplanation = false
	WithRubric = true
	WithHints = 0

	for _, command := range questionTypes {
		t.Run(command, func(t2 *testing.T) {
			template, _ := getTemp(command)
			rendered := template.renderTemplate()

			if !strings.Contains(rendered, rubricTemplate) {
				t2.Errorf("%s did not have rubric block", command)
			}
		})
	}
}

func Test_HintBlockAppearsWhenTrueInMaximal(t *testing.T) {
	PrintTemplate = false
	Minimal = false
	WithExplanation = false
	WithRubric = true
	WithHints = 1

	for _, command := range questionTypes {
		t.Run(command, func(t2 *testing.T) {
			template, _ := getTemp(command)
			rendered := template.renderTemplate()

			if !strings.Contains(rendered, hintTemplate) {
				t2.Errorf("%s did not have hint block", command)
			}
		})
	}
}

func Test_AttributesCommentsAbsentInMinimalWithNoAttributes(t *testing.T) {
	PrintTemplate = false
	Minimal = true
	WithExplanation = false
	WithRubric = false
	WithHints = 0

	for _, command := range questionTypes {
		t.Run(command, func(t2 *testing.T) {
			template, _ := getTemp(command)
			rendered := template.renderTemplate()

			if strings.Contains(rendered, blockHeader) {
				t2.Errorf("%s did have block header", command)
			}
			if strings.Contains(rendered, hintTemplateSilent) {
				t2.Errorf("%s did have hint comment", command)
			}
			if strings.Contains(rendered, rubricTemplateSilent) {
				t2.Errorf("%s did have rubric comment", command)
			}
			if strings.Contains(rendered, explanationTemplateSilent) {
				t2.Errorf("%s did have explanation comment", command)
			}
		})
	}
}

func Test_ExplanationBlockAppearsWhenTrueInMinimal(t *testing.T) {
	PrintTemplate = false
	Minimal = true
	WithExplanation = true
	WithRubric = false
	WithHints = 0

	for _, command := range questionTypes {
		t.Run(command, func(t2 *testing.T) {
			template, _ := getTemp(command)
			rendered := template.renderTemplate()

			if !strings.Contains(rendered, explanationTemplateMin) {
				t2.Errorf("%s did not have explanation block", command)
			}
		})
	}
}

func Test_RubricBlockAppearsWhenTrueInMinimal(t *testing.T) {
	PrintTemplate = false
	Minimal = true
	WithExplanation = false
	WithRubric = true
	WithHints = 0

	for _, command := range questionTypes {
		t.Run(command, func(t2 *testing.T) {
			template, _ := getTemp(command)
			rendered := template.renderTemplate()

			if !strings.Contains(rendered, rubricTemplateMin) {
				t2.Errorf("%s did not have rubric block", command)
			}
		})
	}
}

func Test_HintBlockAppearsWhenTrueInMinimal(t *testing.T) {
	PrintTemplate = false
	Minimal = true
	WithExplanation = false
	WithRubric = true
	WithHints = 1

	for _, command := range questionTypes {
		t.Run(command, func(t2 *testing.T) {
			template, _ := getTemp(command)
			rendered := template.renderTemplate()

			if !strings.Contains(rendered, hintTemplateMin) {
				t2.Errorf("%s did not have hint block", command)
			}
		})
	}
}
