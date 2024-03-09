package questions

import (
	"github.com/spf13/cobra"
)

func NewTaskListCommand(params NewQuestionCommandParams) *cobra.Command {
	params.name = "Task List"
	params.abbr = "tl"
	params.maxTemplate = taskListTemplate
	params.minTemplate = taskListTemplateMin
	params.long = taskListLongDescription

	return createQuestionCommand(params)
}

const taskListLongDescription = `Tasklist Challenges allow you to present tasks the student needs to complete.
Students check off the tasks as they complete them, and the challenge evaluates
as correct when all tasks are completed.`

const taskListTemplate = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->
<!-- Replace everything in square brackets [] and remove brackets  -->

### !challenge

* type: tasklist
* id: %s
* title: [text, a short question title]
<!-- * points: [1] (optional, the number of points for scoring as a checkpoint) -->
<!-- * topics: [python, pandas] (Checkpoints only, optional the topics for analyzing points) -->

##### !question

[optional, markdown, a prompt at the top of the tasklist]

##### !end-question

##### !options

* [Task 1]
* [Task 2]
* [Task 3, etc]

##### !end-options

<optional-attributes>

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const taskListTemplateMin = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->

### !challenge

* type: tasklist
* id: %s
* title:

##### !question



##### !end-question

##### !options

*
*
*

##### !end-options
<optional-attributes>
### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`
