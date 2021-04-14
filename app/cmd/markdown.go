package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var PrintTemplate bool

var Minimal bool

var markdownCmd = &cobra.Command{
	Use:     "markdown",
	Aliases: []string{"md"},
	Short:   "Copy curriculum markdown to clipboard",
	Long:    "Copy curriculum markdown to clipboard. Takes 1-2 arguments, the type of content to copy to clipboard and optionally a file to append.\n\n" + argList,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			t, err := getTemp(args[0])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if PrintTemplate {
				t.printContent()
			} else {
				t.copyContent()
			}

		} else if len(args) == 2 {
			t, err := getTemp(args[0])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if PrintTemplate {
				fmt.Println("-o flag skipped when appending...")
			}
			if err = t.appendContent(args[1]); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

		} else {
			fmt.Println(incorrectNumArgs)
			os.Exit(1)
		}

	},
}

func getTemp(command string) (temp, error) {
	t, ok := templates[command]
	if !ok {
		return temp{}, fmt.Errorf("Unknown arg '%s'. Run 'learn md --help' for options.\n", command)
	}
	return t, nil
}

type temp struct {
	Name        string
	Template    string
	MinTemplate string
	RequireId   bool
}

func (t temp) printContent() {
	template := t.Template
	if Minimal {
		template = t.MinTemplate
	}
	if t.RequireId {
		id := uuid.New().String()
		fmt.Printf(strings.ReplaceAll(template, `~~~`, "```")+"\n", id)
	} else {
		fmt.Println(template)
	}
}

func (t temp) copyContent() {
	template := t.Template
	if Minimal {
		template = t.MinTemplate
	}
	if t.RequireId {
		id := uuid.New().String()
		clipboard.WriteAll(fmt.Sprintf(strings.ReplaceAll(template, `~~~`, "```"), id))
		fmt.Println(t.Name, "generated with id:", id, "\nCopied to clipboard!")
	} else {
		clipboard.WriteAll(template)
		fmt.Println(t.Name, "copied to clipboard!")
	}
}

func (t temp) appendContent(target string) error {
	template := t.Template
	if Minimal {
		template = t.MinTemplate
	}
	if !strings.HasSuffix(target, ".md") {
		return fmt.Errorf("'%s' must have an `.md` extension to append %s content.\n", target, t.Name)
	}

	targetInfo, err := os.Stat(target)
	if err != nil {
		return fmt.Errorf("'%s' is not a file that can be appended!\n%s\n", target, err)
	}
	if targetInfo.IsDir() {
		return fmt.Errorf("'%s' is a directory, please specify a markdown file.\n", target)
	}

	f, err := os.OpenFile(target, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("Cannot open '%s'!\n%s\n", target, err)
	}
	defer f.Close()

	if t.RequireId {
		id := uuid.New().String()
		if _, err = f.WriteString(fmt.Sprintf(strings.ReplaceAll(template, `~~~`, "```"), id) + "\n"); err != nil {
			return fmt.Errorf("Cannot write to '%s'!\n%s\n", target, err)
		}
		fmt.Printf("%s appended to %s!\nid: %s\n", t.Name, target, id)
	} else {
		if _, err = f.WriteString(template + "\n"); err != nil {
			return fmt.Errorf("Cannot write to '%s'!\n%s\n", target, err)
		}
		fmt.Printf("%s appended to %s!\n", t.Name, target)
	}

	return nil
}

