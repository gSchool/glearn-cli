package questions

import (
	"github.com/spf13/cobra"
)

func NewOrderingCommand(params NewQuestionCommandParams) *cobra.Command {
	params.name = "Ordering"
	params.abbr = "or"
	params.maxTemplate = orderingTemplate
	params.minTemplate = orderingTemplateMin
	params.long = orderingLongDescription

	return createQuestionCommand(params)
}

const orderingLongDescription = `These challenges ask students to arrange options in a particular order. They
can be used for sequence-style questions (What's the correct order for
accomplishing a task?) or ranking questions (rank these items from most to
least impactful).`

const orderingTemplate = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->
<!-- Replace everything in square brackets [] and remove brackets  -->

### !challenge

* type: ordering
* id: %s
* title: [text, a short question title]
<!-- * points: [1] (optional, the number of points for scoring as a checkpoint) -->
<!-- * topics: [python, pandas] (Checkpoints only, optional the topics for analyzing points) -->

##### !question

[markdown, your question]

##### !end-question

##### !answer

1. [Option 1 in the correct first position, options will be randomized for students]
1. [Option 2 in the correct second position]
1. [Option 3 in the correct third position]

##### !end-answer

<optional-attributes>

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const orderingTemplateMin = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->

### !challenge

* type: ordering
* id: %s
* title:

##### !question



##### !end-question

##### !answer

1.
2.
3.

##### !end-answer
<optional-attributes>
### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`
