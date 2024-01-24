# github.com/bahner/go-ma-relay

This is a simple libp2p node that runs in server mode. It's meant to just be started
and added to the bootstrap list of your clients. It will keep searching for peers that
advertise the same rendezvous string and connect to them. It will also advertise itself
to other peers that are looking for the same rendezvous string.

It will act as a relay for your clients. It will not relay any traffic that is not
destined for a peer that has the same rendezvous string.

The PeerID along with it's multiAddrs are found on the generated web page.

## TL;DR

```bash
./go-ma-relay # Just run it
```

## Keyset

You probably want a keyset to use with this. One will be generated for you if you don't provide one. You can pipe this directly to an enviornment file, eg:

```bash
./go-ma-relay -generate-keyset > .env
# Then for the next run
. .env
./go-ma-relay
```

## Docker

Docker images are provided. In order to run in docker you have to run in host networking mode. You can experiment with setting -listenPort and exposing that, I guess.

```bash
./go-ma-relay -generate-keyset > .env
docker-compose up 
```

or if you're a hardCoreCoder, run it from the command line:

```bash
# Using host networking directly
docker run --network host --env-file .env bahner/go-ma-relay -keyset kwoo....blahblahblah

# Exposing distinct ports.
# This'll probably not works as you think if you're behind a NAT'ed firewall.
# I suggest running in host network mode
docker run -p 4000-4001:4000-4001 bahner/go-ma-relay 
```

## Configuration

You can configure your settings as command line parameters or as environment variables. The following variables are recognised.

```bash
./go-ma-relay -help
export GO_MA_RELAY_DISCOVERY_SLEEP=10 # Sleep 10 seconds between discovery attempts
export GO_MA_RELAY_HTTP_ADDR="0.0.0.0" # Listen on all interfaces
export GO_MA_RELAY_HTTP_PORT="4000" # Listen on port 4000.
export GO_MA_RELAY_LISTEN_PORT="4001" # Listen on port 4001 for libp2p traffic. 0 = random
export GO_MA_RELAY_LOG_LEVEL="info" # Log level. debug, info, warn, error, fatal, panic
export GO_MA_RELAY_LOW_WATER_MARK=100 # Minimum number of connections to maintain
export GO_MA_RELAY_HIGH_WATER_MARK=1000 # Maximum number of connections to maintain 
export GO_MA_RELAY_CONN_MGR_GRACE_PERIOD=1 # 1 minute to let connections disconnect gracefully
export GO_MA_RELAY_KEYSET="" # Generated secret keyset. If not provided, one will be generated.
./go-ma-relay
```

2023-11-19 bahner
