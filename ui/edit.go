package ui

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/bahner/go-ma-actor/config"
	log "github.com/sirupsen/logrus"
)

func (ui *ChatUI) invokeEditor() ([]byte, error) {
	tmpfile, err := os.CreateTemp("", "edit")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpfile.Name()) // clean up

	var editorErr error
	var contents []byte

	// Use Suspend to stop the TUI and run the external editor
	ui.app.Suspend(func() {
		// Launch external editor
		cmd := exec.Command(config.GetEditor(), tmpfile.Name())
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		editorErr = cmd.Run() // This will wait until the editor is closed

		if editorErr == nil {
			// Read the file only if there was no error with the editor
			contents, editorErr = os.ReadFile(tmpfile.Name())
		}
	})

	// Check if there was an error with the editor or reading the file
	if editorErr != nil {
		return nil, editorErr
	}

	// Remove newlies, We want this to be a single line
	// Trim right to remove trailing newlines
	contents = bytes.TrimRight(contents, "\n")
	// Replace all newlines with spaces (or another character if preferred)
	contents = bytes.ReplaceAll(contents, []byte("\n"), []byte(" "))
	// Append a single newline at the end if desired
	contents = append(contents, '\n')

	return contents, nil
}

func (ui *ChatUI) handleEditCommand(args []string) {

	m, err := ui.invokeEditor()
	if err != nil {
		ui.displaySystemMessage("Error invoking editor: " + err.Error())
		return
	}

	if len(m) == 0 {
		ui.displaySystemMessage("No changes made")
		return
	}

	log.Debugf("Editor returned: %s", m)

}
