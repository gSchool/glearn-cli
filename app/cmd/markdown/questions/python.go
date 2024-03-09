package questions

import (
	"github.com/spf13/cobra"
)

func NewPythonCommand(params NewQuestionCommandParams) *cobra.Command {
	params.name = "Python"
	params.abbr = "py"
	params.maxTemplate = pythonTemplate
	params.minTemplate = pythonTemplateMin
	params.long = pythonLongDescription

	return createQuestionCommand(params)
}

const pythonLongDescription = `Python Code Snippet Challenges allow a student to write Python directly in
Learn. The submission is evaluated against unit tests that are set up as part
of the Challenge. The student then sees the standard output from the test
runner in Learn.

Python snippet challenges on Checkpoints can be worth more than one point. In
that case, Learn will automatically assign points based on the ratio of
passing/failing tests. Thus, all tests are equally weighted, and partial credit
rounds down.

Learn calculates a percentage based on passing tests, then applies that
percentage to possible points on the challenge. So if there are 5 points and 5
tests, a user gets one wrong, they get 4 points. If 10 tests 5 points, a user
gets one wrong, they get 4 points. If they get two wrong, they get 4 points,
etc.

Note that partial credit currently works only for Python snippets, which use
pytest. Partial credit does not work for custom snippets that use a different
testing library to grade Python code.`

const pythonTemplate = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->
<!-- Replace everything in square brackets [] and remove brackets  -->

### !challenge

* type: code-snippet
* language: python3.11
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
~~~py
def doSomething():
  '''
  INPUT: 2 dimensional numpy array
  OUTPUT: boolean
  Return true
  '''
#   return 1
~~~

##### !end-placeholder

##### !tests

[the unit tests below will run against the student submission]
~~~py
import unittest
import main as p
import numpy as np

class TestPython1(unittest.TestCase):
  def test_one(self):
    self.assertEqual(1,p.doSomething())
~~~

##### !end-tests

<optional-attributes>

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const pythonTemplateMin = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->

### !challenge

* type: code-snippet
* language: python3.11
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
