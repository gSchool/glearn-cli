package questions

import (
	"github.com/spf13/cobra"
)

func NewMultipleChoiceCommand(params NewQuestionCommandParams) *cobra.Command {
	params.name = "Multiple Choice"
	params.abbr = "mc"
	params.maxTemplate = multipleChoiceTemplate
	params.minTemplate = multipleChoiceTemplateMin
	params.long = multipleChoiceLongDescription

	return createQuestionCommand(params)
}

const multipleChoiceLongDescription = `Multiple Choice challenges allow a student to submit a single answer to a
multiple-choice question.`

const multipleChoiceTemplate = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->
<!-- Replace everything in square brackets [] and remove brackets  -->

### !challenge

* type: multiple-choice
* id: %s
* title: [text, a short question title]
<!-- * points: [1] (optional, the number of points for scoring as a checkpoint) -->
<!-- * topics: [python, pandas] (Checkpoints only. optional the topics for analyzing points) -->

##### !question

[markdown, your question]

##### !end-question

##### !options

a| [Option 1]
b| [Option 2]
c| [Option 3, etc]

##### !end-options

##### !answer

b|

##### !end-answer

<optional-attributes>

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const multipleChoiceTemplateMin = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->

### !challenge

* type: multiple-choice
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

a|

##### !end-answer
<optional-attributes>
### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`
