package others

import (
	"github.com/spf13/cobra"
)

func NewDistributeCodeCommand(params NewOtherCommandParams) *cobra.Command {
	params.name = "Distribute Code"
	params.abbr = "dc"
	params.maxTemplate = distributeCodeTemplate
	params.minTemplate = distributeCodeTemplateMin

	return createOtherCommand(params)
}

const distributeCodeTemplate = `<!-- Replace everything in square brackets [] and remove brackets  -->
<!-- This button can be added anywhere except inside of a challenge -->
<!-- This can only be used with a single student repository model cohort. -->
### !distribute-code

* student_folder_path: [text, GitLab folder path that code will be distributed to in student's cohort repo URL (can be blank)]
* repository_url: [text, GitLab URL that code will be distributed from]

### !end-distribute-code`

const distributeCodeTemplateMin = `### !distribute-code

* student_folder_path:
* repository_url:

### !end-distribute-code`
