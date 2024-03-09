package others

import (
	"github.com/spf13/cobra"
)

func NewCallOutCommand(params NewOtherCommandParams) *cobra.Command {
	params.name = "Call Out"
	params.abbr = "co"
	params.maxTemplate = callOutTemplate
	params.minTemplate = callOutTemplateMin

	return createOtherCommand(params)
}

const callOutTemplate = `<!-- available callout types: info, success, warning, danger, secondary, star  -->
### !callout-info

## title

body

### !end-callout`

const callOutTemplateMin = `### !callout-info

## title



### !end-callout`
