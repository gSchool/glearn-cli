package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var markdownCmd = &cobra.Command{
	Use:     "markdown",
	Aliases: []string{"md"},
	Short:   "Copy curriculum markdown to clipboard",
	Long:    "Copy curriculum markdown to clipboard. Takes one argument, the type of content to copy to clipboard.\n\n" + argList,
	Args:    cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println(incorrectNumArgs)
			os.Exit(1)
		}
		id := uuid.New().String()
		switch args[0] {
		case "ls", "lesson":
			clipboard.WriteAll(lessonTemplate)
			fmt.Println("Lesson markdown generated\nCopied to clipboard!")
		case "mc", "multiplechoice":
			clipboard.WriteAll(fmt.Sprintf(strings.ReplaceAll(multiplechoiceTemplate, `~~~`, "```"), id))
			fmt.Println("Multiple choice markdown generated with id:", id, "\nCopied to clipboard!")
		case "cb", "checkbox":
			clipboard.WriteAll(fmt.Sprintf(strings.ReplaceAll(checkboxTemplate, `~~~`, "```"), id))
			fmt.Println("Checkbox markdown generated with id:", id, "\nCopied to clipboard!")
		case "sa", "shortanswer":
			clipboard.WriteAll(fmt.Sprintf(strings.ReplaceAll(shortanswerTemplate, `~~~`, "```"), id))
			fmt.Println("Short answer markdown generated with id:", id, "\nCopied to clipboard!")
		case "nb", "number":
			clipboard.WriteAll(fmt.Sprintf(strings.ReplaceAll(numberTemplate, `~~~`, "```"), id))
			fmt.Println("Number markdown generated with id:", id, "\nCopied to clipboard!")
		case "pg", "paragraph":
			clipboard.WriteAll(fmt.Sprintf(strings.ReplaceAll(paragraphTemplate, `~~~`, "```"), id))
			fmt.Println("Paragraph markdown generated with id:", id, "\nCopied to clipboard!")
		case "js", "javascript":
			clipboard.WriteAll(fmt.Sprintf(strings.ReplaceAll(javascriptTemplate, `~~~`, "```"), id))
			fmt.Println("Javascript markdown generated with id:", id, "\nCopied to clipboard!")
		case "ja", "java":
			clipboard.WriteAll(fmt.Sprintf(strings.ReplaceAll(javaTemplate, `~~~`, "```"), id))
			fmt.Println("Java markdown generated with id:", id, "\nCopied to clipboard!")
		case "py", "python":
			clipboard.WriteAll(fmt.Sprintf(strings.ReplaceAll(pythonTemplate, `~~~`, "```"), id))
			fmt.Println("Python markdown generated with id:", id, "\nCopied to clipboard!")
		case "sq", "sql":
			clipboard.WriteAll(fmt.Sprintf(strings.ReplaceAll(sqlTemplate, `~~~`, "```"), id))
			fmt.Println("Sql markdown generated with id:", id, "\nCopied to clipboard!")
		case "pr", "project":
			clipboard.WriteAll(fmt.Sprintf(strings.ReplaceAll(projectTemplate, `~~~`, "```"), id))
			fmt.Println("Project markdown generated with id:", id, "\nCopied to clipboard!")
		case "tpr", "testableproject":
			clipboard.WriteAll(fmt.Sprintf(strings.ReplaceAll(testableProjectTemplate, `~~~`, "```"), id))
			fmt.Println("Testable Project markdown generated with id:", id, "\nCopied to clipboard!")
		case "cfy", "configyaml":
			clipboard.WriteAll(fmt.Sprintf(configyamlTemplate))
			fmt.Println("Config.yaml generated\nCopied to clipboard!")
		case "cry", "courseyaml":
			clipboard.WriteAll(fmt.Sprintf(courseyamlTemplate))
			fmt.Println("Course.yaml generated\nCopied to clipboard!")
		default:
			fmt.Println("Unknown arg " + args[0] + ". Run 'learn cp --help' for options.")
		}
	},
}

const incorrectNumArgs = "Incorrect number of args. Takes one argument, the type of content to copy to clipboard.\n\n" + argList

const argList = `Args, full (abbreviation)--

Files:
  lesson (ls)
Questions:
  multiplechoice (mc)
  checkbox (cb)
  shortanswer (sa)
  number (nb)
  paragraph (pg)
  javascript (js)
  java (ja)
  python (py)
  sql (sq)
  project (pr)
  testableproject (tpr)
Configuration:
  configyaml (cfy)
  courseyaml (cry)`

