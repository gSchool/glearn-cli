package files

import (
	"github.com/spf13/cobra"
)

func NewSurveyCommand(params NewFileCommandParams) *cobra.Command {
	params.name = "Survey"
	params.abbr = "sv"
	params.maxTemplate = surveyTemplate
	params.minTemplate = surveyTemplateMin

	return createFileCommand(params)
}

const surveyTemplate = `---
# BEGIN FILE CONFIGURATION YML HEADER >>>>>
# autoconfig.yml will use these settings. config.yml will override.
Type: Survey
UID: %s
# DefaultVisibility: hidden # Uncomment this line to default Survey to hidden
# END FILE CONFIGURATION YML HEADER <<<<<
---

# Title

<!--A Survey can have any markdown. See examples of markdown formatting by running 'learn walkthrough' and previewing the tutorial. -->
<!--A Survey must include include one or more Challenges, which are the survey questions a student will answer. -->
`

const surveyTemplateMin = `---
Type: Survey
UID: %s
# DefaultVisibility: hidden
---

# Title
`
