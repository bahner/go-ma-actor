package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/spf13/viper"
)

func historySize() int {
	return viper.GetInt("ui.history-size")
}

func (ui *ChatUI) pushToHistory(line string) {

	historySize := historySize()

	if len(ui.inputHistory) == historySize {
		// Remove the oldest entry when we reach max size
		copy(ui.inputHistory, ui.inputHistory[1:])
		ui.inputHistory = ui.inputHistory[:historySize-1]
	}
	ui.inputHistory = append(ui.inputHistory, line)
}

func (ui *ChatUI) setupInputField() *tview.InputField {
	inputField := tview.NewInputField().
		SetLabel(viper.GetString("actor.nick") + ": ").
		SetFieldWidth(0).
		SetLabelColor(tcell.ColorBlack).
		SetText("/help")

	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEscape:
			ui.app.QueueUpdateDraw(func() {
				ui.app.SetFocus(ui.inputField)
			})
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
		if key != tcell.KeyEnter {
			// we don't want to do anything if they just tabbed away
			return
		}

		line := inputField.GetText()
		if len(line) == 0 {
			return
		}

		if line == "/quit" {
			ui.app.Stop()
			return
		}

		ui.pushToHistory(line)
		ui.currentHistoryIndex = -1

		ui.chInput <- line
		inputField.SetText("")
	})

	return inputField
}
