package questions

import (
	"github.com/spf13/cobra"
)

func NewNumberCommand(params NewQuestionCommandParams) *cobra.Command {
	params.name = "Number"
	params.abbr = "nb"
	params.maxTemplate = numberTemplate
	params.minTemplate = numberTemplateMin
	params.long = numberLongDescription

	return createQuestionCommand(params)
}

const numberLongDescription = `Number challenges allow a student to submit a number as the answer to a
question. The answer is evaluated numerically, and the student can answer with
a decimal or a fraction--so that things like 3/10, 30/100, .3, 0.3, 0.300 etc.
are all equivalent.

The number challenge is for floating point values only. To restrict the answer
to integers (no decimal allowed) use a short-answer challenge instead. For the
answer, use the regular expression  /^\s*[0-9]+\s*$/ which only matches whole
numbers.`

const numberTemplate = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->
<!-- Replace everything in square brackets [] and remove brackets  -->

### !challenge

* type: number
* id: %s
* title: [text, a short question title]
<!-- * decimal: [optional number, decimal points to user for answer evaluation] -->
<!-- * points: [1] (optional, the number of points for scoring as a checkpoint) -->
<!-- * topics: [python, pandas] (Checkpoints only, optional the topics for analyzing points) -->

##### !question

[markdown, your question]

##### !end-question

##### !placeholder

[text, placeholder text for input field]

##### !end-placeholder

##### !answer

[number, the correct answer]

##### !end-answer

<optional-attributes>

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const numberTemplateMin = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->

### !challenge

* type: number
* id: %s
* title:
* decimal:

##### !question



##### !end-question

##### !answer



##### !end-answer
<optional-attributes>
### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`
