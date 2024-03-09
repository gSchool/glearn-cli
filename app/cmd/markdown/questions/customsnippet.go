package questions

import (
	"github.com/spf13/cobra"
)

func NewCustomSnippetCommand(params NewQuestionCommandParams) *cobra.Command {
	params.name = "Custom Snippet"
	params.abbr = "cs"
	params.maxTemplate = customSnippetTemplate
	params.minTemplate = customSnippetTemplateMin
	params.long = customSnippetLongDescription

	return createQuestionCommand(params)
}

const customSnippetLongDescription = `The Custom Snippet challenges allow a student to write code directly in Learn
and see the standard output from the test runner. They provide a student
experience similar to the JavaScript and Python Code Snippet Challenges.
However, the Curriculum Developer controls the Dockerfile and has all of the
flexibility of testable-project challenges.

NOTE: All custom-snippet Docker images must have the bash shell available.`

const customSnippetTemplate = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->
<!-- Replace everything in square brackets [] and remove brackets  -->

### !challenge

* type: custom-snippet
* language: [text, one of: csharp, html, java, javascript, json, markdown, python, ruby, or sql]
* id: %s
* title: [text, a short question title]
* docker_directory_path: /[text, the path to the folder with the Docker setup]
<!-- * points: [1] (optional, the number of points for scoring as a checkpoint) -->
<!-- * topics: [python, pandas] (Checkpoints only, optional the topics for analyzing points) -->

##### !question

[markdown, your question]

##### !end-question

##### !placeholder

[the code below is the starting code in the web editor]
~~~
function doSomething() {
}
~~~

##### !end-placeholder

<optional-attributes>

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const customSnippetTemplateMin = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->

### !challenge

* type: custom-snippet
* language: [text, one of: csharp, html, java, javascript, json, markdown, python, or sql]
* id: %s
* title:
* docker_directory_path: /

##### !question



##### !end-question

##### !placeholder



##### !end-placeholder
<optional-attributes>
### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`
