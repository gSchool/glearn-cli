package markdown

import (
	"fmt"
	"os"
	"strings"

	"github.com/atotto/clipboard"
)

type TemplateWriter interface {
	Write(name string, p string) error
}

func getWriter(fileName *string, stdOut bool) TemplateWriter {
	var writer TemplateWriter = clipboardWriter{}
	if fileName != nil && stdOut {
		fmt.Println("INFO: stdout ignored because file specified")
	}
	if fileName != nil {
		writer = fileAppendWriter{*fileName}
	} else if stdOut {
		writer = standardOutWriter{}
	}
	return writer
}

type clipboardWriter struct{}

func (clipboardWriter) Write(name string, p string) error {
	if err := clipboard.WriteAll(p); err != nil {
		return err
	}
	fmt.Println(name, "copied to clipboard!")
	return nil
}

type standardOutWriter struct{}

func (standardOutWriter) Write(name string, p string) error {
	fmt.Println(p)
	return nil
}

type fileAppendWriter struct {
	destination string
}

func (fa fileAppendWriter) Write(name string, template string) error {
	target := fa.destination

	if !(strings.HasSuffix(target, ".md") || strings.HasSuffix(target, ".yaml") || strings.HasSuffix(target, ".yml")) {
		return fmt.Errorf("'%s' must have an `.md`, `.yml`, or `.yaml` extension to append %s content.\n", target, name)
	}

	targetInfo, err := os.Stat(target)
	if err != nil {
		return fmt.Errorf("'%s' is not a file that can be appended!\n%s\n", target, err)
	}
	if targetInfo.IsDir() {
		return fmt.Errorf("'%s' is a directory, please specify a markdown file.\n", target)
	}

	f, err := os.OpenFile(target, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("Cannot open '%s'!\n%s\n", target, err)
	}
	defer f.Close()

	if _, err = f.WriteString(template + "\n"); err != nil {
		return fmt.Errorf("Cannot write to '%s'!\n%s\n", target, err)
	}
	fmt.Printf("%s appended to %s!\n", name, target)

	return nil
}
