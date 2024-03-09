//go:build integration

package markdown

import (
	"testing"
)

func Test_ClipboardWritingDoesNotThrowError(t *testing.T) {
	if err := getWriter(nil, false).Write("Name", "Content"); err != nil {
		t.Errorf("Got an error writing to stdout: %s", err)
	}
}
