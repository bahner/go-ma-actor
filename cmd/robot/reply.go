package main

import (
	"context"

	"github.com/bahner/go-ma/msg"
	log "github.com/sirupsen/logrus"
)

func reply(m *msg.Message) []byte {

	c := client()

	ctx := context.Background()
	res, err := c.SimpleSend(ctx, string(m.Content))
	if err != nil {
		log.Fatal(err)
	}

	return []byte(res.Choices[0].Message.Content)
}
