package actor

import (
	"fmt"

	"github.com/bahner/go-ma/did"
	"github.com/bahner/go-ma/did/doc"
	"github.com/bahner/go-ma/key"
	log "github.com/sirupsen/logrus"
)

// Creates a ne DID Document for the entity. This only applies if
// the entity has a keyset. If the controller is "", the entity
// is the controller alone.
// NB! All entities must have a document, hence we apply this to the Entity struct,
// but we can only create Documents for actors. For entities we fetch them.
func (a *Actor) CreateDocument(controller string) (*doc.Document, error) {

	id := a.Entity.DID.String()

	if controller == "" {
		controller = id
	}

	if a.Keyset == nil {
		return nil, fmt.Errorf("entity: no keyset for entity %s", id)
	}

	// Initialize a new DID Document

	myDoc, err := doc.New(id, controller)
	if err != nil {
		return nil, fmt.Errorf("doc/GetOrCreate: failed to create new document: %w", err)
	}

	// Add the encryption key to the document,
	// and set it as the key agreement key.
	log.Debugf("entity/document: existing keyAgreement: %v", myDoc.KeyAgreement)
	myEncVM, err := doc.NewVerificationMethod(
		id,
		id,
		key.KEY_AGREEMENT_KEY_TYPE,
		did.GetFragment(a.Keyset.EncryptionKey.DID),
		a.Keyset.EncryptionKey.PublicKeyMultibase)
	if err != nil {
		return nil, fmt.Errorf("entity/document: failed to create encryption verification method: %s", err)
	}
	// Add the controller to the verification method
	err = myEncVM.AddController(controller)
	if err != nil {
		return nil, fmt.Errorf("entity/document: failed to add controller to encryption verification method: %s", err)
	}

	// Set the key agreement key verification method
	err = myDoc.AddVerificationMethod(myEncVM)
	if err != nil {
		return nil, fmt.Errorf("entity/document: failed to add encryption verification method to document: %s", err)
	}

	myDoc.KeyAgreement = myEncVM.ID
	log.Debugf("entity/document: set keyAgreement to %v for %s", myDoc.KeyAgreement, myDoc.ID)

	// Add the signing key to the document and set it as the assertion method.
	log.Debugf("entity/document: Creating assertionMethod for document %s", myDoc.ID)
	mySignVM, err := doc.NewVerificationMethod(
		id,
		id,
		key.ASSERTION_METHOD_KEY_TYPE,
		did.GetFragment(a.Keyset.SigningKey.DID),
		a.Keyset.SigningKey.PublicKeyMultibase)
	if err != nil {
		return nil, fmt.Errorf("entity: failed to create signing verification method: %s", err)
	}
	// Add the controller to the verification method if applicable
	err = mySignVM.AddController(controller)
	if err != nil {
		return nil, fmt.Errorf("entity: failed to add controller to signing verification method: %s", err)
	}

	// Set the assertion method verification method
	err = myDoc.AddVerificationMethod(mySignVM)
	if err != nil {
		return nil, fmt.Errorf("entity: failed to add signing verification method to document: %s", err)
	}

	myDoc.AssertionMethod = mySignVM.ID
	log.Debugf("entity/document: Set assertionMethod to %v for %s", myDoc.AssertionMethod, mySignVM.ID)

	// Finally sign the document with the signing key.
	err = myDoc.Sign(a.Keyset.SigningKey, mySignVM)
	if err != nil {
		return nil, fmt.Errorf("entity: failed to sign document: %s", err)
	}

	return myDoc, nil

}

// Creates a new DID Document for the entity, and sets it as the entity's document.
// This only applies if the entity has a keyset. If the controller is "", the entity
// is the controller alone.
func (a *Actor) CreateAndSetDocument(controller string) error {

	doc, err := a.CreateDocument(controller)
	if err != nil {
		return fmt.Errorf("entity: failed to create document: %s", err)
	}

	a.Entity.Doc = doc
	return nil

}
