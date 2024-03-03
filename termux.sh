#!/bin/bash

cd || exit

export PATH="${PATH}:.local/bin"

KUBO_VERSION="${KUBO_VERSION:-v0.26.0}"
KUBO_TARBALL="kubo_${KUBO_VERSION}_linux-arm64.tar.gz"
KUBO_URL="https://dist.ipfs.tech/kubo/${KUBO_VERSION}/${KUBO_TARBALL}"

GO_MA_ACTOR_VERSION="${GO_MA_ACTOR_VERSION:-v0.2.3}"
GO_MA_ACTOR="go-ma-actor-android-arm64"
GO_MA_ACTOR_URL="https://github.com/bahner/go-ma-actor/releases/download/${GO_MA_ACTOR_VERSION}/${GO_MA_ACTOR}.xz"
pkg up

### IPFS

# Fetch
curl  "${KUBO_URL}" -o "${KUBO_TARBALL}"
tar xf "${KUBO_TARBALL}"

# Install
pushd kubo || exit
./install.sh
popd || exit

# Configure
ipfs init

# Run
ipfs daemon &> /dev/null &

### GO_MA_ACTOR

# Fetch
curl "${GO_MA_ACTOR_URL}" -o "${GO_MA_ACTOR}.xz"
xz -d "${GO_MA_ACTOR}.xz"

# Install
chmod +x "${GO_MA_ACTOR}"
mv "${GO_MA_ACTOR}" .local/bin/actor

cat <<EOF
# If you haven't already, now run
actor --api-maddr /ip4/127.0.0.1/tcp/5001 --generate --nick ${USER} --publish
actor --nick ${USER}
EOF
