#!/usr/bin/env bash
# Installs a local Go toolchain and Go cache OUTSIDE the workspace, then verifies.
#
# Why: a Go SDK is ~1.4 GB and ~38k files, and a module cache is often larger.
# Written inside the workspace they blow the session snapshot budget
# (128 MB / 10k files) and your work silently stops being saved.
set -uo pipefail
root="$(cd "$(dirname "$0")/.." && pwd)"
TOOLDIR="${TOOLDIR:-$HOME/.local/share/gotoolchain}"
mkdir -p "$TOOLDIR"

ver() { "$TOOLDIR/go/bin/go" version 2>/dev/null | awk '{print $3}'; }

if [ "$(ver)" != "go1.26.0" ]; then
  echo "== installing go1.26.0 into $TOOLDIR (outside the workspace)"
  ( go_version=go1.26.0
    cd "$TOOLDIR"
    curl -fsSL "https://go.dev/dl/${go_version}.linux-amd64.tar.gz" -o go.tgz
    rm -rf go && tar -xzf go.tgz && rm go.tgz ) || exit 1
fi

export PATH="$TOOLDIR/go/bin:$PATH"
# keep every cache out of the workspace
export GOMODCACHE="${GOMODCACHE:-$HOME/.local/share/gomodcache}"
export GOCACHE="${GOCACHE:-$HOME/.local/share/gocache}"
export GOFLAGS="${GOFLAGS:--mod=mod}"
mkdir -p "$GOMODCACHE" "$GOCACHE"

cd "$root"
echo "go: $(go version)"
go vet -tags server ./...
go test -tags server ./... -count=1
go build -tags server ./...
echo "toolchain ok"
