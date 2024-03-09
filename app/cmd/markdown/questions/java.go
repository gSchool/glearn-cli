package questions

import (
	"github.com/spf13/cobra"
)

func NewJavaCommand(params NewQuestionCommandParams) *cobra.Command {
	params.name = "Java"
	params.abbr = "ja"
	params.maxTemplate = javaTemplate
	params.minTemplate = javaTemplateMin
	params.long = javaLongDescription

	return createQuestionCommand(params)
}

const javaLongDescription = `Java Code Snippet Challenges allow a student to write code directly in Learn.
The submission is evaluated against unit tests that are set up as part of the
Challenge. The student then sees the standard output from the test runner in
Learn.`

const javaTemplate = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->
<!-- Replace everything in square brackets [] and remove brackets  -->

### !challenge

* type: code-snippet
* language: java
* id: %s
* title: [text, a short question title]
<!-- * points: [1] (optional, the number of points for scoring as a checkpoint) -->
<!-- * topics: [python, pandas] (Checkpoints only, optional the topics for analyzing points) -->
<!-- * test_file: [/path/to/file.txt] (External test file, replaces 'tests' section) -->
<!-- * setup_file: [/path/to/file.txt] (External setup file, replaces 'setup' section) -->

##### !question

[markdown, your question]

##### !end-question

##### !setup

[the code below will be added to the beginning of the student submission]
~~~java
// include any imports specific to your tests
import java.io.IOException;

// to allow student to submit simple statements, wrap the submission
//  using the !setup and !tests sections; example below
class VariableChallenge {

    public static String run() {
        // Start Student Code
~~~

##### !end-setup

##### !placeholder

~~~java
[the code below is the starting code in the web editor]
// write code that declares the string foo and sets it to "bar"
// String foo="bar";
~~~

##### !end-placeholder

##### !tests

[the test code below will be added to the end of the student submission]
~~~java
  // End Student Code
  return foo;
  }
}

// public test class name **must** be SnippetTest to match the generated file name
public class SnippetTest {

	@Test
	public void someTest() {
		assertEquals("bar", VariableChallenge.run());
	}
}
~~~

##### !end-tests

<optional-attributes>

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const javaTemplateMin = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->

### !challenge

* type: code-snippet
* language: java
* id: %s
* title:

##### !question



##### !end-question

##### !setup



##### !end-setup

##### !placeholder



##### !end-placeholder

##### !tests



##### !end-tests
<optional-attributes>
### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`
