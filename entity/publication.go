package entity

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/bahner/go-ma/api"
	"github.com/bahner/go-ma/did"
	"github.com/ipfs/boxo/ipns"
	"github.com/ipfs/boxo/path"
	"github.com/ipfs/go-cid"
	caopts "github.com/ipfs/kubo/core/coreiface/options"
	log "github.com/sirupsen/logrus"
)

var ErrAlreadyPublished = errors.New("Document is already published")

type PublishOptions struct {
	Ctx           context.Context
	Pin           bool
	Force         bool
	AllowBigBlock bool
}

func DefaultPublishOptions() *PublishOptions {
	return &PublishOptions{
		Ctx:           context.Background(),
		Pin:           true,
		Force:         true,
		AllowBigBlock: false,
	}
}

// Publish publishes the document to IPFS and returns the CID
// If the opts is nil, the default options are used.
// NB! Publication is more than simply adding the document to IPFS.
// It's about publishing the document to IPNS and possibly pinning it.
func (d *Document) Publish(o ...*PublishOptions) (ipns.Name, error) {

	alreadyPublishedString := "'to' cid was already recursively pinned"
	opts := DefaultPublishOptions()

	// Iterate through all options provided
	for _, opt := range o {
		if opt == nil {
			continue // Skip any nil options
		}
		if opt.Force {
			opts.Force = opt.Force
		}
		if opt.Pin {
			opts.Pin = opt.Pin
		}
		if opt.AllowBigBlock {
			opts.AllowBigBlock = opt.AllowBigBlock
		}
	}
	_did, err := did.New(d.ID)
	if err != nil {
		return ipns.Name{}, fmt.Errorf("DocPublish: %w", err)
	}

	// Make sure a key is available for the document
	ik, err := did.GetOrCreate(_did.Fragment)
	if err != nil {
		return ipns.Name{}, fmt.Errorf("DocPublish: %w", err)
	}

	data, err := d.MarshalToCBOR()
	if err != nil {
		return ipns.Name{}, fmt.Errorf("DocPublish: %w", err)
	}

	// Actually add the document to IPFS and possibly pin it and allow bib blocks.
	c, err := api.IPFSDagPutCBOR(data, opts.Pin, opts.AllowBigBlock)
	if err != nil {
		return ipns.Name{}, fmt.Errorf("DocPublish: %w", err)
	}
	p := path.FromCid(c)

	// Get the IPFS API
	a := api.GetIPFSAPI()

	// If an existing document is already published and Pin is set we need to update the existing pin f asked to force.
	e, err := d.CID()
	if err == nil && c != e && opts.Pin && opts.Force {
		err := a.Pin().Update(opts.Ctx, path.FromCid(e), p)
		if err != nil {
			if err.Error() == alreadyPublishedString {
				return ipns.Name{}, fmt.Errorf("DocPublish: %w", ErrAlreadyPublished)
			}
			return ipns.Name{}, fmt.Errorf("DocPublish: %w", err)
		}
	} else {
		log.Debugf("DocPublish: Document %s is not yet published.", d.ID)
	}

	log.Debugf("DocPublish: Announcing publication of document %s to IPNS. Please wait ...", c.String())
	n, err := a.Name().Publish(opts.Ctx, p, caopts.Name.Key(ik.Fragment))
	if err != nil {

		return ipns.Name{}, fmt.Errorf("DocPublish: %w", err)
	}
	log.Debugf("DocPublish: Successfully announced publication of document %s to %s", c.String(), n.AsPath().String())
	return n, nil

}

// The Goroutine version of Publish must get a cancel function as argument.
// This is to force the caller to use a context with a cancel function.
// Obviously this should probably be a timeout context.
// Other than that it is the same as Publish.
func (d *Document) PublishGoroutine(wg *sync.WaitGroup, cancel context.CancelFunc, opts *PublishOptions) {

	defer cancel()
	defer wg.Done()

	d.Publish(opts)

}

func (d *Document) CID() (cid.Cid, error) {

	maybeDoc, c, err := Fetch(d.ID, false) // Don't accept cached document
	if err != nil {
		log.Debugf("DocumentPublishGoroutine: %v", err)
		return cid.Cid{}, err
	}

	if maybeDoc == nil {
		log.Debugf("DocumentPublishGoroutine: Document is nil")
		return cid.Cid{}, fmt.Errorf("Document is nil")
	}

	return c, nil
}
