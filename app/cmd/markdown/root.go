package markdown

import (
	"fmt"
	"os"

	"github.com/gSchool/glearn-cli/app/cmd/markdown/files"
	"github.com/gSchool/glearn-cli/app/cmd/markdown/others"
	"github.com/gSchool/glearn-cli/app/cmd/markdown/questions"
	"github.com/gSchool/glearn-cli/app/cmd/markdown/templates"
	"github.com/gSchool/glearn-cli/app/cmd/markdown/yaml"
	"github.com/spf13/cobra"
)

var markdownCmd = &cobra.Command{
	Use:     "markdown",
	Aliases: []string{"md"},
	Short:   "Generate markdown and YAML for Learn",
	Long:    "Generate markdown for Learn content and YAML for Learn configuration\n",
}

func NewMarkdownCommand() *cobra.Command {
	markdownCmd.PersistentFlags().BoolP("out", "o", false, "Prints the template to stdout")
	markdownCmd.PersistentFlags().BoolP("min", "m", false, "Uses a terse, minimal version of the template")

	markdownCmd.AddGroup(
		&cobra.Group{ID: "files", Title: "File Creation Commands"},
		&cobra.Group{ID: "questions", Title: "Question Scaffolding Commands"},
		&cobra.Group{ID: "others", Title: "Other Commands"},
		&cobra.Group{ID: "yaml", Title: "YAML Configuration Scaffolding Commands"},
	)

	addFileCommands(markdownCmd)
	addOtherCommands(markdownCmd)
	addYamlCommands(markdownCmd)
	addQuestionCommands(markdownCmd)

	return markdownCmd
}

func addQuestionCommands(c *cobra.Command) {
	params := questions.NewQuestionCommandParams{
		RunCallback: runMarkdownSubcommand,
		Validator:   zeroOrOneFileNames,
	}

	questionCommands := []questions.NewQuestionCommand{
		questions.NewMultipleChoiceCommand,
		questions.NewCheckBoxCommand,
		questions.NewTaskListCommand,
		questions.NewShortAnswerCommand,
		questions.NewNumberCommand,
		questions.NewParagraphCommand,
		questions.NewOrderingCommand,
		questions.NewJavaScriptCommand,
		questions.NewJavaCommand,
		questions.NewPythonCommand,
		questions.NewSqlCommand,
		questions.NewRubyCommand,
		questions.NewUploadCommand,
		questions.NewCustomSnippetCommand,
		questions.NewProjectCommand,
		questions.NewTestableProjectCommand,
	}

	for _, factory := range questionCommands {
		cmd := factory(params)
		cmd.GroupID = "questions"
		c.AddCommand(cmd)
	}
}

func addYamlCommands(c *cobra.Command) {
	params := yaml.NewYamlCommandParams{
		RunCallback: runMarkdownSubcommand,
		Validator:   zeroOrOneFileNames,
	}

	fileCommands := []yaml.NewYamlCommand{
		yaml.NewConfigYamlCommand,
		yaml.NewDescriptionYamlCommand,
		yaml.NewCourseYamlCommand,
	}

	for _, factory := range fileCommands {
		cmd := factory(params)
		cmd.GroupID = "yaml"
		c.AddCommand(cmd)
	}
}

func addOtherCommands(c *cobra.Command) {
	params := others.NewOtherCommandParams{
		RunCallback: runMarkdownSubcommand,
		Validator:   zeroOrOneFileNames,
	}

	fileCommands := []others.NewOtherCommand{
		others.NewCallOutCommand,
		others.NewDistributeCodeCommand,
	}

	for _, factory := range fileCommands {
		cmd := factory(params)
		cmd.GroupID = "others"
		c.AddCommand(cmd)
	}
}

func addFileCommands(c *cobra.Command) {
	params := files.NewFileCommandParams{
		RunCallback: runMarkdownSubcommand,
		Validator:   zeroOrOneFileNames,
	}

	fileCommands := []files.NewFileCommand{
		files.NewLessonCommand,
		files.NewCheckpointCommand,
		files.NewSurveyCommand,
		files.NewInstructorCommand,
		files.NewResourceCommand,
		files.NewFileHeaderCommand,
	}

	for _, factory := range fileCommands {
		cmd := factory(params)
		cmd.GroupID = "files"
		c.AddCommand(cmd)
	}
}

func runMarkdownSubcommand(cmd *cobra.Command, destination *string, t templates.Template) {
	content := t.Render()
	toStdOut, _ := cmd.Flags().GetBool("out")
	writer := getWriter(destination, toStdOut)
	writer.Write(t.GetName(), content)
}

func zeroOrOneFileNames(cmd *cobra.Command, args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("%s takes zero or one arguments", cmd.Name())
	}
	if len(args) == 1 {
		if _, err := os.Stat(args[0]); os.IsNotExist(err) {
			return fmt.Errorf("%s does not exist to append to", args[0])
		}
	}
	return nil
}