var templates = map[string]temp{
	"ls":              {"Lesson markdown", lessonTemplate, lessonTemplateMin, true},
	"lesson":          {"Lesson markdown", lessonTemplate, lessonTemplateMin, true},
	"cp":              {"Checkpoint markdown", checkpointTemplate, checkpointTemplateMin, true},
	"checkpoint":      {"Checkpoint markdown", checkpointTemplate, checkpointTemplateMin, true},
	"sv":              {"Survey markdown", surveyTemplate, surveyTemplateMin, true},
	"survey":          {"Survey markdown", surveyTemplate, surveyTemplateMin, true},
	"in":              {"Instructor markdown", instructorTemplate, instructorTemplateMin, true},
	"instructor":      {"Instructor markdown", instructorTemplate, instructorTemplateMin, true},
	"rs":              {"Resource markdown", resourceTemplate, resourceTemplateMin, true},
	"resource":        {"Resource markdown", resourceTemplate, resourceTemplateMin, true},
	"fh":              {"File header", fileHeaderTemplate, fileHeaderTemplateMin, true},
	"fileheader":      {"File header", fileHeaderTemplate, fileHeaderTemplateMin, true},
	"mc":              {"Multiple Choice markdown", multiplechoiceTemplate, multiplechoiceTemplateMin, true},
	"multiplechoice":  {"Multiple Choice markdown", multiplechoiceTemplate, multiplechoiceTemplateMin, true},
	"cb":              {"Checkbox markdown", checkboxTemplate, checkboxTemplateMin, true},
	"checkbox":        {"Checkbox markdown", checkboxTemplate, checkboxTemplateMin, true},
	"tl":              {"Tasklist markdown", tasklistTemplate, tasklistTemplateMin, true},
	"tasklist":        {"Tasklist markdown", tasklistTemplate, tasklistTemplateMin, true},
	"sa":              {"Short Answer markdown", shortanswerTemplate, shortanswerTemplateMin, true},
	"shortanswer":     {"Short Answer markdown", shortanswerTemplate, shortanswerTemplateMin, true},
	"nb":              {"Number markdown", numberTemplate, numberTemplateMin, true},
	"number":          {"Number markdown", numberTemplate, numberTemplateMin, true},
	"pg":              {"Paragraph markdown", paragraphTemplate, paragraphTemplateMin, true},
	"paragraph":       {"Paragraph markdown", paragraphTemplate, paragraphTemplateMin, true},
	"or":              {"Ordering markdown", orderingTemplate, orderingTemplateMin, true},
	"ordering":        {"Ordering markdown", orderingTemplate, orderingTemplateMin, true},
	"js":              {"Javascript markdown", javascriptTemplate, javascriptTemplateMin, true},
	"javascript":      {"Javascript markdown", javascriptTemplate, javascriptTemplateMin, true},
	"ja":              {"Java markdown", javaTemplate, javaTemplateMin, true},
	"java":            {"Java markdown", javaTemplate, javaTemplateMin, true},
	"py":              {"Python markdown", pythonTemplate, pythonTemplateMin, true},
	"python":          {"Python markdown", pythonTemplate, pythonTemplateMin, true},
	"sq":              {"Sql markdown", sqlTemplate, sqlTemplateMin, true},
	"sql":             {"Sql markdown", sqlTemplate, sqlTemplateMin, true},
	"cs":              {"Custom Snippet markdown", customsnippetTemplate, customsnippetTemplateMin, true},
	"customsnippet":   {"Custom Snippet markdown", customsnippetTemplate, customsnippetTemplateMin, true},
	"pr":              {"Project markdown", projectTemplate, projectTemplateMin, true},
	"project":         {"Project markdown", projectTemplate, projectTemplateMin, true},
	"tpr":             {"Testable Project markdown", testableProjectTemplate, testableProjectTemplateMin, true},
	"testableproject": {"Testable Project markdown", testableProjectTemplate, testableProjectTemplateMin, true},
	"cfy":             {"config.yaml syntax", configyamlTemplate, configyamlTemplateMin, false},
	"configyaml":      {"config.yaml syntax", configyamlTemplate, configyamlTemplateMin, false},
	"cry":             {"course.yaml syntax", courseyamlTemplate, courseyamlTemplateMin, false},
	"courseyaml":      {"course.yaml syntax", courseyamlTemplate, courseyamlTemplateMin, false},
	"callout":         {"Callout markdown", calloutTemplate, calloutTemplateMin, false},
	"co":              {"Callout markdown", calloutTemplate, calloutTemplateMin, false},
	"upload":          {"Upload markdown", uploadTemplate, uploadTemplateMin, false},
	"up":              {"Upload markdown", uploadTemplate, uploadTemplateMin, false},
}

const incorrectNumArgs = "Copy curriculum markdown to clipboard. \n\nTakes 1-2 arguments, the type of content to copy to clipboard and optionally a markdown file to append. Specify -o to print to stdout, -m for a minimal template.\n\n" + argList

const argList = `Args, full (abbreviation)--

Files:
  lesson (ls)
  checkpoint (cp)
  survey (sv)
  instructor (in)
  resource (rs)
  fileheader (fh)
Questions:
  multiplechoice (mc)
  checkbox (cb)
  tasklist (tl)
  shortanswer (sa)
  number (nb)
  paragraph (pg)
  ordering (or)
  javascript (js)
  java (ja)
  python (py)
  sql (sq)
  upload (up)
  customsnippet (cs)
  project (pr)
  testableproject (tpr)
Other Markdown:
  callout (co)
Configuration:
  configyaml (cfy)
  courseyaml (cry)`

