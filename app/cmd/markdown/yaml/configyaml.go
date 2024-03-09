package yaml

import (
	"github.com/spf13/cobra"
)

func NewConfigYamlCommand(params NewYamlCommandParams) *cobra.Command {
	params.name = "configyaml"
	params.fileName = "config.yaml"
	params.abbr = "cfy"
	params.maxTemplate = configYamlTemplate
	params.minTemplate = configYamlTemplateMin
	params.needsUID = false

	return createYamlCommand(params)
}

const configYamlTemplate = `# Config.yaml specifies the content and ordering within a curriculum block repo
#
# Supported Fields
# ==========================
# Standards -- (Standards = Units). An array of Units for a block
# Standard.Title -- The Unit title that shows up on the curriculum overview
# Standard.UID -- A unique ID for the Unit.
# Standard.Description -- The Unit description that shows up on the curriculum overview
# Standard.SuccessCriteria -- An array of success criteria that can be viewed when scoring the checkpoint in a Unit.
# Standard.ContentFiles -- An array of Lessons and (optional) Checkpoint in a Unit.
# Standard.ContentFiles.Type -- 'Lesson' or 'Checkpoint'
# Standard.ContentFiles.UID -- A unique ID for the lesson or checkpoint.
# Standard.ContentFiles.Path -- The absolute path to the Lesson, starting with /.
# Standard.ContentFiles.DefaultVisibility -- (optional) Set to 'hidden' to hide when a course first starts.
# Standard.ContentFiles.Autoscore -- (optional, for Checkpoints only) submit checkpoint scores without review
# Standard.ContentFiles.MaxCheckpointSubmissions -- (optional, for Checkpoints only) limit the number of submissions
# Standard.ContentFiles.EmailOnCompletion -- (optional, for Checkpoints only) Set to 'true' or 'false'. Sends instructors an email once student has completed a checkpoint
# Standard.ContentFiles.TimeLimit -- (optional, for Checkpoints only) the time limit in minutes
#
# Instructions
# ==========================
# Replace everything in square brackets [] and remove brackets
# Add all other Standards, Lessons, and Checkpoints following the pattern below
# All UIDs must be unique within a repo. You can use a uuidgen plugin.

---
Standards:
  - Title: [The unit name]
    UID: [unique-id]
    Description: [The Standard text]
    SuccessCriteria:
      - [The first success criteria]
    ContentFiles:
      - Type: Lesson
        UID: [unique-id]
        Path: /[folder/file.md]
      - Type: Checkpoint
        UID: [unique-id]
        Path: /[folder/file.md]`

const configYamlTemplateMin = `---
Standards:
  - Title:
    UID:
    Description:
    SuccessCriteria:
      -
    ContentFiles:
      - Type: Lesson
        UID:
        Path: /
      - Type: Checkpoint
        UID:
        Path: /`
