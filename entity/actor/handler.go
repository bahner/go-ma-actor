package actor

import (
	"fmt"
	"strings"

	actormsg "github.com/bahner/go-ma-actor/msg"
	"github.com/bahner/go-ma/msg"
)

var ErrInvalidDotContentType = fmt.Errorf("actor: invalid content type for dot message. Must be '%s'", msg.CONTENT_TYPE)

func (a *Actor) defaultMessageHandler(m *msg.Message) error {

	messageType, err := m.MessageType()
	if err != nil {
		return err
	}

	switch messageType {
	case actormsg.CHAT_MESSAGE_TYPE:
		return a.handleAtMessage(m)
	default:
		return msg.ErrInvalidMessageType
	}
}

func (a *Actor) handleAtMessage(m *msg.Message) error {

	// Only receive messages with default content type
	if m.ContentType != msg.CONTENT_TYPE {
		return ErrInvalidDotContentType
	}

	var cmd string

	msgStr := strings.TrimPrefix(string(m.Content), "@")
	elements := strings.Split(msgStr, " ")

	if len(elements) == 0 {
		return fmt.Errorf("actor: empty dot message")
	}

	cmd = elements[0]

	switch cmd {
	default:
		return fmt.Errorf("actor: unknown dot command: %s", cmd)
	}
}
