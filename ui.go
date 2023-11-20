package ui

import (
	"context"
	"fmt"
	"io"

	"github.com/bahner/go-ma/did/doc"
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

	name string
	d    *doc.Document

	app       *tview.Application
	peersList *tview.TextView
	msgBox    *tview.TextView

	msgW    io.Writer
	inputCh chan string
	doneCh  chan struct{}
}

// NewChatUI returns a new ChatUI struct that controls the text UI.
// It won't actually do anything until you call Run().
func NewChatUI(ctx context.Context, n host.Host, did string) *ChatUI {

	var (
		u   ChatUI
		err error
	)

	u.d, err = doc.Fetch(did)
	if err != nil {
		panic(fmt.Sprintf("Failed to fetch DIDDOcument. %v", err))
	}

	app := tview.NewApplication()
	roomNick := n.ID().String()

	// make a text view to contain our chat messages
	msgBox := tview.NewTextView()
	msgBox.SetDynamicColors(true)
	msgBox.SetBorder(true)
	msgBox.SetTitle(fmt.Sprintf("Entity: %s", roomNick))

	// text views are io.Writers, but they don't automatically refresh.
	// this sets a change handler to force the app to redraw when we get
	// new messages to display.
	msgBox.SetChangedFunc(func() {
		app.Draw()
	})

	// an input field for typing messages into
	inputCh := make(chan string, 32)
	input := tview.NewInputField().
		SetLabel(roomNick + " > ").
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
		inputCh <- line
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
		ps:        ps,
		r:         r,
		a:         a,
		app:       app,
		peersList: peersList,
		msgW:      msgBox,
		msgBox:    msgBox,
		inputCh:   inputCh,
		doneCh:    make(chan struct{}, 1),
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
	ui.doneCh <- struct{}{}
}

// refreshPeers pulls the list of peers currently in the chat room and
// displays the last 8 chars of their peer id in the Peers panel in the ui.
func (ui *ChatUI) refreshPeers() {
	peers := ui.r.Actor.Public.ListPeers()

	// clear is thread-safe
	ui.peersList.Clear()

	for _, p := range peers {
		fmt.Fprintln(ui.peersList, p.ShortString())
	}

	ui.app.Draw()
}

// displayChatMessage writes a ChatMessage from the room to the message window,
// with the sender's nick highlighted in green.
func (ui *ChatUI) displayChatMessage(cm *msg.Message) {
	prompt := withColor("green", fmt.Sprintf("<%s>:", cm.From))
	fmt.Fprintf(ui.msgW, "%s %s\n", prompt, cm.Body)
}

// displaySelfMessage writes a message from ourself to the message window,
// with our nick highlighted in yellow.
func (ui *ChatUI) displaySelfMessage(msg string) {
	prompt := withColor("yellow", fmt.Sprintf("<%s>:", ui.r.Entity.DID.Fragment))
	fmt.Fprintf(ui.msgW, "%s %s\n", prompt, msg)
}

// handleEvents runs an event loop that sends user input to the chat room
// and displays messages received from the chat room. It also periodically
// refreshes the list of peers in the UI.
func (ui *ChatUI) handleEvents() {
	peerRefreshTicker := time.NewTicker(time.Second)
	defer peerRefreshTicker.Stop()

	for {
		select {
		case input := <-ui.inputCh:
			if strings.HasPrefix(input, "/") {
				ui.handleCommands(input)
			} else {
				ui.handleChatMessage(input)
			}

		case m := <-ui.a.Messages:
			ui.displayChatMessage(m)

		case <-peerRefreshTicker.C:
			ui.refreshPeers()

		case <-ui.ctx.Done():
			return

		case <-ui.doneCh:
			return
		}
	}
}

func (ui *ChatUI) handleCommands(input string) {
	args := strings.Split(input, " ")

	switch args[0] {
	case "/status":
		ui.handleStatusCommand(args)
	case "/discover":
		ui.triggerDiscovery()
	case "/enter":
		ui.handleEnterCommand(args)
	case "/refresh":
		ui.app.Draw()
	default:
		ui.displaySystemMessage("Unknown command: " + args[0])
	}
}

