package ui

import (
	"bytes"
	"errors"
	"os"
	"os/exec"

	"github.com/bahner/go-ma-actor/config"
	log "github.com/sirupsen/logrus"
)

const editorUsage = "'"
const editorHelp = "Opens the default editor to edit a message"

var ErrEmptyEdit = errors.New("empty edit")

func (ui *ChatUI) invokeEditor() ([]byte, error) {
	tmpfile, err := os.CreateTemp("", "edit")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpfile.Name()) // clean up

	var editorErr error
	var contents []byte
	ui.app.Suspend(func() {
		cmd := exec.Command(config.GetEditor(), tmpfile.Name())
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		editorErr = cmd.Run() // This will wait until the editor is closed
		if editorErr != nil {
			return
		}

		// Check if the file was modified by comparing modification times
		fi, err := os.Stat(tmpfile.Name())
		if err != nil {
			editorErr = err
			return
		}

		if fi.Size() == 0 {
			editorErr = ErrEmptyEdit
			return
		}

		contents, editorErr = os.ReadFile(tmpfile.Name())
	})

	if editorErr != nil {
		return nil, editorErr
	}

	if contents == nil {
		// It's unusual for contents to be nil without an error, but handle gracefully
		return nil, ErrEmptyEdit
	}

	// Remove newlies, We want this to be a single line
	// Trim right to remove trailing newlines
	contents = bytes.TrimRight(contents, "\n")

	// // Replace all newlines with spaces (or another character if preferred)
	// contents = bytes.ReplaceAll(contents, []byte("\n"), []byte(separator))

	// Append a single newline at the end if desired
	contents = append(contents, '\n')

	return contents, nil
}

func (ui *ChatUI) handleEditorCommand() {

	m, err := ui.invokeEditor()
	if err != nil {
		ui.displaySystemMessage("editor: " + err.Error())
		return
	}

	ui.chInput <- string(m)

	log.Debug("Editor command handled")

}
