package entity

import (
	"context"
	"fmt"
	"sync"

	"github.com/bahner/go-ma/did"
	"github.com/bahner/go-ma/did/doc"
	"github.com/bahner/go-ma/key"
	log "github.com/sirupsen/logrus"
)

// Creates a ne DID Document for the entity. This only applies if
// the entity has a keyset. If the controller is "", the entity
// is the controller alone.
func (e *Entity) CreateDocument(controller string) error {

	id := e.DID.String()

	if controller == "" {
		controller = id
	}

	if e.Keyset == nil {
		return fmt.Errorf("entity: no keyset for entity %s", e.DID.String())
	}

	// Initialize a new DID Document

	myDoc, err := doc.New(e.DID.String(), controller)
	if err != nil {
		return fmt.Errorf("doc/GetOrCreate: failed to create new document: %w", err)
	}

	// Add the encryption key to the document,
	// and set it as the key agreement key.
	log.Debugf("entity/document: existing keyAgreement: %v", myDoc.KeyAgreement)
	myEncVM, err := doc.NewVerificationMethod(
		id,
		id,
		key.KEY_AGREEMENT_KEY_TYPE,
		did.GetFragment(e.Keyset.EncryptionKey.DID),
		e.Keyset.EncryptionKey.PublicKeyMultibase)
	if err != nil {
		return fmt.Errorf("entity/document: failed to create encryption verification method: %s", err)
	}
	// Add the controller to the verification method
	err = myEncVM.AddController(controller)
	if err != nil {
		return fmt.Errorf("entity/document: failed to add controller to encryption verification method: %s", err)
	}

	// Set the key agreement key verification method
	err = myDoc.AddVerificationMethod(myEncVM)
	if err != nil {
		return fmt.Errorf("entity/document: failed to add encryption verification method to document: %s", err)
	}

	myDoc.KeyAgreement = myEncVM.ID
	log.Debugf("entity/document: set keyAgreement to %v for %s", myDoc.KeyAgreement, myDoc.ID)

	// Add the signing key to the document and set it as the assertion method.
	log.Debugf("entity/document: Creating assertionMethod for document %s", myDoc.ID)
	mySignVM, err := doc.NewVerificationMethod(
		id,
		id,
		key.ASSERTION_METHOD_KEY_TYPE,
		did.GetFragment(e.Keyset.SigningKey.DID),
		e.Keyset.SigningKey.PublicKeyMultibase)
	if err != nil {
		return fmt.Errorf("entity: failed to create signing verification method: %s", err)
	}
	// Add the controller to the verification method if applicable
	err = mySignVM.AddController(controller)
	if err != nil {
		return fmt.Errorf("entity: failed to add controller to signing verification method: %s", err)
	}

	// Set the assertion method verification method
	err = myDoc.AddVerificationMethod(mySignVM)
	if err != nil {
		return fmt.Errorf("entity: failed to add signing verification method to document: %s", err)
	}

	myDoc.AssertionMethod = mySignVM.ID
	log.Debugf("entity/document: Set assertionMethod to %v for %s", myDoc.AssertionMethod, mySignVM.ID)

	// Finally sign the document with the signing key.
	err = myDoc.Sign(e.Keyset.SigningKey, mySignVM)
	if err != nil {
		return fmt.Errorf("entity: failed to sign document: %s", err)
	}

	e.Doc = myDoc

	return nil

}

// Publish entity document. This needs to be done, when the keyset is new.
// Maybe we can check the assertionMethod and keyAgreement fields to see if
// the document is already published corretly.
func (e *Entity) PublishDocument() error {

	id, err := e.Doc.Publish(nil)
	if err != nil {
		return fmt.Errorf("entity: failed to publish document: %s", err)
	}
	log.Debugf("entity: published document: %v to %s", e.Doc, id)
	return nil
}

func (e *Entity) PublishDocumentGorutine(wg *sync.WaitGroup, cancel context.CancelFunc, opts *doc.PublishOptions) {
	defer wg.Done()
	defer cancel()

	done := make(chan struct{})
	go func() {
		e.Doc.Publish(opts)
		close(done)
	}()

	// Wait for the Publish operation to complete or for the context to be cancelled/timed out
	select {
	case <-opts.Ctx.Done():
		// Context is cancelled or timed out
		if opts.Ctx.Err() == context.DeadlineExceeded {
			log.Errorf("entity: deadline exceeded: %v", opts.Ctx.Err())
		} else {
			log.Errorf("entity: context cancelled: %v", opts.Ctx.Err())
		}
	case <-done:
		// Publish operation completed
		log.Infof("Published document for entity: %s", e.DID.String())
	}
}