const fileHeaderTemplate = `---
# BEGIN FILE CONFIGURATION YML HEADER >>>>>
# autoconfig.yml will use these settings. config.yml will override.
Type: Lesson # Options: Lesson, Checkpoint, Survey, Instructor, Resource
UID: %s
# DefaultVisibility: hidden # Uncomment this line to default Lesson to hidden
# MaxCheckpointSubmissions: 1 # Checkpoints only. Uncomment this line to limit the number of submissions
# TimeLimit: 60 # Checkpoints only. Uncomment this line to set a time limit in minutes
# Autoscore: true # Checkpoints only. Uncomment this line to finalize checkpoint scores without instructor review
# END FILE CONFIGURATION YML HEADER <<<<<
---`

const fileHeaderTemplateMin = `---
Type: Lesson
UID: %s
# DefaultVisibility: hidden
# MaxCheckpointSubmissions: 1
# TimeLimit: 60
# Autoscore: true
---`

const lessonTemplate = `---
# BEGIN FILE CONFIGURATION YML HEADER >>>>>
# autoconfig.yml will use these settings. config.yml will override.
Type: Lesson
UID: %s
# DefaultVisibility: hidden # Uncomment this line to default Lesson to hidden
# END FILE CONFIGURATION YML HEADER <<<<<
---

# Title

<!--Lesson content can be markdown, videos, slides, images, gifs, etc. See examples of markdown formatting -- https://galvanize-learn.zendesk.com/hc/en-us/articles/360061963154-Markdown-Formatting-->
<!--Lessons can include Challenges, which make the content interactive and give instructors visibility into student learning. See -- https://galvanize-learn.zendesk.com/hc/en-us/articles/360061964054-Galvanize-Learn-Markdown-Extensions-Challenges-and-Callouts-->
`

const lessonTemplateMin = `---
Type: Lesson
UID: %s
# DefaultVisibility: hidden
---

# Title
`

const resourceTemplate = `---
# BEGIN FILE CONFIGURATION YML HEADER >>>>>
# autoconfig.yml will use these settings. config.yml will override.
Type: Resource
UID: %s
# END FILE CONFIGURATION YML HEADER <<<<<
---

# Title

<!--A Resource can have all of the same markdown and challenges as a lesson. Resources do not appear in the left nav and don't count toward course completion.-->
`

const resourceTemplateMin = `---
Type: Resource
UID: %s
---

# Title
`

const instructorTemplate = `---
# BEGIN FILE CONFIGURATION YML HEADER >>>>>
# autoconfig.yml will use these settings. config.yml will override.
Type: Instructor
UID: %s
# END FILE CONFIGURATION YML HEADER <<<<<
---

# Title

<!--An Instructor file can have all of the same markdown and challenges as a lesson. Instructor files are only viewable by instructors. -->
`

const instructorTemplateMin = `---
Type: Instructor
UID: %s
---

# Title
`

const surveyTemplate = `---
# BEGIN FILE CONFIGURATION YML HEADER >>>>>
# autoconfig.yml will use these settings. config.yml will override.
Type: Survey
UID: %s
# DefaultVisibility: hidden # Uncomment this line to default Survey to hidden
# END FILE CONFIGURATION YML HEADER <<<<<
---

# Title

<!--A Survey can have any markdown. See examples of markdown formatting -- https://galvanize-learn.zendesk.com/hc/en-us/articles/360061963154-Markdown-Formatting-->
<!--A Survey must include include one or more Challenges, which are the survey questions a student will answer. See -- https://galvanize-learn.zendesk.com/hc/en-us/articles/360061964054-Galvanize-Learn-Markdown-Extensions-Challenges-and-Callouts-->
`

const surveyTemplateMin = `---
Type: Survey
UID: %s
# DefaultVisibility: hidden
---

# Title
`

const checkpointTemplate = `---
# BEGIN FILE CONFIGURATION YML HEADER >>>>>
# autoconfig.yml will use these settings. config.yml will override.
Type: Checkpoint
UID: %s
# DefaultVisibility: hidden # Uncomment this line to default Checkpoint to hidden
# MaxCheckpointSubmissions: 1 # Uncomment this line to limit the number of submissions
# TimeLimit: 60 # Uncomment this line to set a time limit in minutes
# Autoscore: true # Uncomment this line to finalize checkpoint scores without instructor review
# END FILE CONFIGURATION YML HEADER <<<<<
---

# Title

<!--A Checkpoint is an assessment and must include include one or more Challenges. See -- https://galvanize-learn.zendesk.com/hc/en-us/articles/360061964054-Galvanize-Learn-Markdown-Extensions-Challenges-and-Callouts-->
<!--A Checkpoint can have any markdown. See examples of markdown formatting -- https://galvanize-learn.zendesk.com/hc/en-us/articles/360061963154-Markdown-Formatting-->
`

