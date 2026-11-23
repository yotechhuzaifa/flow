//go:build server

// Server mode: the same domain, no webview. Useful for headless CI runs of the
// parity suite and for the future hosted version.
package main

import (
	"log"
	"os"

	"flow/internal/app"
	"flow/internal/transport/httpapi"
)

func main() {
	addr := os.Getenv("FLOW_ADDR")
	if addr == "" {
		addr = "127.0.0.1:8787"
	}
	emit := func(event string, payload any) {
		log.Printf("event %s %v", event, payload)
	}
	a, err := app.New("server", emit)
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()

	if err := httpapi.Serve(addr, a); err != nil {
		log.Fatal(err)
	}
}
