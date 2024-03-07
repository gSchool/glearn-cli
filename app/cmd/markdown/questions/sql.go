package questions

import (
	"github.com/spf13/cobra"
)

func NewSqlCommand(params NewQuestionCommandParams) *cobra.Command {
	params.name = "SQL"
	params.abbr = "sq"
	params.maxTemplate = sqlTemplate
	params.minTemplate = sqlTemplateMin
	params.long = sqlLongDescription

	return createQuestionCommand(params)
}

const sqlLongDescription = `SQL Snippet Challenges allow a student to write SQL SELECT queries against a
database defined by a curriculum developer. The submission is evaluated against
a supplied test query. The student sees a truncated set of rows as if their
query was run from the psql interpreter, and the challenge is graded against
row and column matches, along with data matches. ORDER BY clauses will be
respected if they are supplied in the test query.`

const sqlTemplate = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->
<!-- Replace everything in square brackets [] and remove brackets  -->

### !challenge

* type: code-snippet
* language: sql
* id: %s
* title: [text, a short question title]
* data_path: /[text, the path to the folder with the .sql file]
<!-- * points: [1] (optional, the number of points for scoring as a checkpoint) -->
<!-- * topics: [python, pandas] (Checkpoints only, optional the topics for analyzing points) -->
<!-- * test_file: [/path/to/file.txt] (External test file, replaces 'tests' section) -->

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

<optional-attributes>

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
<optional-attributes>
### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`
