package main

import (
	"context"
	"fmt"
)

func main() {
	ctx := context.Background()

	initConfig()

	actor, err := newActor(ctx, keyset)
	if err != nil {
		panic(fmt.Sprintf("Failed to create actor: %v", err))
	}

	// Create and join the chat room. The "room" is essentially the topic name, so we'll use the IPNS name.
	room := ipnsName.String()
	r, err := newRoom(ctx, ps, nick, room)
	if err != nil {
		panic(err)
	}

	// Draw the UI.
	ui := NewChatUI(ctx, r, actor)
	if err := ui.Run(); err != nil {
		fmt.Errorf("error running text UI: %s", err)
	}
}
