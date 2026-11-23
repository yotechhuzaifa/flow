package flow

import (
	"errors"
	"fmt"
)

type Pos struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Node struct {
	ID     string         `json:"id"`
	Type   string         `json:"type"`
	Pos    Pos            `json:"pos"`
	Config map[string]any `json:"config"`
}

type Port struct {
	Node string `json:"node"`
	Port string `json:"port"`
}

type Edge struct {
	ID   string `json:"id"`
	From Port   `json:"from"`
	To   Port   `json:"to"`
}

type Flow struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	UpdatedAt string            `json:"updatedAt"`
	Nodes     []Node            `json:"nodes"`
	Edges     []Edge            `json:"edges"`
	Vars      map[string]string `json:"vars"`
}

// PortNamer is implemented by whatever knows a node type's ports.
// The domain does not import the node registry: cmd wires it in.
type PortNamer func(nodeType string) []string

// Validate enforces CONTRACT invariants. ports may be nil to skip port checks.
func (f *Flow) Validate(ports PortNamer) error {
	seen := map[string]bool{}
	types := map[string]string{}
	for _, n := range f.Nodes {
		if n.ID == "" {
			return errors.New("node with empty id")
		}
		if seen[n.ID] {
			return fmt.Errorf("duplicate node id %q", n.ID)
		}
		seen[n.ID] = true
		types[n.ID] = n.Type
	}
	for _, e := range f.Edges {
		if !seen[e.From.Node] || !seen[e.To.Node] {
			return fmt.Errorf("edge %q references a missing node", e.ID)
		}
		if e.From.Node == e.To.Node {
			return fmt.Errorf("edge %q connects a node to itself", e.ID)
		}
		if ports != nil {
			declared := ports(types[e.From.Node])
			ok := false
			for _, p := range declared {
				if p == e.From.Port {
					ok = true
				}
			}
			if !ok {
				return fmt.Errorf("edge %q leaves an undeclared port %q of a %s node", e.ID, e.From.Port, types[e.From.Node])
			}
		}
	}
	_, err := f.TopoSort()
	return err
}
