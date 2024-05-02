package files

import (
	"github.com/spf13/cobra"
)

func NewCheckpointCommand(params NewFileCommandParams) *cobra.Command {
	params.name = "Checkpoint"
	params.abbr = "cp"
	params.maxTemplate = checkpointTemplate
	params.minTemplate = checkpointTemplateMin

	return createFileCommand(params)
}

const checkpointTemplate = `---
# BEGIN FILE CONFIGURATION YML HEADER >>>>>
# autoconfig.yml will use these settings. config.yml will override.
Type: Checkpoint
UID: %s
# DefaultVisibility: hidden # Uncomment this line to default Checkpoint to hidden. Please Note, this setting is applied only during the initial sync of a course file with a cohort.
# MaxCheckpointSubmissions: 1 # Uncomment this line to limit the number of submissions
# EmailOnCompletion: true #  Uncomment this line to send instructors an email once a student has completed a checkpoint
# TimeLimit: 60 # Uncomment this line to set a time limit in minutes
# Autoscore: true # Uncomment this line to finalize checkpoint scores without instructor review
# END FILE CONFIGURATION YML HEADER <<<<<
---

# Title

<!--A Checkpoint is an assessment and must include include one or more Challenges. -->
<!--A Checkpoint can have any markdown. See examples of markdown formatting by running 'learn walkthrough' and previewing the tutorial.-->
`

const checkpointTemplateMin = `---
Type: Checkpoint
UID: %s
# DefaultVisibility: hidden
# MaxCheckpointSubmissions: 1
# EmailOnCompletion: true
# TimeLimit: 60
# Autoscore: true
---

# Title
`
