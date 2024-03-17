package config

import (
	"fmt"

	"github.com/bahner/go-ma/did/doc"
	"github.com/bahner/go-ma/key/set"
	log "github.com/sirupsen/logrus"
)

func PublishIdentityFromKeyset(k set.Keyset) error {

	d, err := doc.NewFromKeyset(k)
	if err != nil {
		return fmt.Errorf("config.publishIdentityFromKeyset: failed to create DOC: %w", err)
	}

	assertionMethod, err := d.GetAssertionMethod()
	if err != nil {
		return fmt.Errorf("config.publishIdentityFromKeyset: %w", err)
	}
	d.Sign(k.SigningKey, assertionMethod)

	_, err = d.Publish()
	if err != nil {
		return fmt.Errorf("config.publishIdentityFromKeyset: %w", err)

	}
	log.Debugf("Published identity: %s", d.ID)

	return nil
}
