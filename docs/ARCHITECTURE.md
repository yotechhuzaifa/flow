# ARCHITECTURE

## Shape

```
cmd/app      → internal/wails      → internal/transport (service) → internal/app (wiring)
cmd/server   → internal/transport/httpapi
                        ↓
                   internal/app  (composition root: the only file that knows every layer)
                        ↓
   internal/exec ─ internal/nodes ─ internal/domain/{flow,run} ─ internal/sqlite ─ internal/platform
                        ↑
   internal/events (names)      internal/webui (embedded frontend)
```

Dependency rule: **arrows point inward**. `domain` imports nothing but the standard library.
`transport` and `sqlite` may import `domain`. Nothing imports `cmd`. Only `internal/app` imports
more than two siblings — that is its single job.

## Layers

| Layer | Package | Owns | Never does |
|---|---|---|---|
| domain | `internal/domain/flow`, `internal/domain/run` | Flow, Node, Edge, Run, TopoSort, Validate, Descendants | IO, framework, SQL, HTTP |
| nodes | `internal/nodes` | one file per node type, a registry, `Env` | knows about Wails, storage, or React |
| exec | `internal/exec` | run orchestration, status transitions, event emission | persist, render, or parse config JSON |
| store | `internal/sqlite` | schema migrations, Save/Get/List/Delete/RecordRun | business rules beyond calling Validate |
| transport | `internal/wails`, `internal/transport/httpapi` | Wails app, service methods, SSE/HTTP | business logic |
| webui | `internal/webui` | `go:embed dist` and `AssetFileServerFS` | anything else |
| app | `internal/app` | wiring, node registration, emitter injection | features |
| platform | `internal/platform` | user data dir, GOOS | anything else |

## Frontend

```
frontend/src/
  main.tsx                mount, wails runtime import
  App.tsx                 shell: palette + canvas + run button
  features/canvas/        Canvas.tsx (@xyflow/react), Palette.tsx
  features/nodes/         one component per node type; registry in nodes/index.ts
  features/runs/          run status overlay, log drawer (added when f08 lands)
  shared/lib/             types.ts (mirrors CONTRACT), useFlow.ts (state + autosave)
  shared/api/client.ts    the only module that imports the generated bindings
  shared/ui/              small presentational components
```

Rule: a feature folder may import `shared`, never another feature. Cross-feature needs go to
`shared` or into the store.

## Adding a node type (the growth path)

1. `internal/nodes/<type>.go` — implement `Type/OutPorts/Run`
2. register it in `internal/app/app.go` (`nodes.Register(...)`)
3. `frontend/src/features/nodes/<Type>Node.tsx` + add to `nodes/index.ts` and `OUT_PORTS`
4. one parity row + one Go test + one e2e test

Nothing else changes. If a new node type forces edits in `exec`, the node interface is wrong —
fix the interface, not the executor.

## Adding a feature on top of the clone

New capabilities sit in new packages, not inside the cloned ones:

```
internal/feature/<name>/     service + tests, imports domain only
internal/transport/…         one new method on the existing service
frontend/src/features/<name>/  its own folder
schema/<feature>.json        its own persisted shape
```

The clone is `features/canvas` + `domain/flow`. Everything after it is additive. That is the whole
reason for the layer split: when you add your differentiating feature, you do not touch the
cloned code.

## Cross-platform

One codebase, three targets, no `#ifdef` in feature code.

| Concern | Where it lives |
|---|---|
| Data dir | `internal/platform.DataDir` (macOS `~/Library/Application Support`, Linux `~/.config`, Windows `%AppData%`) |
| Native window chrome | `internal/wails/app.go` — `application.MacWindow` only affects darwin |
| WebView runtime | Windows WebView2, macOS WKWebView, Linux WebKitGTK 6.0 (GTK4) — all provided by Wails |
| Packaging | `build/{windows,darwin,linux}/Taskfile.yml` — to be added with the Wails v3 template; not vendored yet |
| CI | `.github/workflows/build.yml` — planned, not present yet; run `scripts/check.sh` locally |

Never branch on `runtime.GOOS` outside `internal/platform` and `internal/wails`. If you need a
platform behaviour, add a small function to `platform` with one file per OS.

## Server mode

`-tags server` builds the same domain with no webview: `cmd/server` serves `frontend/dist` plus
`/api/*` and `/events` (SSE). This is what CI runs the parity suite against, and it is the seam
where a hosted/team version would attach later — no rewrite.
