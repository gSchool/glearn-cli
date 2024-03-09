package files

import (
	"github.com/spf13/cobra"
)

func NewResourceCommand(params NewFileCommandParams) *cobra.Command {
	params.name = "Resource"
	params.abbr = "rs"
	params.maxTemplate = resourceTemplate
	params.minTemplate = resourceTemplateMin

	return createFileCommand(params)
}

const resourceTemplate = `---
# BEGIN FILE CONFIGURATION YML HEADER >>>>>
# autoconfig.yml will use these settings. config.yml will override.
Type: Resource
UID: %s
# END FILE CONFIGURATION YML HEADER <<<<<
---

# Title

<!--A Resource can have all of the same markdown and challenges as a lesson. Resources do not appear in the left nav and don't count toward course completion. -->
`

const resourceTemplateMin = `---
Type: Resource
UID: %s
---

# Title
`
