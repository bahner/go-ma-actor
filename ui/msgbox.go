package ui

import "github.com/rivo/tview"

func setupMsgbox(app *tview.Application) *tview.TextView {

	// make a text view to contain our chat messages
	msgBox := tview.NewTextView()
	msgBox.SetDynamicColors(true)
	msgBox.SetBorder(true)
	msgBox.SetScrollable(true)
	msgBox.SetTitle(defaultLimbo)

	// text views are io.Writers, but they don't automatically refresh.
	// this sets a change handler to force the app to redraw when we get
	// new messages to display.
	msgBox.SetChangedFunc(func() {
		app.Draw()
		msgBox.ScrollToEnd()
	})

	return msgBox
}