func (ui *ChatUI) handleEnterCommand(args []string) {
	if len(args) > 1 {
		ui.changeTopic(args[1])
	} else {
		ui.displaySystemMessage("Usage: /enter [new_topic_name]")
	}
}

func (ui *ChatUI) handleStatusCommand(args []string) {
	if len(args) > 1 {
		switch args[1] {
		case "sub":
			ui.displaySystemMessage(ui.getStatusSub())
		case "topic":
			ui.displaySystemMessage(ui.getStatusTopic())
		case "host":
			ui.displaySystemMessage(ui.getStatusHost())
		default:
			ui.displaySystemMessage("Unknown status type: " + args[1])
		}
	} else {
		ui.displaySystemMessage("Usage: /status [sub|topic|host]")
	}
}

func (ui *ChatUI) handleChatMessage(input string) error {
	// Wrapping the string message into the msg.Message structure

	msgBytes, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("message serialization error: %s", err)
	}

	msg, err := msg.New(ui.r.Entity.DID.Fragment, ui.r.Entity.DID.Fragment, string(msgBytes), "application/json")
	if err != nil {
		return fmt.Errorf("message creation error: %s", err)
	}

	// FIXME. This should be done in the message.New function
	m, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("message serialization error: %s", err)
	}

	err = ui.r.Public.Publish(ui.ctx, m)
	if err != nil {
		return fmt.Errorf("publish error: %s", err)
	}
	ui.displaySelfMessage(input)

	return nil
}

// Remaining methods like refreshPeers, displayChatMessage, etc., stay unchanged

// withColor wraps a string with color tags for display in the messages text box.
func withColor(color, msg string) string {
	return fmt.Sprintf("[%s]%s[-]", color, msg)
}

// displayStatusMessage writes a status message to the message window.
func (ui *ChatUI) displayStatusMessage(msg string) {
	prompt := withColor("cyan", "<STATUS>:")
	fmt.Fprintf(ui.msgW, "%s %s\n", prompt, msg)
}

func (ui *ChatUI) getStatusSub() string {
	// Return whatever status you'd like about the subscription.
	// Just an example below:
	return fmt.Sprintf("Subscription Topic: %s", ui.r.Public.String())
}

func (ui *ChatUI) getStatusTopic() string {
	// Return whatever status you'd like about the topic.
	// Fetching peers as an example below:
	peers := ui.r.Public.ListPeers()
	return fmt.Sprintf("Topic Status: %s | Peers: %v", ui.r.Entity.DID.Fragment, peers)
}

func (ui *ChatUI) getStatusHost() string {
	// Return whatever status you'd like about the host.
	// Just an example below:
	var result string
	result += "Peer ID: " + ui.n.ID().String() + "\n"
	result += "Peers:\n"
	for _, p := range ui.n.Network().Peers() {
		result += p.String() + "\n"
	}
	return result
}

func (ui *ChatUI) triggerDiscovery() {

	// go ui.n.StartPeerDiscovery(ui.ctx, config.GetRendezvous())
	p2p.StartPeerDiscovery(ui.ctx, ui.n)
	ui.displaySystemMessage("Discovery process started...")
}

func (ui *ChatUI) displaySystemMessage(msg string) {
	prompt := withColor("cyan", "[SYSTEM]:")
	fmt.Fprintf(ui.msgW, "%s %s\n", prompt, msg)
}

func (ui *ChatUI) changeTopic(d string) {

	// Create a new Room instance with the new topic
	room, err := room.New(ui.a)
	if err != nil {
		ui.displaySystemMessage(fmt.Sprintf("Failed to enter room '%s': %s", d, err))
		return
	}
	// If successful, assign the new Room instance to ui.cr
	ui.r = room
	// ui.msgW.SetTitle(fmt.Sprintf("Topic: %s", newTopic))
	ui.msgBox.SetTitle("Room: " + ui.r.Entity.DID.Fragment)

	// Notify the user
	ui.displaySystemMessage(fmt.Sprintf("Entered the new Room: %s", ui.r.Entity.DID.Fragment))

	ui.app.Draw()
}