const lessonTemplate = `# Title

## Learning Objectives

By the end of this lesson you will be able to:

* First Objective
* [at least one]
* [no more than four]

## Lesson Content

[Can be written content, videos, slides, images, gifs, etc. Think about including a rationale as the first few sentences/paragraph if you feel the lesson requires significant motivation or context.]

## Challenges

[It's recommended that each lesson has at least one challenge. Challenges make the content interactive and give instructors visibility into student learning. These challenge can be spread out in between content, or can be at the end of the lesson]`

const multiplechoiceTemplate = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->
<!-- Replace everything in square brackets [] and remove brackets  -->

### !challenge

* type: multiple-choice
* id: %s
* title: [text, a short question title]
<!-- * points: [1] (optional, the number of points for scoring as a checkpoint) -->
<!-- * topics: [python, pandas] (optional the topics for analyzing points) -->

##### !question

[markdown, your question]

##### !end-question

##### !options

[markdown, replace the list below with the question options]
* Option 1
* Option 2
* Option 3

##### !end-options

##### !answer

[markdown, replace the item below with the correct answer]
* Option 2

##### !end-answer

<!-- other optional sections -->
<!-- !hint - !end-hint (markdown, users can see after a failed attempt) -->
<!-- !rubric - !end-rubric (markdown, instructors can see while scoring a checkpoint) -->
<!-- !explanation - !end-explanation (markdown, students can see after answering correctly) -->

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const checkboxTemplate = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->
<!-- Replace everything in square brackets [] and remove brackets  -->

### !challenge

* type: checkbox
* id: %s
* title: [text, a short question title]
<!-- * points: [1] (optional, the number of points for scoring as a checkpoint) -->
<!-- * topics: [python, pandas] (optional the topics for analyzing points) -->

##### !question

[markdown, your question]

##### !end-question

##### !options

[markdown, replace the list below with the question options]
* Option 1
* Option 2
* Option 3

##### !end-options

##### !answer

[markdown, replace the list below with the correct set of options]
* Option 2
* Option 3

##### !end-answer

<!-- other optional sections -->
<!-- !hint - !end-hint (markdown, users can see after a failed attempt) -->
<!-- !rubric - !end-rubric (markdown, instructors can see while scoring a checkpoint) -->
<!-- !explanation - !end-explanation (markdown, students can see after answering correctly) -->

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const shortanswerTemplate = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->
<!-- Replace everything in square brackets [] and remove brackets  -->

### !challenge

* type: short-answer
* id: %s
* title: [text, a short question title]
<!-- * points: [1] (optional, the number of points for scoring as a checkpoint) -->
<!-- * topics: [python, pandas] (optional the topics for analyzing points) -->

##### !question

[markdown, your question]

##### !end-question

##### !placeholder

[text, placeholder text for the input field]

##### !end-placeholder

##### !answer

[text or regex, the answer (if regex wrap in /)]

##### !end-answer

<!-- other optional sections -->
<!-- !hint - !end-hint (markdown, users can see after a failed attempt) -->
<!-- !rubric - !end-rubric (markdown, instructors can see while scoring a checkpoint) -->
<!-- !explanation - !end-explanation (markdown, students can see after answering correctly) -->

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const numberTemplate = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->
<!-- Replace everything in square brackets [] and remove brackets  -->

### !challenge

* type: number
* id: %s
* title: [text, a short question title]
* decimal: [optional number, decimal points to user for answer evaluation]
<!-- * points: [1] (optional, the number of points for scoring as a checkpoint) -->
<!-- * topics: [python, pandas] (optional the topics for analyzing points) -->

##### !question

[markdown, your question]

##### !end-question

##### !placeholder

[text, placeholder text for input field]

##### !end-placeholder

##### !answer

[number, the correct answer]

##### !end-answer

<!-- other optional sections -->
<!-- !hint - !end-hint (markdown, users can see after a failed attempt) -->
<!-- !rubric - !end-rubric (markdown, instructors can see while scoring a checkpoint) -->
<!-- !explanation - !end-explanation (markdown, students can see after answering correctly) -->

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const paragraphTemplate = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->
<!-- Replace everything in square brackets [] and remove brackets  -->

### !challenge

* type: paragraph
* id: %s
* title: [text, a short question title]
<!-- * points: [1] (optional, the number of points for scoring as a checkpoint) -->
<!-- * topics: [python, pandas] (optional the topics for analyzing points) -->

##### !question

[markdown, your question]

##### !end-question

##### !placeholder

[text, placeholder text for input field]

##### !end-placeholder

