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
7. `docs/RESEARCH.md` — the research loop, Layer 0 sources, and OSS references.

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
spec/                   features.jsonl, target, proposals, test plan, audit
traces/                 playwright evidence, one folder per scenario
docs/RESEARCH.md        the six-layer research pipeline
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
- Load the `grabbing-flows` skill before any observation session; load `spec-compiler` to hand
  research to the build loop.
- Layer 0 is free intel first (`docs/RESEARCH.md`): docs, changelog, community threads, public
  API docs. Use `firecrawl` for pages and `playwright` only for the gaps the docs leave open.
- Extract behaviour, never implementation. No copied CSS, assets, icons, class names, microcopy.
- The target's shipped JavaScript and sourcemaps are off-limits. Open-source reference code may
  be read, but never copied; record repo + commit + path + license in `spec/INTEL.md`.
- Every feature claim cites a trace path. No citation, no claim. Doc-only claims stay in
  `spec/INTEL.md`.
- Run `node scripts/spec.mjs check` after extraction and `audit` before marking a row pass.
- Never edit `parity.csv` or `docs/CONTRACT.md` from research. Write
  `spec/parity.proposed.csv` or `spec/CONTRACT-PROPOSAL.md`; a human merges with a
  `[contract]` commit.
- Never bypass authentication or touch another user's data.

## Stop conditions
- `./scripts/dod.sh` exits 0 → write `DONE.md`, stop.
- Same check failing three times → write the question to `BLOCKED.md`, stop.
