package actor

import (
	"fmt"
	"strings"

	"github.com/bahner/go-ma/msg"
)

var ErrInvalidDotContentType = fmt.Errorf("actor: invalid content type for dot message. Must be '%s'", msg.DEFAULT_CONTENT_TYPE)

func (a *Actor) defaultMessageHandler(m *msg.Message) error {

	switch m.Type {
	case msg.REQUEST:
		return a.handleAtMessage(m)
	default:
		return msg.ErrInvalidMessageType
	}
}

func (a *Actor) handleAtMessage(m *msg.Message) error {

	// Only receive messages with default content type
	if m.ContentType != msg.DEFAULT_CONTENT_TYPE {
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
	case "location":
		return a.HandleLocationMessage(m)
	default:
		return fmt.Errorf("actor: unknown dot command: %s", cmd)
	}
}
