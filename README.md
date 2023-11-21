# github.com/bahner/go-ma-actor

This is go-ma-actor based on [an example from go-libp2p][src].

Now you can either run with `go run`, or build and run the binary:

```bash
# Generate persistent environment variables of *SECRET* keysets
eval $(go run . -genenv -forcePublish | tee .env)
./go-ma-actor # Share and enjoy!
```

## Configuration

type `./go-ma-actor -help`. Most config settings can be set with environment variables, as follows:

```bash
export GO_ACTOR_LOG_LEVEL="error"
export GO_ACTOR_DISCOVERY_TIMEOUT="300"
export GO_ACTOR_KEYSET="myBase58EncodedPrivkeyGeneratedByGenerate"
```

## Identity

A `-generate` or `genenv` parameter to generate a text version of a secret key.
The key is text formatted privKey for your node.

This key can and should be kept safely on a PostIt note on your monitor :-)
Just don't store somewhere insecure. It's your future identity.

```bash
unset HISTFILE
 export GO_ACTOR_KEYSET=FooBarABCDEFbase58
```

or specified on the command line:

```bash
./go-ma-actor -keyset FooBarABCDEFbase58
```

The first is the best. (Noticed that in most shells the empty space before the command, means that the line isn't saved in history.)

## Usage

To quit, hit `Ctrl-C`, or type `/quit` into the input field.

## Commands

- /status [sub|topic|host]
- /discover
- /enter room
- /nick Name
- /refresh

[src]: https://github.com/libp2p/go-libp2p/tree/master/examples/pubsub/chat
