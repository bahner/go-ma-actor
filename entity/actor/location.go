package actor

import (
	"fmt"

	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma/did"
)

const defaultLocationReply = "Location unknown"

var ErrNoLocation = fmt.Errorf("actor: no location set")

func (a *Actor) SetLocation(e *entity.Entity) error {

	err := e.Verify()
	if err != nil {
		return err
	}

	a.Location = e

	return nil
}

func (a *Actor) SetLocationFromDID(did did.DID) error {

	e, err := entity.New(did)
	if err != nil {
		return err
	}

	return a.SetLocation(e)
}

func (a *Actor) GetLocation() (string, error) {

	if a.Location == nil {
		return "", ErrNoLocation
	}

	return a.Location.DID.Id, nil
}

// func (a *Actor) HandleLocationMessage(m *msg.Message) error {

// 	ctx := context.Background()
// 	replyBytes := []byte(defaultLocationReply)

// 	// Set the reply to the currentLocation, if it is set.
// 	loc, err := a.GetLocation()
// 	if err == nil {
// 		replyBytes = []byte(loc)
// 	}

// 	e, err := entity.GetOrCreate(m.From)
// 	if err != nil {
// 		return fmt.Errorf("failed to get or create entity: %w", errors.Cause(err))
// 	}

// 	log.Debugf("Sending location to %s over %s", m.From, a.Entity.Topic.String())

// 	return actormsg.Reply(
// 		ctx,
// 		*m,
// 		replyBytes,
// 		a.Keyset.SigningKey.PrivKey,
// 		e.Topic,
// 	)

// }
