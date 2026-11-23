package exec

import (
	"context"
	"errors"
	"testing"

	"flow/internal/domain/flow"
	"flow/internal/domain/run"
	"flow/internal/nodes"
)

func init() {
	nodes.Register(nodes.Log{}, nodes.Send{}, nodes.Delay{}, nodes.If{})
}

func chain(t *testing.T) *flow.Flow {
	t.Helper()
	return &flow.Flow{
		ID: "f1",
		Nodes: []flow.Node{
			{ID: "a", Type: "log", Config: map[string]any{"message": "hi {{who}}"}},
			{ID: "b", Type: "log"},
			{ID: "c", Type: "log"},
		},
		Edges: []flow.Edge{
			{ID: "e1", From: flow.Port{Node: "a", Port: "out"}, To: flow.Port{Node: "b", Port: "in"}},
			{ID: "e2", From: flow.Port{Node: "b", Port: "out"}, To: flow.Port{Node: "c", Port: "in"}},
		},
		Vars: map[string]string{"who": "there"},
	}
}

func TestRunOrderAndStatuses(t *testing.T) {
	f := chain(t)
	var order []string
	e := New(func(ev string, p any) {
		if ev == "node.started" {
			order = append(order, p.(map[string]any)["nodeId"].(string))
		}
	})
	r, err := e.Start(context.Background(), f)
	if err != nil {
		t.Fatal(err)
	}
	if r.Status != run.Done {
		t.Fatalf("status %s", r.Status)
	}
	if len(order) != 3 || order[0] != "a" || order[2] != "c" {
		t.Fatalf("order %v", order)
	}
}

func TestFailureSkipsDescendants(t *testing.T) {
	f := chain(t)
	f.Nodes[1].Type = "send"
	f.Nodes[1].Config = map[string]any{"url": ""} // fails validation of config
	r, err := New(nil).Start(context.Background(), f)
	if err != nil {
		t.Fatal(err)
	}
	if r.Status != run.Failed {
		t.Fatalf("expected failed run, got %s", r.Status)
	}
	if r.NodeStatus["c"] != run.Skipd {
		t.Fatalf("downstream should be skipped, got %s", r.NodeStatus["c"])
	}
}

func TestUnknownNodeType(t *testing.T) {
	f := &flow.Flow{ID: "x", Nodes: []flow.Node{{ID: "a", Type: "nope"}}}
	r, err := New(nil).Start(context.Background(), f)
	if err != nil {
		t.Fatal(err)
	}
	if r.NodeStatus["a"] != run.Fail {
		t.Fatalf("expected fail, got %s", r.NodeStatus["a"])
	}
}

var _ = errors.New
