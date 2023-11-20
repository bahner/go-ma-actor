package room

import (
	"github.com/bahner/go-home/actor"
)

type Room struct {
	*actor.Actor
}

func New(a *actor.Actor) (*Room, error) {

	r := &Room{Actor: a}

	return r, nil

}
