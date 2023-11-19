package room

import (
	"fmt"

	"github.com/bahner/go-home/actor"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
)

type Room struct {
	*actor.Actor
}

func New(a *actor.Actor) (*Room, error) {

	r := &Room{Actor: a}

	return r, nil

}

// This is a very spcific function that is only used in the home package.
func (r *Room) Enter(ps *pubsub.PubSub, a *actor.Actor) error {

	var err error

	roomTopic := r.Public.String()

	a.Public, err = ps.Join(roomTopic)
	if err != nil {
		return fmt.Errorf("home: %v failed to join room: %v", a, err)
	}

	return nil

}
