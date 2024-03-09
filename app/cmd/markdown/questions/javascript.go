package questions

import (
	"github.com/spf13/cobra"
)

func NewJavaScriptCommand(params NewQuestionCommandParams) *cobra.Command {
	params.name = "JavaScript"
	params.abbr = "js"
	params.maxTemplate = javascriptTemplate
	params.minTemplate = javascriptTemplateMin
	params.long = javascriptLongDescription

	return createQuestionCommand(params)
}

const javascriptLongDescription = `Code Snippet Challenges allow a student to write code directly in Learn. The
submission is evaluated against unit tests that are set up as part of the
Challenge. The student then sees the standard output from the test runner in
Learn.
`

const javascriptTemplate = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->
<!-- Replace everything in square brackets [] and remove brackets  -->

### !challenge

* type: code-snippet
* language: javascript18
* id: %s
* title: [text, a short question title]
<!-- * points: [1] (optional, the number of points for scoring as a checkpoint) -->
<!-- * topics: [python, pandas] (Checkpoints only, optional the topics for analyzing points) -->
<!-- * test_file: [/path/to/file.txt] (External test file, replaces 'tests' section) -->
<!-- * setup_file: [/path/to/file.txt] (External setup file, replaces 'setup' section) -->

##### !question

[markdown, your question]

##### !end-question

##### !placeholder

[the code below is the starting code in the web editor]
~~~js
// notes on what to return, etc
function doSomething() {
  // return true
}
~~~

##### !end-placeholder

##### !tests

[the mocha tests below will run against the student submission]
~~~js
describe('doSomething', function() {

  it("does what it is supposed to do", function() {
    expect(doSomething(), "Error message").to.deep.eq(true)
  })
})
~~~

##### !end-tests

<optional-attributes>

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const javascriptTemplateMin = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->

### !challenge

* type: code-snippet
* language: javascript18
* id: %s
* title:

##### !question



##### !end-question

##### !placeholder



##### !end-placeholder

##### !tests



##### !end-tests
<optional-attributes>
### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`
