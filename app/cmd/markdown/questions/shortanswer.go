package questions

import (
	"github.com/spf13/cobra"
)

func NewShortAnswerCommand(params NewQuestionCommandParams) *cobra.Command {
	params.name = "Short Answer"
	params.abbr = "sa"
	params.maxTemplate = shortAnswerTemplate
	params.minTemplate = shortAnswerTemplateMin
	params.long = shortAnswerLongDescription

	return createQuestionCommand(params)
}

const shortAnswerLongDescription = `Short-answer challenges allow a student to submit a short answer, usually a
single word, to answer a question. By default, the answer is evaluated as a
case-insensitive exact match but can also be evaluated as a regex.`

const shortAnswerTemplate = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->
<!-- Replace everything in square brackets [] and remove brackets  -->

### !challenge

* type: short-answer
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

##### !answer

[text or regex, the answer (if regex wrap in /)]

##### !end-answer

<optional-attributes>

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const shortAnswerTemplateMin = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->

### !challenge

* type: short-answer
* id: %s
* title:

##### !question



##### !end-question

##### !answer



##### !end-answer
<optional-attributes>
### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`
