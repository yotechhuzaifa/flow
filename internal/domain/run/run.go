package run

import "flow/internal/domain/flow"

type Status string

const (
	Running Status = "running"
	Done    Status = "done"
	Failed  Status = "failed"
)

type NodeStatus string

const (
	Idle  NodeStatus = "idle"
	Exec  NodeStatus = "running"
	Ok    NodeStatus = "ok"
	Fail  NodeStatus = "failed"
	Skipd NodeStatus = "skipped"
)

// Run is the record of one execution. Pure data: no IO, no framework.
type Run struct {
	ID         string                `json:"id"`
	FlowID     string                `json:"flowId"`
	Status     Status                `json:"status"`
	NodeStatus map[string]NodeStatus `json:"nodeStatus"`
}

// Failure marks a node failed and its descendants skipped.
func (r *Run) Failure(f *flow.Flow, nodeID string) {
	r.NodeStatus[nodeID] = Fail
	for _, d := range f.Descendants(nodeID) {
		if r.NodeStatus[d] == Idle || r.NodeStatus[d] == "" {
			r.NodeStatus[d] = Skipd
		}
	}
}

// Finalise sets the run status from its node statuses.
func (r *Run) Finalise(f *flow.Flow) {
	r.Status = Done
	for _, n := range f.Nodes {
		if r.NodeStatus[n.ID] == Fail {
			r.Status = Failed
			return
		}
	}
}
