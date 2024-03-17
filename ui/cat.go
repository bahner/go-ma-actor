package ui

import (
	"context"
	"io"

	"github.com/bahner/go-ma/api"
	"github.com/ipfs/boxo/path"
)

const (
	catUsage = "/cat CID"
	catHelp  = "Fetches an object from IPFS and displays it"
)

func (ui *ChatUI) handleCatCommand(args []string) {

	if len(args) == 2 {
		cid := args[1]

		p, err := path.NewPath("/ipfs/" + cid)
		if err != nil {
			ui.displaySystemMessage("Invalid CID: " + err.Error())
			return
		}

		ia := api.GetIPFSAPI()

		r, err := ia.Unixfs().Get(context.Background(), p)
		if err != nil {
			ui.displaySystemMessage("Error fetching object: " + err.Error())
			return
		}

		f, ok := r.(io.ReadCloser)
		if !ok {
			ui.displaySystemMessage("Error fetching object: not a ReadCloser")
			return
		}
		// Read all data from the reader
		data, err := io.ReadAll(f)
		if err != nil {
			ui.displaySystemMessage("Error reading data: " + err.Error())
			return
		}

		ui.displaySystemMessage(string(data))
		return
	}

	ui.handleHelpCommand(catUsage, catHelp)

}
