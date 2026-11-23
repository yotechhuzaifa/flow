// Package exec runs a flow. Framework-free: give it a Flow and an emitter, it runs headless.
package exec

import (
	"context"
	"fmt"
	"log/slog"
	"sync/atomic"
	"time"

	"flow/internal/domain/flow"
	"flow/internal/domain/run"
	"flow/internal/nodes"
)

// Emitter receives run.* and node.* events. Wired to Wails in internal/transport,
// to a slice in tests, to a websocket in server mode.
type Emitter func(event string, payload any)

// Executor runs flows sequentially. Concurrent runs are allowed; each gets its own state.
type Executor struct {
	emit Emitter
	seq  atomic.Int64
}

func New(emit Emitter) *Executor {
	if emit == nil {
		emit = func(string, any) {}
	}
	return &Executor{emit: emit}
}

// Start validates, runs in topological order, and returns the finished Run.
func (e *Executor) Start(ctx context.Context, f *flow.Flow) (*run.Run, error) {
	if err := f.Validate(nodes.OutPorts); err != nil {
		return nil, err
	}
	order, err := f.TopoSort()
	if err != nil {
		return nil, err
	}

	r := &run.Run{
		ID:         fmt.Sprintf("run_%d_%d", time.Now().UnixNano(), e.seq.Add(1)),
		FlowID:     f.ID,
		Status:     run.Running,
		NodeStatus: map[string]run.NodeStatus{},
	}
	for _, n := range f.Nodes {
		r.NodeStatus[n.ID] = run.Idle
	}
	e.emit("run.started", map[string]any{"runId": r.ID, "flowId": f.ID})

	env := &nodes.Env{
		Vars: f.Vars,
		Ctx:  ctx,
		Logf: slog.Info,
	}

	for _, id := range order {
		if ctx.Err() != nil {
			r.NodeStatus[id] = run.Skipd
			continue
		}
		// a failed upstream means this node is skipped
		if upstreamFailed(f, id, r) {
			r.NodeStatus[id] = run.Skipd
			e.emit("node.finished", map[string]any{"runId": r.ID, "nodeId": id, "status": run.Skipd, "ms": 0})
			continue
		}

		node := findNode(f, id)
		impl, err := nodes.Lookup(node.Type)
		if err != nil {
			r.NodeStatus[id] = run.Fail
			e.emit("node.finished", map[string]any{"runId": r.ID, "nodeId": id, "status": run.Fail, "ms": 0, "error": err.Error()})
			continue
		}

		r.NodeStatus[id] = run.Exec
		e.emit("node.started", map[string]any{"runId": r.ID, "nodeId": id})

		start := time.Now()
		runErr := impl.Run(env, node.Config)
		ms := time.Since(start).Milliseconds()

		if runErr != nil {
			r.Failure(f, id)
			e.emit("node.finished", map[string]any{"runId": r.ID, "nodeId": id, "status": run.Fail, "ms": ms, "error": runErr.Error()})
			continue
		}
		r.NodeStatus[id] = run.Ok
		e.emit("node.finished", map[string]any{"runId": r.ID, "nodeId": id, "status": run.Ok, "ms": ms})
	}

	r.Finalise(f)
	e.emit("run.finished", map[string]any{"runId": r.ID, "status": r.Status})
	return r, nil
}

func findNode(f *flow.Flow, id string) flow.Node {
	for _, n := range f.Nodes {
		if n.ID == id {
			return n
		}
	}
	return flow.Node{}
}

func upstreamFailed(f *flow.Flow, id string, r *run.Run) bool {
	for _, e := range f.Edges {
		if e.To.Node == id && r.NodeStatus[e.From.Node] == run.Fail {
			return true
		}
	}
	return false
}
