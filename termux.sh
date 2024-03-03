#!/bin/bash

set -eu

trap cleanup EXIT

BINDIR="${PREFIX}/bin"
WGET_OPTS="-qO"

KUBO_VERSION="${KUBO_VERSION:-v0.26.0}"
KUBO_TARBALL="kubo_${KUBO_VERSION}_linux-arm64.tar.gz"
KUBO_URL="https://dist.ipfs.tech/kubo/${KUBO_VERSION}/${KUBO_TARBALL}"

GO_MA_ACTOR_VERSION="${GO_MA_ACTOR_VERSION:-v0.2.3}"
GO_MA_ACTOR="go-ma-actor-android-arm64"
GO_MA_ACTOR_URL="https://github.com/bahner/go-ma-actor/releases/download/${GO_MA_ACTOR_VERSION}/${GO_MA_ACTOR}.xz"

# These are for generating config
export GO_MA_ACTOR_ACTOR_NICK="${GO_MA_ACTOR_ACTOR_NICK:-termux}"
export GO_MA_ACTOR_API_MADDR="/ip4/127.0.0.1/tcp/5001"

cleanup()  {
  rm -rf "${GO_MA_ACTOR}"*
  rm -rf i"${KUBO_TARBALL}"*
}

install_wget() {
  pkg install wget
}

install_ipfs() {
    wget "${WGET_OPTS}" "${KUBO_TARBALL}" "${KUBO_URL}"

    tar xf "${KUBO_TARBALL}"
    mv kubo/ipfs "${BINDIR}"

    rm -rf "${KUBO_TARBALL}" kubo

    ipfs init || true

} 

run_ipfs() {
    ipfs daemon &> /dev/null || true &
}

install_actor() {
    wget "${WGET_OPTS}" "${GO_MA_ACTOR}.xz" "${GO_MA_ACTOR_URL}"
    xz -d "${GO_MA_ACTOR}.xz"

    chmod +x "${GO_MA_ACTOR}"
    mv "${GO_MA_ACTOR}" "${BINDIR}/actor"
}

generate_actor() {
  actor --generate --publish || true # No --force. Idempotency, you know.
}

# Install ipfs if not already installed
ipfs version &> /dev/null || install_ipfs
run_ipfs
test -x "${BINDIR}/wget" || install_wget
install_actor
actor
