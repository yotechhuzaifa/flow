// Package app is the composition root: it knows every layer and wires them
// together. Nothing else in the repo is allowed to import both transport and
// storage. This is the only place a dependency direction is decided.
package app

import (
	"log/slog"
	"path/filepath"

	"flow/internal/domain/flow"
	"flow/internal/exec"
	"flow/internal/nodes"
	"flow/internal/platform"
	"flow/internal/sqlite"
)

type App struct {
	Store   *sqlite.Store
	Exec    *exec.Executor
	DBPath  string
	DataDir string
	Emitter func(event string, payload any) // set by the transport layer
	Version string
}

// Ports answers "what ports does this node type expose", for validation.
func Ports(nodeType string) []string { return nodes.OutPorts(nodeType) }

// New builds the whole app graph. Call Open before, Close after.
func New(version string, emit func(string, any)) (*App, error) {
	dir, err := platform.DataDir("flow")
	if err != nil {
		return nil, err
	}
	dbPath := filepath.Join(dir, "app.db")

	st, err := sqlite.Open(dbPath)
	if err != nil {
		return nil, err
	}
	st.Ports = flow.PortNamer(Ports)

	nodes.Register(nodes.Send{}, nodes.Log{}, nodes.Delay{}, nodes.If{})

	return &App{
		Store:   st,
		Exec:    exec.New(emit),
		DBPath:  dbPath,
		DataDir: dir,
		Emitter: emit,
		Version: version,
	}, nil
}

func (a *App) Close() error { return a.Store.Close() }

// SetEmitter swaps where run events go (Wails events, SSE, test slice).
func (a *App) SetEmitter(fn func(string, any)) {
	a.Emitter = fn
	a.Exec = exec.New(fn)
}

// Logf is what nodes use to write to the app log.
func (a *App) Logf(format string, args ...any) { slog.Info(format, args...) }
