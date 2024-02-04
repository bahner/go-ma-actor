package entity

import (
	"context"
	"fmt"

	p2ppubsub "github.com/libp2p/go-libp2p-pubsub"
	log "github.com/sirupsen/logrus"
)

func (e *Entity) Subscribe(a *Entity) (chan *p2ppubsub.Message, error) {

	// I believe only the entity context is need for cancellation.
	// Everything else should be backgroud processes. If the pubsub
	// is broken, the entity is broken.
	ctx := context.Background()

	sub, err := e.Topic.Subscribe()
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	messageChan := make(chan *p2ppubsub.Message)

	// Start a goroutine to receive messages from sub.Next and pass them to messageChan
	go func() {
		for {
			select {
			case <-e.Ctx.Done():
				tn := sub.Topic()
				log.Errorf("entity/subscribe: Entity context done. Closing the topic %s", tn)
				sub.Cancel()
				log.Infof("entity/subscribe: Closed topic subscrption %s", tn)
				return
			default:
				message, err := sub.Next(ctx)
				if err != nil {
					continue
				}
				messageChan <- message
			}
		}
	}()

	return messageChan, nil
}
