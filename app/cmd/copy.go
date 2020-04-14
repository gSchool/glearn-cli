package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/atotto/clipboard"
	"github.com/google/uuid"
)

var copyCmd = &cobra.Command{
	Use:   "copy",
	Aliases: []string{"cp"},
	Short: "Copy curriculum markdown to clipboard",
	Long: `Copy curriculum markdown to clipboard. Takes one argument, the type of content to copy to clipboard.

File Args (full / abbreviation):
  lesson / ls
  checkpoint / cp
Question Args:
  multiplechoice / mc
  checkbox / cb
  shortanswer / sa
  number / nb
  paragraph / pg
  javascript / js
  python / py
  sql / sql
  project / pj
  testableproject / tpj
Configuration Args--
  configyaml / cfy
  courseyaml / cry
	`,
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println(incorrectNumArgs)
			os.Exit(1)
		}
		id := uuid.New()
		switch args[0] {
		case "ls", "lesson":
			clipboard.WriteAll(lessonTemplate)
			fmt.Println("Copied lesson markdown to clipboard.")
		case "cp", "checkpoint":
			clipboard.WriteAll(lessonTemplate) // update
			fmt.Println("Copied checkpoint markdown to clipboard.")
		case "mc", "multiplechoice":
			clipboard.WriteAll(fmt.Sprintf(multiplechoiceTemplate, id.String()))
			fmt.Println("Copied multiple choice markdown to clipboard -- " + id.String())
		case "cb", "checkbox":
			clipboard.WriteAll(fmt.Sprintf(checkboxTemplate, id.String()))
			fmt.Println("Copied checkbox markdown to clipboard -- " + id.String())
		default:
			fmt.Println("Unknown arg " + args[0] + ". Run 'learn cp --help' for options.")
		}
	},
}


const incorrectNumArgs = `Incorrect number of args. Takes one argument, the type of content to copy to clipboard.

File Args (full / abbreviation):
  lesson / ls
  checkpoint / cp
Question Args:
  multiplechoice / mc
  checkbox / cb
  shortanswer / sa
  number / nb
  paragraph / pg
  javascript / js
  python / py
  sql / sql
  project / pj
  testableproject / tpj
Configuration Args--
  configyaml / cfy
  courseyaml / cry`

const lessonTemplate =`# Title

## Learning Objectives

By the end of this lesson you will be able to:

* First Objective
* [at least one]
* [no more than four]

## Lesson Content

[Include a rationale as the first few sentences/paragraph if you feel the lesson requires significant motivation or context.]

## Challenges

[Each lesson must have one or more challenges. These challenge can be spread out in between content, or can be at the end of the lesson]`

const multiplechoiceTemplate =`<!-- ▼▼▼▼▼▼▼▼▼▼▼▼▼▼▼▼▼▼▼▼ BEGIN CHALLENGE ▼▼▼▼▼▼▼▼▼▼▼▼▼▼▼▼▼▼▼▼ -->
<!-- Replace everything in square brackets [] and remove brackets -->

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

[* Option 1
* Option 2
* Option 3]

##### !end-options

##### !answer

[* Option 2]

##### !end-answer

<!-- other optional sections -->
<!-- !hint - !end-hint (markdown, users can see after a failed attempt) -->
<!-- !rubric - !end-rubric (markdown, instructors can see while scoring a checkpoint) -->
<!-- !explanation - !end-explanation (markdown, students can see after answering correctly) -->

### !end-challenge

<!-- ▲▲▲▲▲▲▲▲▲▲▲▲▲▲▲▲▲▲▲▲ END CHALLENGE ▲▲▲▲▲▲▲▲▲▲▲▲▲▲▲▲▲▲▲▲ -->`

const checkboxTemplate =`<!-- ▼▼▼▼▼▼▼▼▼▼▼▼▼▼▼▼▼▼▼▼ BEGIN CHALLENGE ▼▼▼▼▼▼▼▼▼▼▼▼▼▼▼▼▼▼▼▼ -->
<!-- Replace everything in square brackets [] and remove brackets -->

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

[* Option 1
* Option 2
* Option 3]

##### !end-options

##### !answer

[* Option 2
* Option 3]

##### !end-answer

<!-- other optional sections -->
<!-- !hint - !end-hint (markdown, users can see after a failed attempt) -->
<!-- !rubric - !end-rubric (markdown, instructors can see while scoring a checkpoint) -->
<!-- !explanation - !end-explanation (markdown, students can see after answering correctly) -->

### !end-challenge

<!-- ▲▲▲▲▲▲▲▲▲▲▲▲▲▲▲▲▲▲▲▲ END CHALLENGE ▲▲▲▲▲▲▲▲▲▲▲▲▲▲▲▲▲▲▲▲ -->`
