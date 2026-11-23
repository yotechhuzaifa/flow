# Scaling after the clone

The clone is a floor, not a ceiling. This is the plan for what comes after.

## Phase 0 — clone (now)
Canvas, nodes, edges, run, live status, SQLite persistence. 11 MUST rows in `parity.csv`.
Success = `dod.sh` green.

## Phase 1 — make the clone honest
- Node config editors with schema validation
- Retry / timeout per node, error taxonomy (`network | timeout | validation | internal`)
- Run history in SQLite (`runs` table already exists) + a runs drawer
- Import/export flow JSON

## Phase 2 — differentiate (this is why you cloned)
Pick one. Not three.

| Direction | Where it plugs in |
|---|---|
| Local-first AI node ("transform with a prompt") | new `internal/nodes/ai.go`, calls an OpenAI-compatible endpoint from settings |
| Batch/collection runner over a CSV | new `internal/exec/batch.go` — the executor already takes a Flow, add a dataset |
| Assertions as first-class nodes | new node type + `run` record extension |
| Offline environment/workspace switching | `sqlite` gains a workspace column; UI gets a switcher |

The layer rules make each of these additive. If a feature requires editing `domain/flow`, it is
probably a new entity — add a package, do not widen the core.

## Phase 3 — team / hosted
`-tags server` is already the seam. That path needs:
- auth (out of scope today, but `transport/httpapi` is the only place to add middleware)
- SQLite → Postgres behind the same `Store` interface (extract an interface in `internal/app`)
- SSE → websocket if you need bidirectional control
- object storage for captured runs

## Guardrails that keep scaling cheap
1. `domain` stays dependency-free. It is the part that survives every rewrite.
2. One node type = one file + one registration + one component + one row. No exception.
3. Every feature folder owns its own state. `shared` grows only when two features need the same thing.
4. `internal/app` is the only place that knows the whole graph. When wiring gets painful, that is the
   signal to add a package, not to add globals.
5. `check.sh` under 60s, forever. The moment it is slow, iteration speed dies and quality follows.
6. Every new platform behaviour goes through `internal/platform`, one file per OS.
