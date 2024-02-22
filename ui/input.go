package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/spf13/viper"
)

func (ui *ChatUI) setupInputField() *tview.InputField {

	inputField := tview.NewInputField().
		SetLabel(viper.GetString("actor.nick") + ": ").
		SetFieldWidth(0).
		SetLabelColor(tcell.ColorBlack)

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
			ui.app.Stop()
			return
		}

		ui.inputHistory = append(ui.inputHistory, line)
		ui.currentHistoryIndex = -1

		// send the line onto the input chan and reset the field text
		ui.chInput <- line
		inputField.SetText("")
	})

	return inputField

}
