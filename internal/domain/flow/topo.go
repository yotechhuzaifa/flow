package flow

import "fmt"

// TopoSort returns node ids in execution order. Cycles are an error.
func (f *Flow) TopoSort() ([]string, error) {
	indeg := map[string]int{}
	next := map[string][]string{}
	for _, n := range f.Nodes {
		indeg[n.ID] = 0
	}
	for _, e := range f.Edges {
		next[e.From.Node] = append(next[e.From.Node], e.To.Node)
		indeg[e.To.Node]++
	}

	var queue []string
	for _, n := range f.Nodes {
		if indeg[n.ID] == 0 {
			queue = append(queue, n.ID)
		}
	}

	var order []string
	for len(queue) > 0 {
		id := queue[0]
		queue = queue[1:]
		order = append(order, id)
		for _, d := range next[id] {
			indeg[d]--
			if indeg[d] == 0 {
				queue = append(queue, d)
			}
		}
	}

	if len(order) != len(f.Nodes) {
		return nil, fmt.Errorf("cycle detected: sorted %d of %d nodes", len(order), len(f.Nodes))
	}
	return order, nil
}

// Descendants returns every node reachable downstream of id.
func (f *Flow) Descendants(id string) []string {
	next := map[string][]string{}
	for _, e := range f.Edges {
		next[e.From.Node] = append(next[e.From.Node], e.To.Node)
	}
	seen := map[string]bool{}
	var walk func(string)
	walk = func(cur string) {
		for _, d := range next[cur] {
			if seen[d] {
				continue
			}
			seen[d] = true
			walk(d)
		}
	}
	walk(id)
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	return out
}