const checkpointTemplateMin = `---
Type: Checkpoint
UID: %s
# DefaultVisibility: hidden
# MaxCheckpointSubmissions: 1
# TimeLimit: 60
# Autoscore: true
---

# Title
`

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

* [Option 1]
* [Option 2]
* [Option 3, etc]

##### !end-options

##### !answer

* [Option 2 (the correct answer)]

##### !end-answer

<!-- other optional sections -->
<!-- !hint - !end-hint (markdown, hidden, students click to view) -->
<!-- !rubric - !end-rubric (markdown, instructors can see while scoring a checkpoint) -->
<!-- !explanation - !end-explanation (markdown, students can see after answering correctly) -->

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const multiplechoiceTemplateMin = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->

### !challenge

* type: multiple-choice
* id: %s
* title:

##### !question



##### !end-question

##### !options

*
*
*

##### !end-options

##### !answer

*

##### !end-answer

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

* [Option 1]
* [Option 2]
* [Option 3, etc]

##### !end-options

##### !answer

* [Option 2]
* [Option 3 (the correct answer set)]

##### !end-answer

<!-- other optional sections -->
<!-- !hint - !end-hint (markdown, hidden, students click to view) -->
<!-- !rubric - !end-rubric (markdown, instructors can see while scoring a checkpoint) -->
<!-- !explanation - !end-explanation (markdown, students can see after answering correctly) -->

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const checkboxTemplateMin = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->

### !challenge

* type: checkbox
* id: %s
* title:

##### !question



##### !end-question

##### !options

*
*
*

##### !end-options

##### !answer

*
*

##### !end-answer

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const tasklistTemplate = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->
<!-- Replace everything in square brackets [] and remove brackets  -->

### !challenge

* type: tasklist
* id: %s
* title: [text, a short question title]
<!-- * points: [1] (optional, the number of points for scoring as a checkpoint) -->
<!-- * topics: [python, pandas] (optional the topics for analyzing points) -->

##### !question

[optional, markdown, a prompt at the top of the tasklist]

##### !end-question

##### !options

* [Task 1]
* [Task 2]
* [Task 3, etc]

##### !end-options

<!-- other optional sections -->
<!-- !hint - !end-hint (markdown, hidden, students click to view) -->
<!-- !rubric - !end-rubric (markdown, instructors can see while scoring a checkpoint) -->
<!-- !explanation - !end-explanation (markdown, students can see after answering correctly) -->

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const tasklistTemplateMin = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->

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
<!-- !hint - !end-hint (markdown, hidden, students click to view) -->
<!-- !rubric - !end-rubric (markdown, instructors can see while scoring a checkpoint) -->
<!-- !explanation - !end-explanation (markdown, students can see after answering correctly) -->

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const shortanswerTemplateMin = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->

### !challenge

* type: short-answer
* id: %s
* title:

##### !question



##### !end-question

##### !answer



##### !end-answer

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
<!-- !hint - !end-hint (markdown, hidden, students click to view) -->
<!-- !rubric - !end-rubric (markdown, instructors can see while scoring a checkpoint) -->
<!-- !explanation - !end-explanation (markdown, students can see after answering correctly) -->

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const numberTemplateMin = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->

### !challenge

* type: number
* id: %s
* title:
* decimal:

##### !question



##### !end-question

##### !answer



##### !end-answer

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
<!-- !hint - !end-hint (markdown, hidden, students click to view) -->
<!-- !rubric - !end-rubric (markdown, instructors can see while scoring a checkpoint) -->
<!-- !explanation - !end-explanation (markdown, students can see after answering correctly) -->

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

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const orderingTemplate = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->
<!-- Replace everything in square brackets [] and remove brackets  -->

### !challenge

* type: ordering
* id: %s
* title: [text, a short question title]
<!-- * points: [1] (optional, the number of points for scoring as a checkpoint) -->
<!-- * topics: [python, pandas] (optional the topics for analyzing points) -->

##### !question

[markdown, your question]

##### !end-question

##### !answer

1. [Option 1 in the correct first position, options will be randomized for students]
1. [Option 2 in the correct second position]
1. [Option 3 in the correct third position]

##### !end-answer

<!-- other optional sections -->
<!-- !hint - !end-hint (markdown, hidden, students click to view) -->
<!-- !rubric - !end-rubric (markdown, instructors can see while scoring a checkpoint) -->
<!-- !explanation - !end-explanation (markdown, students can see after answering correctly) -->

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

