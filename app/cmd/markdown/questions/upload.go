package questions

import (
	"github.com/spf13/cobra"
)

func NewUploadCommand(params NewQuestionCommandParams) *cobra.Command {
	params.name = "Upload"
	params.abbr = "up"
	params.maxTemplate = uploadTemplate
	params.minTemplate = uploadTemplateMin
	params.long = uploadLongDescription

	return createQuestionCommand(params)
}

const uploadLongDescription = `Upload Challenges allow a student to submit a file in response to a question.
The instructor can then access the uploaded file for scoring.`

const uploadTemplate = `### !challenge

* type: upload
* id: %s
* title:
<!-- * points: [1] (optional, the number of points for scoring as a checkpoint) -->
<!-- * topics: [python, pandas] (Checkpoints only, optional the topics for analyzing points) -->

##### !question
[markdown, your question]
##### !end-question

<optional-attributes>

### !end-challenge`

const uploadTemplateMin = `### !challenge

* type: upload
* id: %s
* title:

##### !question

##### !end-question
<optional-attributes>
### !end-challenge`
