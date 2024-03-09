package yaml

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/gSchool/glearn-cli/app/cmd/markdown/templates"
	"github.com/spf13/cobra"
	"go.uber.org/mock/gomock"
)

type testCase struct {
	name        string
	abbr        string
	command     NewYamlCommand
	maxTemplate string
	minTemplate string
	hasUID      bool
}

func getTestCases() []testCase {
	return []testCase{
		{"configyaml", "cfy", NewConfigYamlCommand, configYamlTemplate, configYamlTemplateMin, false},
		{"courseyaml", "cry", NewCourseYamlCommand, courseYamlTemplate, courseYamlTemplateMin, false},
		{"descyaml", "dsy", NewDescriptionYamlCommand, descYamlTemplate, descYamlTemplateMin, true},
	}
}

func Test_NoDestinationWorksAsExpected(t *testing.T) {
	testCases := getTestCases()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			callback := NewMockRunCallback(ctrl)
			validator := NewMockValidator(ctrl)

			var template templates.Template

			callback.
				EXPECT().
				Call(gomock.Any(), gomock.Any(), gomock.Any()).
				Do(func(_ *cobra.Command, _ *string, t templates.Template) {
					template = t
				})
			validator.EXPECT().Call(gomock.Any(), gomock.Any()).Return(nil)

			params := NewYamlCommandParams{
				RunCallback: callback.Call,
				Validator:   validator.Call,
			}
			cmd := tc.command(params)

			if cmd.Use != tc.name+" [file-to-append-to]" {
				t.Error("Command use is not what is expected")
			}
			if len(cmd.Aliases) != 1 || cmd.Aliases[0] != tc.abbr {
				t.Error("Command alias is not what is expected")
			}
			err := cmd.Execute()
			if err != nil {
				t.Error("Command erred when it should not have")
			}
			if template.GetUnrenderedContent() != tc.maxTemplate {
				t.Errorf("Failed to get the expected template")
			}
			re := regexp.MustCompile(`UID: [[:alnum:]]{8}-[[:alnum:]]{4}-[[:alnum:]]{4}-[[:alnum:]]{4}-[[:alnum:]]{12}`)
			if tc.hasUID {
				rendered := template.Render()
				if !re.MatchString(rendered) {
					t.Error("Did not find a UID GUID in the rendered file:\n" + rendered)
				}
			}
		})
	}
}

func Test_MinimumTemplateProvidedByFlag(t *testing.T) {
	testCases := getTestCases()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			callback := NewMockRunCallback(ctrl)
			validator := NewMockValidator(ctrl)

			var template templates.Template

			callback.
				EXPECT().
				Call(gomock.Any(), gomock.Any(), gomock.Any()).
				Do(func(_ *cobra.Command, _ *string, t templates.Template) {
					template = t
				})
			validator.EXPECT().Call(gomock.Any(), gomock.Any()).Return(nil)

			params := NewYamlCommandParams{
				RunCallback: callback.Call,
				Validator:   validator.Call,
			}
			cmd := tc.command(params)

			if cmd.Use != tc.name+" [file-to-append-to]" {
				t.Error("Command use is not what is expected")
			}
			if len(cmd.Aliases) != 1 || cmd.Aliases[0] != tc.abbr {
				t.Error("Command alias is not what is expected")
			}
			cmd.Flags().BoolP("min", "m", false, "")
			cmd.SetArgs([]string{"--min"})
			err := cmd.Execute()
			if err != nil {
				t.Error("Command erred when it should not have")
			}
			if template.GetUnrenderedContent() != tc.minTemplate {
				t.Errorf("Failed to get the expected template")
			}
		})
	}
}

func Test__DestinationIsPassedToValidator(t *testing.T) {
	testCases := getTestCases()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			callback := NewMockRunCallback(ctrl)
			validator := NewMockValidator(ctrl)

			callback.EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any())
			validator.EXPECT().Call(gomock.Any(), []string{"output.md"}).Return(nil)

			params := NewYamlCommandParams{
				RunCallback: callback.Call,
				Validator:   validator.Call,
			}
			cmd := tc.command(params)

			cmd.SetArgs([]string{"output.md"})
			err := cmd.Execute()
			if err != nil {
				t.Error("Command erred when it should not have")
			}
		})
	}
}

func Test_WhenValidatorErrsCommandDoesNotCallCallback(t *testing.T) {
	testCases := getTestCases()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			callback := NewMockRunCallback(ctrl)
			validator := NewMockValidator(ctrl)
			srcErr := fmt.Errorf("Just your average error")

			callback.EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			validator.EXPECT().Call(gomock.Any(), []string{"output.md"}).Return(srcErr)

			params := NewYamlCommandParams{
				RunCallback: callback.Call,
				Validator:   validator.Call,
			}
			cmd := tc.command(params)

			cmd.SetArgs([]string{"output.md"})
			err := cmd.Execute()
			if err == nil {
				t.Error("Command did not err on argument validator")
			}
			if err != srcErr {
				t.Error("Returned error did not match the source error")
			}
		})
	}
}
