---
name: wails-v3
description: Rules for Wails v3 code in this repo. Load before touching cmd/ or wiring events.
autoinvoke: true
---

# Wails v3

v3.0.0 went stable on 2026-10-26. Anything you remember about v3 is probably alpha. Verify locally.

## Generate the truth, once, at the start of any Wails task
```bash
go doc -all github.com/wailsapp/wails/v3/pkg/application > docs/WAILS_V3.md
wails3 generate bindings -ts -d frontend/bindings
```
`docs/WAILS_V3.md` and `frontend/bindings/**` are the only sources you may code against.

## Rules
- Services are the unit. The frontend calls generated methods; never fetch a blank URL.
- Never hand-edit `frontend/bindings/**`. Regenerate instead.
- Never change the Wails version in `go.mod`.
- Events: emit only the names in `docs/CONTRACT.md`.
- The frontend must `import "@wailsio/runtime"` for side effects, or events never arrive.
- If a signature is not in the generated docs, stop and write the question to `BLOCKED.md`.
  Do not invent one.

## After any Go change
```bash
go build ./... && go vet ./...
```
