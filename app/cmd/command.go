// Package command provides wrappers around bash commands in order to support
// commands across multiple operating systems. The commands in this script
// should detect if a feature is supported, and if it is, use the feature that
// is best suited for the platform.
package cmd

import (
	"fmt"
	"os/exec"
)

// isCommandAvailable checks to see if a command is available.
// It returns true if the command is available, false otherwise
func isCommandAvailable(command string) bool {
	result := exec.Command("bash", "-c", fmt.Sprintf("command -v %s", command))
	if err := result.Run(); err != nil {
		return false
	}
	return true
}

// openURL attempts to open the specified URL in a new browser window.
// It checks to see if the open command is supported first, and if not, the
// start command. If either of the commands is supported, the appropriate
// command will be executed. If neither one is supported, a browser window will
// not be opened. No message is returned.
func openURL(url string) {
	var openBrowserCommand string

	if isCommandAvailable("open") {
		openBrowserCommand = "open"
	} else if isCommandAvailable("start") {
		openBrowserCommand = "start"
	}

	if len(openBrowserCommand) > 0 {
		exec.Command("bash", "-c", fmt.Sprintf("%s %s", openBrowserCommand, url)).Output()
	}
}
