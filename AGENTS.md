# AGENTS.md

## Objective
A cross-platform desktop app (Windows, macOS, Linux) that reimplements a visual flow builder's
**behaviour**: palette, canvas, nodes, edges, run with live per-node status, SQLite persistence.
Frontend first; backend exists only to serve it. Later phases add features on top — see `docs/SCALE.md`.

## Read order
1. `GOAL.md`
2. `.agents/STATE.md`
3. `parity.csv` — work `todo` + `MUST` rows only, top to bottom
4. `docs/CONTRACT.md` — frozen schema and event names
5. `docs/ARCHITECTURE.md` — where code goes. Obey the dependency direction.
6. `docs/SCALE.md` — the phase plan after the clone. Do not start phase 2 work now.

## Frozen (needs a human `[contract]` commit)
`docs/CONTRACT.md` · `parity.csv` · `scripts/check.sh` · `scripts/dod.sh` · `GOAL.md`

## Layout and the rule that matters
```
cmd/app                 desktop entry (wails)
cmd/server              headless entry (-tags server)
internal/domain/flow    Flow, Node, Edge, TopoSort, Validate   ← stdlib only
internal/domain/run     Run, statuses
internal/nodes/         one file per node type + registry
internal/exec/          run orchestration + events
internal/sqlite/        persistence, migrations
internal/wails/         desktop transport
internal/transport/httpapi  server transport
internal/webui/         go:embed frontend dist
internal/app/           composition root
internal/platform/      data dir, GOOS
internal/events/        event names (allow-list)
frontend/src/features/  one folder per feature
frontend/src/shared/    types, api client, hooks
```

Dependency arrows point inward. `domain` imports nothing but stdlib. Only `internal/app` may import
several layers. Never import `cmd` from anywhere.

## Rules
- One parity row per iteration. Small change, verify, commit.
- New node type = `internal/nodes/<type>.go` + one line in `internal/app` + one component in
  `frontend/src/features/nodes/` + one parity row + one Go test. No edits to `internal/exec`.
- Canvas is `@xyflow/react`. Never implement drag, zoom, or edge routing.
- Every interactive element gets a `data-testid`. The parity suite depends on it.
- Executor must run headless: `go test -tags server ./...` passes with no UI and no database file
  in the repo.
- SQLite only, through `internal/sqlite`. One file, `app.db`, path from `internal/platform`.
- No `runtime.GOOS` outside `internal/platform` and `internal/wails`.
- No new dependency without recording the reason in `.agents/STATE.md` first.
- Keep `bash scripts/check.sh` under 60 seconds.

## Research rules
- Load the `grabbing-flows` skill before any observation session.
- Extract behaviour, never implementation. No copied CSS, assets, icons, class names, microcopy.
- Every spec claim cites a trace path. No citation, no claim.
- Never bypass authentication or touch another user's data.

## Stop conditions
- `./scripts/dod.sh` exits 0 → write `DONE.md`, stop.
- Same check failing three times → write the question to `BLOCKED.md`, stop.
