package actor

import (
	"github.com/bahner/go-ma/msg"
)

func (a *Actor) DefaultMessageHandler(m *msg.Message) error {

	switch m.Type {
	case msg.DOT:
		return handleDotMessage(m)
	case msg.BROADCAST:
		return handleBroadcastMessage(m)

	default:
		return msg.ErrInvalidType
	}
}

func handleDotMessage(m *msg.Message) error {
	// Do something with the message
	return nil
}

func handleBroadcastMessage(m *msg.Message) error {
	// Do something with the message
	return nil
}
