package yaml

import (
	"github.com/spf13/cobra"
)

func NewDescriptionYamlCommand(params NewYamlCommandParams) *cobra.Command {
	params.name = "descyaml"
	params.fileName = "description.yaml"
	params.abbr = "dsy"
	params.maxTemplate = descYamlTemplate
	params.minTemplate = descYamlTemplateMin
	params.needsUID = true

	return createYamlCommand(params)
}

const descYamlTemplate = `# description.yaml file template, defines Unit/Standard details when generating an autoconfig.yaml
# Populates a file called 'description.yaml'. When placed in a unit directory, autoconfig.yaml generation will read this file
# and apply these settings to the block's autoconfig.
---
Title: [Unit Title] (appears when viewing content files or curriculum overview)
Description: [Description] (longer text shown with each Title on curriculum overview)
UID: %s
SuccessCriteria:
 - [Success Criteria] (define what success for this Unit means)
`
const descYamlTemplateMin = `---
Title:
Description:
UID: %s
SuccessCriteria:
 -`
