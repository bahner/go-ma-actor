package ui

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/bahner/go-ma/msg"
	"github.com/bahner/go-ma/p2p"
)

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
	peers := ui.n.Network().Peers()

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
	prompt := withColor("yellow", fmt.Sprintf("<%s>:", ui.room.Entity.DID.Fragment))
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

		case m := <-ui.actor.Messages:
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

	msg, err := msg.New(ui.room.Entity.DID.Fragment, ui.room.Entity.DID.Fragment, string(msgBytes), "application/json")
	if err != nil {
		return fmt.Errorf("message creation error: %s", err)
	}

	// FIXME. This should be done in the message.New function
	m, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("message serialization error: %s", err)
	}

	err = ui.room.Public.Publish(ui.ctx, m)
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
	return fmt.Sprintf("Subscription Topic: %s", ui.room.Public.String())
}

func (ui *ChatUI) getStatusTopic() string {
	// Return whatever status you'd like about the topic.
	// Fetching peers as an example below:
	peers := ui.room.Public.ListPeers()
	return fmt.Sprintf("Topic Status: %s | Peers: %v", ui.room.Entity.DID.Fragment, peers)
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

	ui.actor.Enter(d)

	ui.msgBox.SetTitle("Room: " + ui.room.Entity.DID.Fragment)

	// Notify the user
	ui.displaySystemMessage(fmt.Sprintf("Entered the new Room: %s", ui.room.Entity.DID.Fragment))

	ui.app.Draw()
}
