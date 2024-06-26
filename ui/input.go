package ui

import (
	"github.com/bahner/go-ma-actor/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (ui *ChatUI) setupInputField() *tview.InputField {
	inputField := tview.NewInputField().
		SetLabel(config.ActorNick() + ": ").
		SetFieldWidth(0).
		SetLabelColor(tcell.ColorBlack).
		SetText("/help")

	inputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyCtrlC:
			ui.displaySystemMessage("Ctrl-C is disabled. Use /quit to exit.")
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
