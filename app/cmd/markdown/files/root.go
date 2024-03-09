package files

import (
	"fmt"
	"strings"

	"github.com/gSchool/glearn-cli/app/cmd/markdown/templates"
	"github.com/spf13/cobra"
)

type NewFileCommand func(NewFileCommandParams) *cobra.Command

type NewFileCommandParams struct {
	RunCallback func(*cobra.Command, *string, templates.Template)
	Validator   cobra.PositionalArgs
	name        string
	abbr        string
	maxTemplate string
	minTemplate string
}

func createFileCommand(params NewFileCommandParams) *cobra.Command {
	commandName := strings.ReplaceAll(strings.ToLower(params.name), " ", "")
	use := fmt.Sprintf("%s [file-to-append-to]", commandName)
	short := fmt.Sprintf("(%s) Generate the content of a %s file", params.abbr, params.name)
	long := fmt.Sprintf("Create a %s file to the clipboard, standard out, or appended to a file.", params.name)

	return &cobra.Command{
		Use:     use,
		Aliases: []string{params.abbr},
		Short:   short,
		Long:    long,
		Args:    params.Validator,
		Run: func(cmd *cobra.Command, args []string) {
			template := params.maxTemplate
			if isMin, err := cmd.Flags().GetBool("min"); err == nil && isMin {
				template = params.minTemplate
			}
			templateName := fmt.Sprintf("%s Template", params.name)
			t := templates.NewIdTemplate(templateName, template)
			var fileName *string
			if len(args) == 1 {
				fileName = &args[0]
			}
			params.RunCallback(cmd, fileName, t)
		},
	}
}
