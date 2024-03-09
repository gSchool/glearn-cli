package questions

import (
	"github.com/spf13/cobra"
)

func NewParagraphCommand(params NewQuestionCommandParams) *cobra.Command {
	params.name = "Paragraph"
	params.abbr = "pg"
	params.maxTemplate = paragraphTemplate
	params.minTemplate = paragraphTemplateMin
	params.long = paragraphLongDescription

	return createQuestionCommand(params)
}

const paragraphLongDescription = `Paragraph Challenges allow a student to submit a long free-form text answer to
a question, such as a definition of explanation. The answers to these
Challenges are not evaluated by Learn but are available for the instructor to
view.`

const paragraphTemplate = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->
<!-- Replace everything in square brackets [] and remove brackets  -->

### !challenge

* type: paragraph
* id: %s
* title: [text, a short question title]
<!-- * points: [1] (optional, the number of points for scoring as a checkpoint) -->
<!-- * topics: [python, pandas] (Checkpoints only, optional the topics for analyzing points) -->

##### !question

[markdown, your question]

##### !end-question

##### !placeholder

[text, placeholder text for input field]

##### !end-placeholder

<optional-attributes>

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const paragraphTemplateMin = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->
<!-- Replace everything in square brackets [] and remove brackets  -->

### !challenge

* type: paragraph
* id: %s
* title:

##### !question



##### !end-question
<optional-attributes>
### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`
