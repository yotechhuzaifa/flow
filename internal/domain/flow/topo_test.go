package flow

import "testing"

func TestTopoSortLinear(t *testing.T) {
	f := &Flow{
		Nodes: []Node{{ID: "a", Type: "send"}, {ID: "b", Type: "log"}},
		Edges: []Edge{{ID: "e", From: Port{"a", "out"}, To: Port{"b", "in"}}},
	}
	got, err := f.TopoSort()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Fatalf("unexpected order %v", got)
	}
}

func TestTopoSortCycle(t *testing.T) {
	f := &Flow{
		Nodes: []Node{{ID: "a", Type: "send"}, {ID: "b", Type: "send"}},
		Edges: []Edge{
			{ID: "e1", From: Port{"a", "out"}, To: Port{"b", "in"}},
			{ID: "e2", From: Port{"b", "out"}, To: Port{"a", "in"}},
		},
	}
	if _, err := f.TopoSort(); err == nil {
		t.Fatal("expected a cycle error")
	}
}

func TestValidateRejectsUndeclaredPort(t *testing.T) {
	f := &Flow{
		Nodes: []Node{{ID: "a", Type: "log"}, {ID: "b", Type: "log"}},
		Edges: []Edge{{ID: "e", From: Port{"a", "nope"}, To: Port{"b", "in"}}},
	}
	ports := func(string) []string { return []string{"out"} }
	if err := f.Validate(ports); err == nil {
		t.Fatal("expected an undeclared port error")
	}
}

func TestValidateNilPortsSkipsPortChecks(t *testing.T) {
	f := &Flow{
		Nodes: []Node{{ID: "a", Type: "log"}, {ID: "b", Type: "log"}},
		Edges: []Edge{{ID: "e", From: Port{"a", "anything"}, To: Port{"b", "in"}}},
	}
	if err := f.Validate(nil); err != nil {
		t.Fatal(err)
	}
}

func TestDescendants(t *testing.T) {
	f := &Flow{
		Nodes: []Node{{ID: "a"}, {ID: "b"}, {ID: "c"}, {ID: "d"}},
		Edges: []Edge{
			{ID: "1", From: Port{"a", "out"}, To: Port{"b", "in"}},
			{ID: "2", From: Port{"b", "out"}, To: Port{"c", "in"}},
		},
	}
	got := f.Descendants("a")
	if len(got) != 2 {
		t.Fatalf("want b,c got %v", got)
	}
}
