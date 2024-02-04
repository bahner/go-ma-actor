package ui

import (
	"fmt"
	"io"

	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma-actor/p2p"
	"github.com/bahner/go-ma/msg"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ChatUI is a Text User Interface (TUI) for a Room.
// The Run method will draw the UI to the terminal in "fullscreen"
// mode. You can quit with Ctrl-C, or by typing "/quit" into the
// chat prompt.
type ChatUI struct {
	p *p2p.P2P

	e *entity.Entity

	// The actor is need to encrypt and sign messages in the event loop.
	a *entity.Entity

	// The Topic is used for publication of messages after encryption and signing.
	// The names are obviously, from the corresponding DIDDocument.

	app       *tview.Application
	peersList *tview.TextView
	msgBox    *tview.TextView

	msgW      io.Writer
	chInput   chan string
	chMessage chan *msg.Message
	chDone    chan struct{}
}

// NewChatUI returns a new ChatUI struct that controls the text UI.
// It won't actually do anything until you call Run().
// The enity is the "room" we are convering with.
func NewChatUI(p *p2p.P2P, a *entity.Entity, e *entity.Entity) (*ChatUI, error) {

	err := e.Verify()
	if err != nil {
		return nil, fmt.Errorf("entity verification error: %w", err)
	}

	app := tview.NewApplication()

	// make a text view to contain our chat messages
	msgBox := tview.NewTextView()
	msgBox.SetDynamicColors(true)
	msgBox.SetBorder(true)
	msgBox.SetScrollable(true)
	msgBox.SetTitle(fmt.Sprintf("Entity: %s", e.Nick))

	// text views are io.Writers, but they don't automatically refresh.
	// this sets a change handler to force the app to redraw when we get
	// new messages to display.
	msgBox.SetChangedFunc(func() {
		app.Draw()
		msgBox.ScrollToEnd()
	})

	// an input field for typing messages into
	chInput := make(chan string, 32)
	input := tview.NewInputField().
		SetLabel(e.Nick + " > ").
		SetFieldWidth(0).
		SetFieldBackgroundColor(tcell.ColorBlack)

	// the done func is called when the user hits enter, or tabs out of the field
	input.SetDoneFunc(func(key tcell.Key) {
		if key != tcell.KeyEnter {
			// we don't want to do anything if they just tabbed away
			return
		}
		line := input.GetText()
		if len(line) == 0 {
			// ignore blank lines
			return
		}

		// bail if requested
		if line == "/quit" {
			app.Stop()
			return
		}

		// send the line onto the input chan and reset the field text
		chInput <- line
		input.SetText("")
	})

	// make a text view to hold the list of peers in the room, updated by ui.refreshPeers()
	peersList := tview.NewTextView()
	peersList.SetBorder(true)
	peersList.SetTitle("Peers")
	peersList.SetChangedFunc(func() { app.Draw() })

	// chatPanel is a horizontal box with messages on the left and peers on the right
	// the peers list takes 20 columns, and the messages take the remaining space
	chatPanel := tview.NewFlex().
		AddItem(msgBox, 0, 1, false).
		AddItem(peersList, 20, 1, false)

	// flex is a vertical box with the chatPanel on top and the input field at the bottom.

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(chatPanel, 0, 1, false).
		AddItem(input, 1, 1, true)

	app.SetRoot(flex, true)

	return &ChatUI{
		a:         a,
		p:         p,
		e:         e,
		app:       app,
		peersList: peersList,
		msgW:      msgBox,
		msgBox:    msgBox,
		chInput:   chInput,
		chDone:    make(chan struct{}, 1),
	}, nil
}

// Run starts the chat event loop in the background, then starts
// the event loop for the text UI.
func (ui *ChatUI) Run() error {

	defer ui.end()

	go ui.setEntity(ui.e.DID.String())
	go ui.setActor(ui.a)
	go ui.handleEvents()

	return ui.app.Run()
}

// end signals the event loop to exit gracefully
func (ui *ChatUI) end() {
	ui.chDone <- struct{}{}
}
