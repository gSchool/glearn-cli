package markdown

import (
	"os"
	"reflect"
	"testing"
)

func Test_getWriterReturnsClipboardWriterForNoFileAndNoStdOut(t *testing.T) {
	expected := reflect.TypeOf(clipboardWriter{})

	writer := getWriter(nil, false)

	writerType := reflect.TypeOf(writer)
	if writerType != expected {
		t.Errorf("Did not get expected clipboard writer: %s != %s", writerType, expected)
	}
}

func Test_getWriterReturnsFileAppendWriterForFileName(t *testing.T) {
	expected := reflect.TypeOf(fileAppendWriter{})
	name := "foo"

	writer := getWriter(&name, false)

	writerType := reflect.TypeOf(writer)
	if writerType != expected {
		t.Errorf("Did not get expected file append writer: %s != %s", writerType, expected)
	}

	faw := writer.(fileAppendWriter)

	if faw.destination != "foo" {
		t.Error("File append writer did not have the correct destination")
	}
}

func Test_getWriterReturnsStandardOutWriterForNoFileNameAndStdOut(t *testing.T) {
	expected := reflect.TypeOf(standardOutWriter{})

	writer := getWriter(nil, true)

	writerType := reflect.TypeOf(writer)
	if writerType != expected {
		t.Errorf("Did not get expected standard out writer: %s != %s", writerType, expected)
	}
}

func Test_StandardOutWritingDoesNotThrowError(t *testing.T) {
	if err := getWriter(nil, true).Write("Name", "Content"); err != nil {
		t.Errorf("Got an error writing to stdout: %s", err)
	}
}

func Test_ClipboardWritingDoesNotThrowError(t *testing.T) {
	if err := getWriter(nil, false).Write("Name", "Content"); err != nil {
		t.Errorf("Got an error writing to stdout: %s", err)
	}
}

func Test_FileAppendWritingFailsForNonMdYamlFile(t *testing.T) {
	destination := "not-md-yaml"
	if err := getWriter(&destination, false).Write("Name", "Content"); err == nil {
		t.Error("Failed to get error for bad file name")
	}
}

func Test_FileAppendWritingFailsForFileNotExisting(t *testing.T) {
	destination := "does-not-exist.yaml"
	if err := getWriter(&destination, false).Write("Name", "Content"); err == nil {
		t.Error("Failed to get error for bad file name")
	}
}

func Test_FileAppendWritingFailsForDirectoryTarget(t *testing.T) {
	destination := "./foo.yaml"
	os.Mkdir(destination, 0750)
	defer os.RemoveAll(destination)
	if err := getWriter(&destination, false).Write("Name", "Content"); err == nil {
		t.Error("Failed to get error for bad file name")
	}
}
