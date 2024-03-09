package questions

import (
	"github.com/spf13/cobra"
)

func NewCheckBoxCommand(params NewQuestionCommandParams) *cobra.Command {
	params.name = "Check Box"
	params.abbr = "cb"
	params.maxTemplate = checkboxTemplate
	params.minTemplate = checkboxTemplateMin
	params.long = checkboxLongDescription

	return createQuestionCommand(params)
}

const checkboxLongDescription = `Checkbox challenges allow a student to submit multiple answers to a
multiple-choice question.`

const checkboxTemplate = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->
<!-- Replace everything in square brackets [] and remove brackets  -->

### !challenge

* type: checkbox
* id: %s
* title: [text, a short question title]
<!-- * points: [1] (optional, the number of points for scoring as a checkpoint) -->
<!-- * topics: [python, pandas] (Checkpoints only, optional the topics for analyzing points) -->

##### !question

[markdown, your question]

##### !end-question

##### !options

a| [Option 1]
b| [Option 2]
c| [Option 3]

##### !end-options

##### !answer

b|
c|

##### !end-answer

<optional-attributes>

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const checkboxTemplateMin = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->

### !challenge

* type: checkbox
* id: %s
* title:

##### !question



##### !end-question

##### !options

a|
b|
c|

##### !end-options

##### !answer

b|
c|

##### !end-answer
<optional-attributes>
### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`
