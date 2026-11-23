# CONTRACT v1

Frozen. Everything codes against this file.

```jsonc
Flow  { id: string, name: string, updatedAt: string, nodes: Node[], edges: Edge[], vars: Record<string,string> }
Node  { id: string, type: "send" | "delay" | "log" | "if", pos: { x: number, y: number }, config: object }
Edge  { id: string, from: { node: string, port: string }, to: { node: string, port: string } }
Run   { id: string, flowId: string, status: "running" | "done" | "failed", nodeStatus: Record<string, NodeStatus> }
```

Ports per type: `send | delay | log → out`, `if → true, false`. Two ports are implicit:
every node accepts `in`; edges always start at a declared output port.

`NodeStatus`: `idle` | `running` | `ok` | `failed` | `skipped`

## Events
| Event | Payload |
|---|---|
| `run.started` | `{ runId, flowId }` |
| `node.started` | `{ runId, nodeId }` |
| `node.finished` | `{ runId, nodeId, status, ms, error? }` |
| `run.finished` | `{ runId, status }` |

Names are enforced by `internal/events`. Emitting anything else is a bug.

## Invariants
1. An edge references existing node ids and a declared output port.
2. Self-edges are rejected. Cycles are rejected at save time, not run time.
3. A failed node marks its descendants `skipped`. The run becomes `failed`.
4. The frontend never mutates a `Run`; it renders events.
5. `{{name}}` resolves from `flow.vars` at run time, not at save time.

## Persistence
SQLite, one file. `flows(id, name, updated_at, data)` stores the whole Flow as JSON in `data`;
`runs(id, flow_id, status, started_at, finished_at)`. Migrations live in
`internal/sqlite/migrations/` and are applied on open.
