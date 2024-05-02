package files

import (
	"github.com/spf13/cobra"
)

func NewFileHeaderCommand(params NewFileCommandParams) *cobra.Command {
	params.name = "File Header"
	params.abbr = "fh"
	params.maxTemplate = fileHeaderTemplate
	params.minTemplate = fileHeaderTemplateMin

	return createFileCommand(params)
}

const fileHeaderTemplate = `---
# BEGIN FILE CONFIGURATION YML HEADER >>>>>
# autoconfig.yml will use these settings. config.yml will override.
Type: Lesson # Options: Lesson, Checkpoint, Survey, Instructor, Resource
UID: %s
# DefaultVisibility: hidden # Uncomment this line to default Lesson to hidden. Please note, the default visibility setting is applied only during the initial sync of a course file within a cohort.
# MaxCheckpointSubmissions: 1 # Checkpoints only. Uncomment this line to limit the number of submissions
# EmailOnCompletion: true # Checkpoints only. Uncomment this line to send instructors an email once a student has completed a checkpoint
# TimeLimit: 60 # Checkpoints only. Uncomment this line to set a time limit in minutes
# Autoscore: true # Checkpoints only. Uncomment this line to finalize checkpoint scores without instructor review
# END FILE CONFIGURATION YML HEADER <<<<<
---`

const fileHeaderTemplateMin = `---
Type: Checkpoint
UID: %s
# DefaultVisibility: hidden
# MaxCheckpointSubmissions: 1
# EmailOnCompletion: true
# TimeLimit: 60
# Autoscore: true
---`
