package questions

import (
	"github.com/spf13/cobra"
)

func NewGoCommand(params NewQuestionCommandParams) *cobra.Command {
	params.name = "Go"
	params.abbr = "go"
	params.maxTemplate = goTemplate
	params.minTemplate = goTemplateMin
	params.long = goLongDescription

	return createQuestionCommand(params)
}

const goLongDescription = `Go Code Snippet Challenges allow a student to write Go code directly in
Learn. The submission is evaluated against unit tests that are set up as part
of the Challenge. The student then sees the standard output from the test
runner in Learn.

The test runner automatically wraps the submission in ` + "`package main`" + `, so the
placeholder and tests blocks should not declare a package — just write the
function, any imports, and the tests.`

const goTemplate = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->
<!-- Replace everything in square brackets [] and remove brackets  -->

### !challenge

* type: code-snippet
* language: golang
* id: %s
* title: [text, a short question title]
<!-- * points: [1] (optional, the number of points for scoring as a checkpoint) -->
<!-- * topics: [goroutines, channels] (Checkpoints only, optional the topics for analyzing points) -->
<!-- * test_file: [/path/to/file.txt] (External test file, replaces 'tests' section) -->
<!-- * setup_file: [/path/to/file.txt] (External setup file, replaces 'setup' section) -->

##### !question

[markdown, your question]

##### !end-question

##### !placeholder

[the code below is the starting code in the web editor]
~~~go
// DoSomething returns true
func DoSomething() bool {
	// return true
	return false
}
~~~

##### !end-placeholder

##### !tests

[the go tests below will run against the student submission]
~~~go
import (
	"strings"
	"testing"
)

func TestDoSomething(t *testing.T) {
	if !DoSomething() {
		t.Errorf("%s should return true", strings.ToLower("DoSomething"))
	}
}
~~~

##### !end-tests

<optional-attributes>

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const goTemplateMin = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->

### !challenge

* type: code-snippet
* language: golang
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
