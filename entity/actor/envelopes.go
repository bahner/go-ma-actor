package actor

import (
	"context"
	"fmt"

	"github.com/bahner/go-ma-actor/entity"
	log "github.com/sirupsen/logrus"
)

// Lets actor listen for incoming envelopes and open them, if possible. The actor can decrypt the
// messages, which are then sent to the intended recipient. This is a little extra step., but it
// makes the design simpler and easier to understand /methinks.
// Sent the messages to the message channel in the input. Then you can handle the messages in the
// UI or wherever you want to handle them and you know they were private messages.
func (a *Actor) HandleIncomingEnvelopes(ctx context.Context, messages chan *entity.Message) {

	mesg := fmt.Sprintf("Handling incoming envelopes to " + a.Entity.DID.Id)
	log.Info(mesg)

	for {
		select {
		case <-ctx.Done():
			log.Debug("actor.HandleIncomingEnvelopes: Context cancelled, exiting handleIncomingEnvelopes...")
			return
		case envelope, ok := <-a.Envelopes:
			if !ok {
				log.Debug("actor.HandleIncomingEnvelopes: Actor envelope channel closed, exiting...")
				return
			}
			log.Debugf("actor.HandleIncomingEnvelopes: Received actor envelope: %v", envelope)

			err := a.Keyset.Verify()
			if err != nil {
				log.Errorf("actor.HandleIncomingEnvelopes: : %s: %v", a.Entity.DID.Id, err)
				continue
			}

			// Process the envelope and send a pong response
			m, err := envelope.Open(a.Keyset.EncryptionKey.PrivKey[:])
			if err != nil {
				log.Errorf("actor.HandleIncomingEnvelopes: Error opening actor envelope: %v\n", err)
				continue
			}

			// Deliver message to the requested message channel.
			log.Debugf("actor.HandleIncomingEnvelopes: Accepted message %s from %s to %s", m.Id, m.From, m.To)
			messages <- entity.NewMessage(m, true)
		}
	}
}
