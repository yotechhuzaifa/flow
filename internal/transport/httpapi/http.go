// Package httpapi is the server-mode adapter: same domain, no webview.
// Used by `cmd/server` and by the headless parity suite.
package httpapi

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"flow/internal/app"
	"flow/internal/domain/flow"
)

type sseBroker struct {
	mu      sync.Mutex
	clients map[chan []byte]struct{}
}

func newSSEBroker() *sseBroker { return &sseBroker{clients: map[chan []byte]struct{}{}} }

func (b *sseBroker) emit(event string, payload any) {
	body, _ := json.Marshal(map[string]any{"type": event, "data": payload})
	b.mu.Lock()
	defer b.mu.Unlock()
	for ch := range b.clients {
		select {
		case ch <- body:
		default: // slow client: drop rather than block a run
		}
	}
}

func (b *sseBroker) serveHTTP(w http.ResponseWriter, r *http.Request) {
	fl, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")

	ch := make(chan []byte, 64)
	b.mu.Lock()
	b.clients[ch] = struct{}{}
	b.mu.Unlock()
	defer func() {
		b.mu.Lock()
		delete(b.clients, ch)
		b.mu.Unlock()
	}()

	for {
		select {
		case msg := <-ch:
			_, _ = w.Write([]byte("data: " + string(msg) + "\n\n"))
			fl.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

// Serve exposes the app over HTTP with SSE events, plus the built frontend.
func Serve(addr string, a *app.App) error {
	broker := newSSEBroker()
	a.SetEmitter(broker.emit)

	mux := http.NewServeMux()

	mux.HandleFunc("/api/flows", func(w http.ResponseWriter, r *http.Request) {
		list, err := a.Store.List()
		writeJSON(w, list, err)
	})

	mux.HandleFunc("/api/flow", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		switch r.Method {
		case http.MethodGet:
			f, err := a.Store.Get(id)
			writeJSON(w, f, err)
		case http.MethodPut:
			var f flow.Flow
			if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			writeJSON(w, map[string]bool{"ok": true}, a.Store.Save(&f))
		case http.MethodDelete:
			writeJSON(w, map[string]bool{"ok": true}, a.Store.Delete(id))
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/run", func(w http.ResponseWriter, r *http.Request) {
		f, err := a.Store.Get(r.URL.Query().Get("id"))
		if err != nil {
			writeJSON(w, nil, err)
			return
		}
		run, err := a.Exec.Start(context.Background(), f)
		if err != nil {
			writeJSON(w, nil, err)
			return
		}
		_ = a.Store.RecordRun(run.ID, run.FlowID, string(run.Status))
		writeJSON(w, map[string]string{"runId": run.ID, "status": string(run.Status)}, nil)
	})

	mux.HandleFunc("/api/about", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, map[string]string{"version": a.Version, "dbPath": a.DBPath}, nil)
	})

	mux.HandleFunc("/events", broker.serveHTTP)
	mux.Handle("/", http.FileServer(http.Dir("frontend/dist")))

	log.Printf("flow server on http://%s", addr)
	return http.ListenAndServe(addr, mux)
}

func writeJSON(w http.ResponseWriter, v any, err error) {
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	_ = json.NewEncoder(w).Encode(v)
}
