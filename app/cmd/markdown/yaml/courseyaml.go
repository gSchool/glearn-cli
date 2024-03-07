package yaml

import (
	"github.com/spf13/cobra"
)

func NewCourseYamlCommand(params NewYamlCommandParams) *cobra.Command {
	params.name = "courseyaml"
	params.fileName = "course.yaml"
	params.abbr = "cry"
	params.maxTemplate = courseYamlTemplate
	params.minTemplate = courseYamlTemplateMin

	return createYamlCommand(params)
}

const courseYamlTemplate = `# Course.yaml files specify the grouping and ordering of repos that define a course.
#
# Supported Fields
# ===================
# DefaultUnitVisibility -- (optional) set to 'hidden' to hide all units when a course first starts.
# Course -- The top level array containing the sections of a course
# Course.Section -- An array contining a single array of repos. Content in the same section is grouped together on curriculum homepage.
# Course.Repos --  An array containing block repos that have been published in Learn.
# Course.Repos.URL -- The URL to a block repo that has been published in Learn.
# DefaultUpdates: -- Optional. Sets branches to automatically or manually receive updates. Leaving blank defaults to auto. This setting can be overridden or set through the UI.
# Instructions
# ==========================
# Replace everything in square brackets [] and remove brackets
# Add all other Sections and Repos following the pattern below
# All UIDs must be unique within a repo. You can use a uuidgen plugin.

---
# DefaultUnitVisibility: hidden
Course:
  - Section: [Section name]
    Repos:
      -  URL: https://example-repo-1
         DefaultUpdates: [manual | auto]
      -  URL: https://example-repo-2
         DefaultUpdates: [manual | auto]`

const courseYamlTemplateMin = `---
# DefaultUnitVisibility: hidden
Course:
  - Section:
    Repos:
      - URL:
        DefaultUpdates: auto`
