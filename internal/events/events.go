// Package events is the single place where event names live.
// docs/CONTRACT.md documents them; this file enforces them.
package events

const (
	RunStarted   = "run.started"
	NodeStarted  = "node.started"
	NodeFinished = "node.finished"
	RunFinished  = "run.finished"
)

// All is the allow-list. An emitter that sends anything else is a bug.
var All = []string{RunStarted, NodeStarted, NodeFinished, RunFinished}

func Valid(name string) bool {
	for _, n := range All {
		if n == name {
			return true
		}
	}
	return false
}
