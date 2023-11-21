package ui

import (
	"context"
	"fmt"
	"io"

	"github.com/bahner/go-ma-actor/actor"
	"github.com/bahner/go-ma/did"
	"github.com/bahner/go-ma/did/doc"
	"github.com/bahner/go-ma/msg"
	"github.com/gdamore/tcell/v2"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/rivo/tview"
)

// ChatUI is a Text User Interface (TUI) for a Room.
// The Run method will draw the UI to the terminal in "fullscreen"
// mode. You can quit with Ctrl-C, or by typing "/quit" into the
// chat prompt.
type ChatUI struct {
	ctx context.Context
	n   host.Host

	nick string
	// We need this because for the acor, the chatui *is* the room.
	d *doc.Document

	// The actor is need to encrypt and sign messages in the event loop.
	a *actor.Actor

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
func NewChatUI(ctx context.Context, n host.Host, a *actor.Actor, id string) *ChatUI {

	var (
		u   ChatUI
		err error
	)

	// Assign actor to the room, so the room can handle messages on behalf of the actor.
	// This is because the ui has the event loop, and the actor doesn't.
	u.a = a

	u.d, err = doc.Fetch(id)
	if err != nil {
		panic(fmt.Sprintf("Failed to fetch DIDDOcument. %v", err))
	}

	u.nick = did.GetFragment(id)

	app := tview.NewApplication()

	// make a text view to contain our chat messages
	msgBox := tview.NewTextView()
	msgBox.SetDynamicColors(true)
	msgBox.SetBorder(true)
	msgBox.SetTitle(fmt.Sprintf("Entity: %s", u.nick))

	// text views are io.Writers, but they don't automatically refresh.
	// this sets a change handler to force the app to redraw when we get
	// new messages to display.
	msgBox.SetChangedFunc(func() {
		app.Draw()
	})

	// an input field for typing messages into
	chInput := make(chan string, 32)
	input := tview.NewInputField().
		SetLabel(u.nick + " > ").
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
	peersList.SetTitle("TODO")
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
		ctx:       ctx,
		n:         n,
		app:       app,
		peersList: peersList,
		msgW:      msgBox,
		msgBox:    msgBox,
		chInput:   chInput,
		chDone:    make(chan struct{}, 1),
	}
}

// Run starts the chat event loop in the background, then starts
// the event loop for the text UI.
func (ui *ChatUI) Run() error {
	go ui.handleEvents()
	defer ui.end()

	return ui.app.Run()
}

// end signals the event loop to exit gracefully
func (ui *ChatUI) end() {
	ui.chDone <- struct{}{}
}
