//go:build !server

// Package wails adapts the app to the desktop shell. Services live here, never in cmd.
package wails

import (
	"context"
	"errors"

	"flow/internal/app"
	"flow/internal/domain/flow"
	"flow/internal/platform"
)

// Flows is the one service the frontend talks to.
// Method names become the generated TypeScript bindings.
type Flows struct {
	A   *app.App
	Ctx context.Context
}

func NewFlows(a *app.App) *Flows { return &Flows{A: a} }

func (s *Flows) ListFlows() ([]flow.Flow, error) { return s.A.Store.List() }

func (s *Flows) GetFlow(id string) (*flow.Flow, error) {
	if id == "" {
		return nil, errors.New("empty id")
	}
	return s.A.Store.Get(id)
}

func (s *Flows) SaveFlow(f *flow.Flow) error { return s.A.Store.Save(f) }

func (s *Flows) DeleteFlow(id string) error { return s.A.Store.Delete(id) }

// StartRun runs synchronously. Events stream out through the emitter while it runs.
// Returns the run id so the frontend can correlate events.
func (s *Flows) StartRun(flowID string) (string, error) {
	f, err := s.A.Store.Get(flowID)
	if err != nil {
		return "", err
	}
	ctx := s.Ctx
	if ctx == nil {
		ctx = context.Background()
	}
	r, err := s.A.Exec.Start(ctx, f)
	if err != nil {
		return "", err
	}
	_ = s.A.Store.RecordRun(r.ID, r.FlowID, string(r.Status))
	return r.ID, nil
}

// About is used by the frontend footer and by tests that want a smoke probe.
type About struct {
	Version string `json:"version"`
	OS      string `json:"os"`
	DBPath  string `json:"dbPath"`
}

func (s *Flows) About() About {
	return About{Version: s.A.Version, OS: platform.GOOS(), DBPath: s.A.DBPath}
}
