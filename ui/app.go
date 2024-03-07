package ui

import (
	"github.com/bahner/go-ma-actor/config"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (ui *ChatUI) setupApp() {

	// Global style
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorNavajoWhite
	tview.Styles.PrimaryTextColor = tcell.ColorBlack
	tview.Styles.ContrastBackgroundColor = tcell.ColorNavajoWhite
	tview.Styles.BorderColor = tcell.ColorDarkGray
	tview.Styles.TitleColor = tcell.ColorDarkSlateGray

	// make a text view to contain our chat messages
	msgBox := setupMsgbox(ui.app)
	ui.msgBox = msgBox
	ui.msgW = msgBox

	// make a text view to hold the list of peers in the room, updated by ui.refreshPeers()
	peersList := tview.NewTextView()
	peersList.SetBorder(true)
	peersList.SetTitle("Peers")
	peersList.SetChangedFunc(func() { ui.app.Draw() })
	ui.peersList = peersList

	// chatPanel is a horizontal box with messages on the left and peers on the right
	// the peers list takes 20 columns, and the messages take the remaining space
	ui.chatPanel = tview.NewFlex().
		AddItem(msgBox, 0, 1, false).
		AddItem(peersList, config.UIPeerslistWidth(), 1, false)

	// The ordering here is a little kludgy, but acceptable for now.
	// the input fiield setup became rather verbose, so it was moved to its own file.
	ui.inputField = ui.setupInputField()

	ui.updateRoot()

}

func (ui *ChatUI) updateRoot() {
	ui.screen = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(ui.chatPanel, 0, 1, false).
		AddItem(ui.inputField, 1, 1, true) // True meaans it has focus.
	ui.app.SetRoot(ui.screen, true)
}

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