<!-- other optional sections -->
<!-- !hint - !end-hint (markdown, users can see after a failed attempt) -->
<!-- !rubric - !end-rubric (markdown, instructors can see while scoring a checkpoint) -->
<!-- !explanation - !end-explanation (markdown, students can see after answering correctly) -->

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const javascriptTemplate = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->
<!-- Replace everything in square brackets [] and remove brackets  -->

### !challenge

* type: code-snippet
* language: javascript
* id: %s
* title: [text, a short question title]
<!-- * points: [1] (optional, the number of points for scoring as a checkpoint) -->
<!-- * topics: [python, pandas] (optional the topics for analyzing points) -->

##### !question

[markdown, your question]

##### !end-question

##### !placeholder

[javascript, replace the block below with placeholder text for the code editor]

~~~js
// notes on what to return, etc
function doSomething() {
  // return true
}
~~~

##### !end-placeholder

##### !tests

[javascript, replace the block below with mocha tests the student code will run against]

~~~js
describe('doSomething', function() {

  it("does what it is suposed to do", function() {
    expect(doSomething(), "Error message").to.deep.eq(true)
  })
})
~~~

##### !end-tests

<!-- other optional sections -->
<!-- !hint - !end-hint (markdown, users can see after a failed attempt) -->
<!-- !rubric - !end-rubric (markdown, instructors can see while scoring a checkpoint) -->
<!-- !explanation - !end-explanation (markdown, students can see after answering correctly) -->

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const javaTemplate = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->
<!-- Replace everything in square brackets [] and remove brackets  -->

### !challenge

* type: code-snippet
* language: java
* id: %s
* title: [text, a short question title]
<!-- * points: [1] (optional, the number of points for scoring as a checkpoint) -->
<!-- * topics: [python, pandas] (optional the topics for analyzing points) -->

##### !question

[markdown, your question]

##### !end-question

##### !setup

[java, replace the block below with any code that should be added to the beginning of the student submission]
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

[java, replace the block below with placeholder text for the code editor]

~~~java
// write code that declares the string foo and sets it to "bar"
// String foo="bar";
~~~

##### !end-placeholder

##### !tests

[java, replace the block below with test code that should be added to the end of the student submission]

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

<!-- other optional sections -->
<!-- !hint - !end-hint (markdown, users can see after a failed attempt) -->
<!-- !rubric - !end-rubric (markdown, instructors can see while scoring a checkpoint) -->
<!-- !explanation - !end-explanation (markdown, students can see after answering correctly) -->

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const pythonTemplate = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->
<!-- Replace everything in square brackets [] and remove brackets  -->

### !challenge

* type: code-snippet
* language: python3.6
* id: %s
* title: [text, a short question title]
<!-- * points: [1] (optional, the number of points for scoring as a checkpoint) -->
<!-- * topics: [python, pandas] (optional the topics for analyzing points) -->

##### !question

[markdown, your question]

##### !end-question

##### !placeholder

[python, replace the block below with placeholder text for the code editor]

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

[python, replace the block below with unit tests the student code will run against]

~~~py
import unittest
import main as p
import numpy as np

class TestPython1(unittest.TestCase):
  def test_one(self):
    self.assertEqual(1,p.doSomething())
~~~

##### !end-tests

<!-- other optional sections -->
<!-- !hint - !end-hint (markdown, users can see after a failed attempt) -->
<!-- !rubric - !end-rubric (markdown, instructors can see while scoring a checkpoint) -->
<!-- !explanation - !end-explanation (markdown, students can see after answering correctly) -->

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const sqlTemplate = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->
<!-- Replace everything in square brackets [] and remove brackets  -->

### !challenge

* type: code-snippet
* language: sql
* id: %s
* title: [text, a short question title]
* data_path: /[text, the path to the folder with the .sql file]
<!-- * points: [1] (optional, the number of points for scoring as a checkpoint) -->
<!-- * topics: [python, pandas] (optional the topics for analyzing points) -->

##### !question

[markdown, your question]

##### !end-question

##### !placeholder

[sql, replace the block below with placeholder text for the code editor]

~~~sql
-- write a statement to select...
~~~

##### !end-placeholder

##### !tests

[sql, replace the block below with a sql statement to returns the correct answer]

~~~sql
SELECT these
FROM that
JOIN other
WHERE this
GROUP BY logic
ORDER BY something
~~~

##### !end-tests

