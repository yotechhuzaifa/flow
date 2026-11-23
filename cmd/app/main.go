//go:build !server

package main

import (
	"log"

	"flow/internal/wails"
)

// version is injected by the build: -ldflags "-X main.version=1.2.3"
var version = "dev"

func main() {
	a, err := wails.Build(version)
	if err != nil {
		log.Fatal(err)
	}
	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
