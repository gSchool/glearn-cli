package markdown

import "testing"

func Test_MarkdownCommandIsConfigured(t *testing.T) {
	cmd := NewMarkdownCommand()

	if len(cmd.Groups()) == 0 {
		t.Error("Markdown command did not have any groups")
	}

	if len(cmd.Commands()) == 0 {
		t.Error("Markdown command did not have any commands")
	}
}

func Test_ZeroOrOneFileNamesValidatesArgs(t *testing.T) {
	if err := zeroOrOneFileNames(nil, []string{}); err != nil {
		t.Error("Validator fails with no arguments")
	}

	if err := zeroOrOneFileNames(nil, []string{"does-not-exist.png"}); err == nil {
		t.Error("Validator passed for non-existent file does-not-exist.png")
	}

	if err := zeroOrOneFileNames(nil, []string{"root_test.go"}); err != nil {
		t.Error("Validator failed for existing file root_test.go")
	}
}
