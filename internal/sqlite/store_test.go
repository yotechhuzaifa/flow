package sqlite

import (
	"path/filepath"
	"testing"

	"flow/internal/domain/flow"
)

func TestSaveGetRoundTrip(t *testing.T) {
	s, err := Open(filepath.Join(t.TempDir(), "app.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	in := &flow.Flow{
		ID:   "f1",
		Name: "test",
		Nodes: []flow.Node{
			{ID: "a", Type: "send", Pos: flow.Pos{X: 10, Y: 20}, Config: map[string]any{"url": "https://x"}},
			{ID: "b", Type: "log"},
		},
		Edges: []flow.Edge{{ID: "e", From: flow.Port{Node: "a", Port: "out"}, To: flow.Port{Node: "b", Port: "in"}}},
	}
	if err := s.Save(in); err != nil {
		t.Fatal(err)
	}
	out, err := s.Get("f1")
	if err != nil {
		t.Fatal(err)
	}
	if len(out.Nodes) != 2 || out.Nodes[0].Pos.X != 10 {
		t.Fatalf("round trip lost data: %+v", out.Nodes)
	}
	if err := s.Delete("f1"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Get("f1"); err == nil {
		t.Fatal("expected not found after delete")
	}
}
