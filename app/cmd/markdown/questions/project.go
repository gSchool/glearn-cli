package questions

import (
	"github.com/spf13/cobra"
)

func NewProjectCommand(params NewQuestionCommandParams) *cobra.Command {
	params.name = "Project"
	params.abbr = "pr"
	params.maxTemplate = projectTemplate
	params.minTemplate = projectTemplateMin
	params.long = projectLongDescription

	return createQuestionCommand(params)
}

const projectLongDescription = `Project Challenges allow a student to do work outside of Learn. After
completing the work, the student submits a link (typically from Github) to
their work for tracking and review within Learn.`

const projectTemplate = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->
<!-- Replace everything in square brackets [] and remove brackets  -->

### !challenge

* type: project
* id: %s
* title: [text, a short question title]
<!-- * points: [1] (optional, the number of points for scoring as a checkpoint) -->
<!-- * topics: [python, pandas] (Checkpoints only, optional the topics for analyzing points) -->

##### !question

[markdown, your question]

##### !end-question

##### !placeholder

[text, placeholder text for the input field]

##### !end-placeholder

<optional-attributes>

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const projectTemplateMin = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->

### !challenge

* type: project
* id: %s
* title:

##### !question



##### !end-question

##### !placeholder



##### !end-placeholder
<optional-attributes>
### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`
