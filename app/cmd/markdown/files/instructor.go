package files

import (
	"github.com/spf13/cobra"
)

func NewInstructorCommand(params NewFileCommandParams) *cobra.Command {
	params.name = "Instructor"
	params.abbr = "in"
	params.maxTemplate = instructorTemplate
	params.minTemplate = instructorTemplateMin

	return createFileCommand(params)
}

const instructorTemplate = `---
# BEGIN FILE CONFIGURATION YML HEADER >>>>>
# autoconfig.yml will use these settings. config.yml will override.
Type: Instructor
UID: %s
# END FILE CONFIGURATION YML HEADER <<<<<
---

# Title

<!--An Instructor file can have all of the same markdown and challenges as a lesson. Instructor files are only viewable by instructors. -->
`

const instructorTemplateMin = `---
Type: Instructor
UID: %s
---

# Title
`
