package files

import (
	"github.com/spf13/cobra"
)

func NewLessonCommand(params NewFileCommandParams) *cobra.Command {
	params.name = "Lesson"
	params.abbr = "ls"
	params.maxTemplate = lessonTemplate
	params.minTemplate = lessonTemplateMin

	return createFileCommand(params)
}

const lessonTemplate = `---
# BEGIN FILE CONFIGURATION YML HEADER >>>>>
# autoconfig.yml will use these settings. config.yml will override.
Type: Lesson
UID: %s
# DefaultVisibility: hidden # Uncomment this line to default Lesson to hidden. Please note, the default visibility setting is applied only during the initial sync of a course file within a cohort.
# END FILE CONFIGURATION YML HEADER <<<<<
---

# Title

<!--Lesson content can be markdown, videos, slides, images, gifs, etc. See examples of markdown formatting by running 'learn walkthrough' and previewing the tutorial. -->
<!--Lessons can include Challenges, which make the content interactive and give instructors visibility into student learning. -->
`

const lessonTemplateMin = `---
Type: Lesson
UID: %s
# DefaultVisibility: hidden
---

# Title
`
