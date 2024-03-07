package questions

import (
	"github.com/spf13/cobra"
)

func NewTestableProjectCommand(params NewQuestionCommandParams) *cobra.Command {
	params.name = "Testable Project"
	params.abbr = "tpr"
	params.maxTemplate = testableProjectTemplate
	params.minTemplate = testableProjectTemplateMin
	params.long = testableProjectLongDescription

	return createQuestionCommand(params)
}

const testableProjectLongDescription = `Testable project challenges allow students to submit entire repositories of
code to be evaluated against automated actions of any kind. This includes tests
passing, code coverage, linting, cyclomatic complexity, and more.

To do this you need two things --

  1. An upstream repo in Docker
  2. A project challenge where the students can submit their exercise repo
     to Learn

More information about creating an upstream repository can be found here:
https://learn-2.galvanize.com/cohorts/667/blocks/13/content_files/Testing-Project-Challenges.md`

const testableProjectTemplate = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->
<!-- Replace everything in square brackets [] and remove brackets  -->

### !challenge

* type: testable-project
* id: %s
* title: [text, a short question title]
* upstream: [URL, the upstream repo URL like https://github.com/gSchool/simple-compose-upstream]
* validate_fork: false
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

const testableProjectTemplateMin = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->

### !challenge

* type: testable-project
* id: %s
* title:
* upstream:
* validate_fork: false

##### !question



##### !end-question

##### !placeholder



##### !end-placeholder
<optional-attributes>
### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`
