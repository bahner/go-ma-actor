package ui

import (
	"context"
	"fmt"
	"io"

	"github.com/bahner/go-ma-actor/config"
	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma-actor/entity/actor"
	"github.com/bahner/go-ma-actor/p2p"
	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/rivo/tview"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	UI_MESSAGES_CHANNEL_BUFFERSIZE = 32
	PUBSUB_MESSAGES_BUFFERSIZE     = 32

	defaultLimbo = "closet"

	defaultPeerslistWidth = 20
	defaultHistorySize    = 50

	separator = " "
)

func init() {

	// UI
	pflag.Int("peerslist-width", defaultPeerslistWidth, "Sets the width of the peerslist pane in the UI")
	viper.BindPFlag("ui.peerslist-width", pflag.Lookup("peerslist-width"))

	pflag.Int("history-size", defaultHistorySize, "Sets the size of the history buffer in the UI.")
	viper.BindPFlag("ui.history-size", pflag.Lookup("history-size"))

}

// ChatUI is a Text User Interface (TUI) for a Room.
// The Run method will draw the UI to the terminal in "fullscreen"
// mode. You can quit with Ctrl-C, or by typing "/quit" into the
// chat prompt.
type ChatUI struct {
	p *p2p.P2P
	c config.Config

	// The actor is need to encrypt and sign messages in the event loop.
	a                  *actor.Actor
	currentActorCtx    context.Context
	currentActorCancel context.CancelFunc

	// The current entity is the "room" we are convering with.
	e *entity.Entity
	// Context for the current entity - NOT the actor!
	currentEntityCtx    context.Context
	currentEntityCancel context.CancelFunc

	// Broadcasts
	b               *p2ppubsub.Topic
	broadcastCtx    context.Context
	broadcastCancel context.CancelFunc

	// DiscoveryLoop
	discoveryLoopCtx    context.Context
	discoveryLoopCancel context.CancelFunc

	// History of entries
	inputField          *tview.InputField
	inputHistory        []string
	currentHistoryIndex int

	// The Topic is used for publication of messages after encryption and signing.
	// The names are obviously, from the corresponding DIDDocument.

	app       *tview.Application
	peersList *tview.TextView
	msgBox    *tview.TextView
	chatPanel *tview.Flex
	screen    *tview.Flex

	msgW       io.Writer
	chInput    chan string
	chMessages chan *entity.Message
	chDone     chan struct{}
}

// New returns a new ChatUI struct that controls the text UI.
// It won't actually do anything until you call Run().
// The enity is the "room" we are convering with.
func New(p *p2p.P2P, a *actor.Actor) (*ChatUI, error) {

	app := tview.NewApplication()

	ui := &ChatUI{
		a:          a,
		p:          p,
		app:        app,
		chInput:    make(chan string, 32),
		chMessages: make(chan *entity.Message, UI_MESSAGES_CHANNEL_BUFFERSIZE),
		chDone:     make(chan struct{}, 1),
	}

	// Start the broadcast subscription first, so
	// actors can announce themselves.
	fmt.Println("Initialising broadcasts...")
	ui.initBroadcast()

	// Since tview is global we can just run this which sets the style for the whole app.
	ui.setupApp()

	return ui, nil
}

// end signals the event loop to exit gracefully
func (ui *ChatUI) end() {
	ui.chDone <- struct{}{}
}
