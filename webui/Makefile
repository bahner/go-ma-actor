#!/usr/bin/make -ef

CID = .cid
WASM = static/main.wasm
export GOOS=js
export GOARCH=wasm

default: clean $(WASM) publish

$(WASM):
	go build -o $(WASM) main.go

publish: $(CID)

$(CID): $(WASM)
	./.publish > $(CID)

open: $(CID)
	echo $$(cat $(CID)| xargs -I{} xdg-open ipfs://{})

clean:
	rm -f $(WASM) 

.PHONY: publish open clean tidy
