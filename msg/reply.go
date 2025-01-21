package msg

import (
	"crypto/ed25519"
	"mime"

	"github.com/bahner/go-ma-actor/entity"
	"github.com/bahner/go-ma/msg"
	"github.com/fxamacker/cbor/v2"
)

const REPLY_SERIALIZATION = "cbor"

var REPLY_CONTENT_TYPE_PARAMS = map[string]string{
	"type": "reply",
}

type ReplyContent struct {
	// Id of the messagew to reply to
	RequestID string `cbor:"requestID"`
	// Reply content
	Reply []byte `cbor:"reply"`
}

func NewReply(m msg.Message, reply []byte) ([]byte, error) {
	return cbor.Marshal(
		&ReplyContent{
			RequestID: m.Id,
			Reply:     reply,
		})
}

// Reply to a message. requires the message to create a reply containing the id of the requesting message
// The message is not a pointer, as we only need the ID and then throw it away.
func Reply(m msg.Message, replyBytes []byte, privKey ed25519.PrivateKey) error {

	replyContent, err := NewReply(m, replyBytes)
	if err != nil {
		return err
	}

	mimeType := msg.CONTENT_TYPE + "+" + REPLY_SERIALIZATION
	contentType := mime.FormatMediaType(mimeType, REPLY_CONTENT_TYPE_PARAMS)

	reply, err := msg.New(m.To, m.From, contentType, replyContent, privKey)
	if err != nil {
		return err
	}

	err = reply.Sign(privKey)
	if err != nil {
		return err
	}

	envelope, err := reply.Enclose()
	if err != nil {
		return err
	}

	recipient, err := entity.GetOrCreate(m.From)
	if err != nil {
		return err
	}

	return recipient.Publish(envelope)
}