<!-- other optional sections -->
<!-- !hint - !end-hint (markdown, hidden, students click to view) -->
<!-- !rubric - !end-rubric (markdown, instructors can see while scoring a checkpoint) -->
<!-- !explanation - !end-explanation (markdown, students can see after answering correctly) -->

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const javascriptTemplateMin = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->

### !challenge

* type: code-snippet
* language: javascript
* id: %s
* title:

##### !question



##### !end-question

##### !placeholder



##### !end-placeholder

##### !tests



##### !end-tests



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

<!-- other optional sections -->
<!-- !hint - !end-hint (markdown, hidden, students click to view) -->
<!-- !rubric - !end-rubric (markdown, instructors can see while scoring a checkpoint) -->
<!-- !explanation - !end-explanation (markdown, students can see after answering correctly) -->

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

<!-- other optional sections -->
<!-- !hint - !end-hint (markdown, hidden, students click to view) -->
<!-- !rubric - !end-rubric (markdown, instructors can see while scoring a checkpoint) -->
<!-- !explanation - !end-explanation (markdown, students can see after answering correctly) -->

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const pythonTemplateMin = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->

### !challenge

* type: code-snippet
* language: python3.6
* id: %s
* title:

##### !question



##### !end-question

##### !placeholder



##### !end-placeholder

##### !tests



##### !end-tests

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

[the code below is the starting code in the web editor]
~~~sql
-- write a statement to select...
~~~

##### !end-placeholder

##### !tests

[the code below is the sql statement that returns the correct answer]
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
<!-- !hint - !end-hint (markdown, hidden, students click to view) -->
<!-- !rubric - !end-rubric (markdown, instructors can see while scoring a checkpoint) -->
<!-- !explanation - !end-explanation (markdown, students can see after answering correctly) -->

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const sqlTemplateMin = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->

### !challenge

* type: code-snippet
* language: sql
* id: %s
* title:
* data_path: /

##### !question



##### !end-question

##### !placeholder



##### !end-placeholder

##### !tests



##### !end-tests

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const customsnippetTemplate = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->
<!-- Replace everything in square brackets [] and remove brackets  -->

### !challenge

* type: custom-snippet
* language: [text, one of: csharp, html, java, javascript, json, markdown, python, or sql]
* id: %s
* title: [text, a short question title]
* docker_directory_path: /[text, the path to the folder with the Docker setup]
<!-- * points: [1] (optional, the number of points for scoring as a checkpoint) -->
<!-- * topics: [python, pandas] (optional the topics for analyzing points) -->

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

<!-- other optional sections -->
<!-- !hint - !end-hint (markdown, hidden, students click to view) -->
<!-- !rubric - !end-rubric (markdown, instructors can see while scoring a checkpoint) -->
<!-- !explanation - !end-explanation (markdown, students can see after answering correctly) -->

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const customsnippetTemplateMin = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->

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
<!-- !hint - !end-hint (markdown, hidden, students click to view) -->
<!-- !rubric - !end-rubric (markdown, instructors can see while scoring a checkpoint) -->
<!-- !explanation - !end-explanation (markdown, students can see after answering correctly) -->

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const projectTemplateMin = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->

### !challenge

* type: project
* id: %s
* title:

##### !question



##### !end-question

##### !placeholder



##### !end-placeholder

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
<!-- !hint - !end-hint (markdown, hidden, students click to view) -->
<!-- !rubric - !end-rubric (markdown, instructors can see while scoring a checkpoint) -->
<!-- !explanation - !end-explanation (markdown, students can see after answering correctly) -->

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const testableProjectTemplateMin = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->

### !challenge

* type: testable-project
* id: %s
* title:
* upstream:
* validate_fork: true

##### !question



##### !end-question

##### !placeholder



##### !end-placeholder

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
        Path: /[folder/file.md]`

const configyamlTemplateMin = `---
Standards:
  - Title:
    UID:
    Description:
    SuccessCriteria:
      -
    ContentFiles:
      - Type: Lesson
        UID:
        Path: /
      - Type: Checkpoint
        UID:
        Path: /`

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

const courseyamlTemplateMin = `---
# DefaultUnitVisibility: hidden
Course:
  - Section:
    Repos:
      - URL:`

const calloutTemplate = `<!-- available callout types: info, success, warning, danger, secondary  -->
### !callout-info

## title

body

### !end-callout`

const calloutTemplateMin = `### !callout-info

## title



### !end-callout`

const uploadTemplate = `### !challenge

* type: upload
* id: %s
* title:
* points:

##### !question
[markdown, your question]
##### !end-question

### !end-challenge`

const uploadTemplateMin = `### !challenge

* type: upload
* id:
* title:
* points:

##### !question

##### !end-question

### !end-challenge`
