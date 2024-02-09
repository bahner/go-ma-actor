# github.com/bahner/go-ma-actor

This is go-ma-actor based on [an example from go-libp2p][src].

## Requirements

This is a distributed app that relies heavily on the [libp2p](https://libp2p.io/) stack
and [IPFS][ipfs] in particular. It's unusable unless you have a running IPFS node.

I suggest using [Brave Browser][brave] or [IPFS Desktop][desktop] to run and IPFS node.

*By using Brave browser your can run an IPFS node without installing anything.
And you can investigate the IPFS network with the built-in IPFS node.
It provides The ability to browse IPFS properly, and to pin files and directories.*

## TL;DR

```bash

# Generate persistent config file with *SECRETS*
# The public parts needs published to the IPFS network to be useful, hence
# this takes a while. Sometimes 10 seconds, but also 3 minutes.
./go-ma-actor --generate --nick "asj" --publish > actor.yaml
./go-ma-actor # Share and enjoy!
```

## Configuration

type `./go-ma-actor -help`. Most config settings can be set with environment variables, as follows:

```bash
export GO_MA_LOG_LEVEL="error"
export GO_MA_LIBP2P_DISCOVERY_TIMEOUT="300"
export GO_MA_ACTOR_IDENTITY="myBase58EncodedPrivkeyGeneratedByGenerate"
```

## Identity

A `-generate` flag is available to generate a new identity.
It uses defaults, BUT it generates a new random identity.

You can use the output as your future identity, but keep it secret.
Those identities are used to sign messages, and to encrypt and decrypt private messages.

## Usage

To quit, hit `Ctrl-C`, or type `/quit` into the input field.

## Commands

- /status [sub|topic|host]
- /discover
- /alias [node|entity] set [DID|NAME] NAME
- /aliases
- /whereis [DID|NAME]
- /msg Name Message
- /enter room
- /refresh

[src]: https://github.com/libp2p/go-libp2p/tree/master/examples/pubsub/chat
[brave]: <https://brave.com/> (Recommended Browser for 間)
[desktop]: <https://docs.ipfs.tech/install/ipfs-desktop/> (IPFS Desktop)
[ipfs]: <https://ipfs.io/> (IPFS)
