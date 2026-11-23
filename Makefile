# Thin wrapper over wails3 + npm. The Wails build system stays authoritative;
# these targets are the ones the agent and CI use.
.PHONY: help setup dev bind build build-server test check e2e package clean doctor webui

APP      := flow
BIN      := bin
VERSION  ?= 0.1.0

help:            ## show this help
	@grep -E '^[a-z-]+:.*##' $(MAKEFILE_LIST) | sed 's/:.*##/\t/' | column -t -s $$'\t'

setup:           ## one-time toolchain check + deps
	bash scripts/setup.sh

dev:             ## wails3 dev: hot reload (frontend + go)
	wails3 dev -config ./build/config.yml -port 9245

bind:            ## regenerate TS bindings from Go services
	wails3 generate bindings -ts -clean=true -d frontend/bindings

build: bind      ## desktop build for the host OS
	wails3 task build PRODUCTION=true

build-server:    ## headless HTTP build (no webview), for CI and parity tests
	go build -tags server -o $(BIN)/flow-server ./cmd/server

test:            ## go tests + frontend typecheck
	go test -tags server ./... -count=1
	npm --prefix frontend run -s typecheck

check:           ## the fast gate the loop uses
	bash scripts/check.sh

webui:           ## build the frontend and sync it into internal/webui/dist
	npm --prefix frontend run -s build

e2e:             ## parity suite (starts vite; run npm install once at the repo root)
	npm --prefix frontend run -s build
	npx playwright test

package:         ## native installers for the current OS
	wails3 task package

clean:
	rm -rf $(BIN) frontend/dist

doctor:          ## what is missing on this machine
	wails3 doctor

# cross-compilation helpers (Wails v3 supports it through docker; see docs/platform-setup.md)
cross-mac:       ## from Linux/Windows
	wails3 task darwin:build GOOS=darwin
cross-win:
	wails3 task windows:build GOOS=windows
cross-linux:
	wails3 task linux:build GOOS=linux
