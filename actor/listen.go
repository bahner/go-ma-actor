package actor

import (
	"fmt"
)

func (a *Actor) Listen(outputChannel chan<- string) error {
	// Subscribe to Inbox topic
	inboxSub, err := a.Inbox.Subscribe()
	if err != nil {
		return fmt.Errorf("failed to subscribe to Inbox topic: %v", err)
	}
	defer inboxSub.Cancel()

	// Start a goroutine for Inbox subscription
	go a.handlePrivateMessages(inboxSub)

	// Wait for context cancellation (or other exit conditions)
	<-a.ctx.Done()
	return a.ctx.Err()
}
