package yaml

import (
	"fmt"
	"strings"

	"github.com/gSchool/glearn-cli/app/cmd/markdown/templates"
	"github.com/spf13/cobra"
)

type NewYamlCommand func(NewYamlCommandParams) *cobra.Command

type NewYamlCommandParams struct {
	RunCallback func(*cobra.Command, *string, templates.Template)
	Validator   cobra.PositionalArgs
	name        string
	fileName    string
	abbr        string
	maxTemplate string
	minTemplate string
}

func createYamlCommand(params NewYamlCommandParams) *cobra.Command {
	commandName := strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(params.name), " ", ""), ".", "")
	use := fmt.Sprintf("%s [file-to-append-to]", commandName)
	short := fmt.Sprintf("(%s) Generate the content for a %s", params.abbr, params.fileName)
	long := fmt.Sprintf("Create the content for a %s to the clipboard, standard out, or appended to a file.", params.fileName)

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
			t := templates.NewStaticTemplate(templateName, template)
			var fileName *string
			if len(args) == 1 {
				fileName = &args[0]
			}
			params.RunCallback(cmd, fileName, t)
		},
	}
}
