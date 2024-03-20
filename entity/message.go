package entity

import "github.com/bahner/go-ma/msg"

type Message struct {
	Message   *msg.Message
	Enveloped bool
}

// Wraps a message and a boolean indicating if it was enveloped.
// Enveloped messages should probably be sent directly to the sender,
// whereas non-enveloped messages should be sent to the entity it was received from.
func NewMessage(m *msg.Message, enveloped bool) *Message {
	return &Message{
		Message:   m,
		Enveloped: enveloped,
	}
}
