package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/bahner/go-ma"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// // displaySelfMessage writes a message from ourself to the message window,
// // with our nick highlighted in yellow.
// func (ui *ChatUI) displaySelfMessage(msg string) {
// 	prompt := withColor("yellow", fmt.Sprintf("<%s>:", ui.e.Nick))
// 	fmt.Fprintf(ui.msgW, "%s %s\n", prompt, msg)
// }

// withColor wraps a string with color tags for display in the messages text box.
func withColor(color, msg string) string {
	return fmt.Sprintf("[%s]%s[-]", color, msg)
}

func (ui *ChatUI) displaySystemMessage(msg string) {
	prompt := withColor("cyan", "[SYSTEM]:")
	fmt.Fprintf(ui.msgW, "%s %s\n", prompt, msg)
}

func (ui *ChatUI) setupInputField(inputField *tview.InputField, app *tview.Application) {
	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyUp:
			if ui.currentHistoryIndex < len(ui.inputHistory)-1 {
				ui.currentHistoryIndex++
				input := ui.inputHistory[len(ui.inputHistory)-1-ui.currentHistoryIndex]
				inputField.SetText(input)
				return nil // event handled
			}
		case tcell.KeyDown:
			if ui.currentHistoryIndex > 0 {
				ui.currentHistoryIndex--
				input := ui.inputHistory[len(ui.inputHistory)-1-ui.currentHistoryIndex]
				inputField.SetText(input)
				return nil // event handled
			} else if ui.currentHistoryIndex == 0 {
				ui.currentHistoryIndex = -1
				inputField.SetText("") // Clear the input field
				return nil             // event handled
			}
		}
		return event // let other keys pass through
	})

	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			input := inputField.GetText()
			if input != "" {
				ui.inputHistory = append(ui.inputHistory, input) // Add to history
				ui.currentHistoryIndex = -1                      // Reset index
				ui.chInput <- input                              // Send input to be handled
				inputField.SetText("")                           // Clear the input field
			}
		}
	})

	// the done func is called when the user hits enter, or tabs out of the field
	inputField.SetDoneFunc(func(key tcell.Key) {
		if key != tcell.KeyEnter {
			// we don't want to do anything if they just tabbed away
			return
		}

		line := inputField.GetText()
		if len(line) == 0 {
			// ignore blank lines
			return
		}

		// bail if requested
		if line == "/quit" {
			app.Stop()
			return
		}

		ui.inputHistory = append(ui.inputHistory, line)
		ui.currentHistoryIndex = -1

		// send the line onto the input chan and reset the field text
		ui.chInput <- line
		inputField.SetText("")
	})

}

// handleEvents runs an event loop that sends user input to the chat room
// and displays messages received from the chat room. It also periodically
// refreshes the list of peers in the UI.

func (ui *ChatUI) handleEvents() {
	peerRefreshTicker := time.NewTicker(time.Second)
	defer peerRefreshTicker.Stop()

	for {
		select {
		case input := <-ui.chInput:
			if strings.HasPrefix(input, "/") {
				ui.handleCommands(input)
			} else {
				ui.handleChatMessage(input)
			}

		case m := <-ui.chMessage:
			if m.MimeType == ma.BROADCAST_MIME_TYPE {
				ui.displayBroadcastMessage(m)
				continue
			}
			ui.displayChatMessage(m)

		case <-peerRefreshTicker.C:
			ui.refreshPeers()

		case <-ui.chDone:
			return
		}
	}
}
