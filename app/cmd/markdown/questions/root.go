package questions

import (
	"fmt"
	"strings"

	"github.com/gSchool/glearn-cli/app/cmd/markdown/templates"
	"github.com/spf13/cobra"
)

type NewQuestionCommand func(NewQuestionCommandParams) *cobra.Command

type NewQuestionCommandParams struct {
	RunCallback func(*cobra.Command, *string, templates.Template)
	Validator   cobra.PositionalArgs
	name        string
	abbr        string
	maxTemplate string
	minTemplate string
	long        string
}

func createQuestionCommand(params NewQuestionCommandParams) *cobra.Command {
	commandName := strings.ReplaceAll(strings.ToLower(params.name), " ", "")
	use := fmt.Sprintf("%s [file-to-append-to]", commandName)
	short := fmt.Sprintf("(%s) Generate the content for a %s block", params.abbr, params.name)
	long := fmt.Sprintf("Create the content for a %s block to the clipboard, standard out, or appended to a file.", params.name)
	if len(params.long) > 0 {
		long = params.long
	}

	cmd := &cobra.Command{
		Use:     use,
		Aliases: []string{params.abbr},
		Short:   short,
		Long:    long,
		Args:    params.Validator,
		Run: func(cmd *cobra.Command, args []string) {
			template := params.maxTemplate
			isMin := true
			var err error = nil
			if isMin, err = cmd.Flags().GetBool("min"); err == nil && isMin {
				template = params.minTemplate
			}
			withExplanation := false
			withRubric := false
			withHints := 0
			if we, err := cmd.Flags().GetBool("with-explanation"); err == nil {
				withExplanation = we
			}
			if wr, err := cmd.Flags().GetBool("with-rubric"); err == nil {
				withRubric = wr
			}
			if wh, err := cmd.Flags().GetInt("with-hints"); err == nil {
				withHints = wh
			}
			templateName := fmt.Sprintf("%s Template", params.name)
			t := templates.NewAttributeTemplate(
				templateName,
				template,
				isMin,
				withExplanation,
				withRubric,
				withHints,
			)
			var fileName *string
			if len(args) == 1 {
				fileName = &args[0]
			}
			params.RunCallback(cmd, fileName, t)
		},
	}

	cmd.Flags().BoolP("with-explanation", "e", false, "Include explanation blocks")
	cmd.Flags().BoolP("with-rubric", "r", false, "Include rubric blocks")
	cmd.Flags().IntP("with-hints", "n", 0, "Include n number of hint blocks")

	return cmd
}
