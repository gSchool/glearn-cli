package questions

import (
	"github.com/spf13/cobra"
)

func NewRubyCommand(params NewQuestionCommandParams) *cobra.Command {
	params.name = "Ruby"
	params.abbr = "rb"
	params.maxTemplate = rubyTemplate
	params.minTemplate = rubyTemplateMin
	params.long = rubyLongDescription

	return createQuestionCommand(params)
}

const rubyLongDescription = `Ruby Code Snippet Challenges allow a student to write Ruby code directly in
Learn. The submission is evaluated against unit tests that are set up as part
of the Challenge. The student then sees the standard output from the test
runner in Learn.`

const rubyTemplate = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->
<!-- Replace everything in square brackets [] and remove brackets  -->

### !challenge

* type: code-snippet
* language: ruby3
* id: %s
* title: [text, a short question title]
<!-- * points: [1] (optional, the number of points for scoring as a checkpoint) -->
<!-- * topics: [block, proc, lambda] (Checkpoints only, optional the topics for analyzing points) -->
<!-- * test_file: [/path/to/file.txt] (External test file, replaces 'tests' section) -->
<!-- * setup_file: [/path/to/file.txt] (External setup file, replaces 'setup' section) -->

##### !question

[markdown, your question]

##### !end-question

##### !placeholder

~~~ruby
class Foo
  def initialize
  end

  def truthy
    return true
  end
end
~~~

##### !end-placeholder

##### !tests

~~~ruby
describe Foo do
  it "does something" do
    f = Foo.new
    expect(f.truthy).to eq true
  end
end
~~~

##### !end-tests

<optional-attributes>

### !end-challenge

<!-- ======================= END CHALLENGE ======================= -->`

const rubyTemplateMin = `<!-- >>>>>>>>>>>>>>>>>>>>>> BEGIN CHALLENGE >>>>>>>>>>>>>>>>>>>>>> -->

### !challenge

* type: code-snippet
* language: ruby3
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
