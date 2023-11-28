package main

import (
	"crypto"
	_ "crypto/sha512"
	"encoding/hex"
	"syscall/js"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	done := make(chan struct{})
	js.Global().Set("wasmHash", js.FuncOf(hash))
	<-done
}

func hash(this js.Value, args []js.Value) interface{} {
	log.Debugf("hashing: %s", args[0].String())
	h := crypto.SHA512.New()
	h.Write([]byte(args[0].String()))

	return hex.EncodeToString(h.Sum(nil))
}