<!-- other optional sections -->
<!-- !hint - !end-hint (markdown, users can see after a failed attempt) -->
<!-- !rubric - !end-rubric (markdown, instructors can see while scoring a checkpoint) -->
<!-- !explanation - !end-explanation (markdown, students can see after answering correctly) -->

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const projectTemplate = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->
<!-- Replace everything in square brackets [] and remove brackets  -->

### !challenge

* type: project
* id: %s
* title: [text, a short question title]
<!-- * points: [1] (optional, the number of points for scoring as a checkpoint) -->
<!-- * topics: [python, pandas] (optional the topics for analyzing points) -->

##### !question

[markdown, your question]

##### !end-question

##### !placeholder

[text, placeholder text for the input field]

##### !end-placeholder

<!-- other optional sections -->
<!-- !hint - !end-hint (markdown, users can see after a failed attempt) -->
<!-- !rubric - !end-rubric (markdown, instructors can see while scoring a checkpoint) -->
<!-- !explanation - !end-explanation (markdown, students can see after answering correctly) -->

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const testableProjectTemplate = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->
<!-- Replace everything in square brackets [] and remove brackets  -->

### !challenge

* type: testable-project
* id: %s
* title: [text, a short question title]
* upstream: [URL, the upstream repo URL like https://github.com/gSchool/js-native-array-methods/]
* validate_fork: true
<!-- * points: [1] (optional, the number of points for scoring as a checkpoint) -->
<!-- * topics: [python, pandas] (optional the topics for analyzing points) -->

##### !question

[markdown, your question]

##### !end-question

##### !placeholder

[text, placeholder text for the input field]

##### !end-placeholder

<!-- other optional sections -->
<!-- !hint - !end-hint (markdown, users can see after a failed attempt) -->
<!-- !rubric - !end-rubric (markdown, instructors can see while scoring a checkpoint) -->
<!-- !explanation - !end-explanation (markdown, students can see after answering correctly) -->

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const configyamlTemplate = `# Config.yaml specifies the content and ordering within a curriculum block repo
#
# Supported Fields
# ==========================
# Standards -- (Standards = Units). An array of Units for a block
# Standard.Title -- The Unit title that shows up on the curriculum overview
# Standard.UID -- A unique ID for the Unit.
# Standard.Description -- The Unit description that shows up on the curriculum overview
# Standard.SuccessCriteria -- An array of success criteria that can be viewed when scoring the checkpoint in a Unit.
# Standard.ContentFiles -- An array of Lessons and (optional) Checkpoint in a Unit.
# Standard.ContentFiles.Type -- 'Lesson' or 'Checkpoint'
# Standard.ContentFiles.UID -- A unique ID for the lesson or checkpoint.
# Standard.ContentFiles.Path -- The absolute path to the Lesson, starting with /.
# Standard.ContentFiles.DefaultVisibility -- (optional) Set to 'hidden' to hide when a course first starts.
# Standard.ContentFiles.Autoscore -- (optional, for Checkpoints only) submit checkpoint scores without review
# Standard.ContentFiles.MaxCheckpointSubmissions -- (optional, for Checkpoints only) limit the number of submissions
# Standard.ContentFiles.TimeLimit -- (optional, for Checkpoints only) the time limit in minutes
#
# Instructions
# ==========================
# Replace everything in square brackets [] and remove brackets
# Add all other Standards, Lessons, and Checkpoints following the pattern below
# All UIDs must be unique within a repo. You can use a uuidgen plugin.

---
Standards:
  - Title: [The unit name]
    UID: [unique-id]
    Description: [The Standard text]
    SuccessCriteria:
      - [The first success criteria]
    ContentFiles:
      - Type: Lesson
        UID: [unique-id]
        Path: /[folder/file.md]
      - Type: Checkpoint
        UID: [unique-id]
        Path: /[folder/file.md]
`

const courseyamlTemplate = `# Course.yaml files specify the grouping and ordering of repos that define a course.
#
# Supported Fields
# ===================
# DefaultUnitVisibility -- (optional) set to 'hidden' to hide all units when a course first starts.
# Course -- The top level array containing the sections of a course
# Course.Section -- An array contining a single array of repos. Content in the same section is grouped together on curriculum homepage.
# Course.Repos --  An array containing block repos that have been published in Learn.
# Course.Repos.URL -- The URL to a block repo that has been published in Learn.
#
# Instructions
# ==========================
# Replace everything in square brackets [] and remove brackets
# Add all other Sections and Repos following the pattern below
# All UIDs must be unique within a repo. You can use a uuidgen plugin.

---
# DefaultUnitVisibility: hidden
Course:
  - Section: [Section name]
    Repos:
      - URL: https://github.com/gSchool/[Repo name]`
