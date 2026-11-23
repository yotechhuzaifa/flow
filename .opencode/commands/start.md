---
description: Start here. Boots the project, then works the parity queue until dod.sh is green.
agent: build
---

You are starting a fresh session on this repo. Do this in order. Do not skip steps.

1. Sanity check. Run `./scripts/check.sh`. Read `GOAL.md`, `.agents/STATE.md`, `parity.csv`,
   `docs/CONTRACT.md`. Do not change any of them.

2. Research, if `spec/features.jsonl` does not exist or is empty:
   read the `grabbing-flows` skill, then do exactly one scenario and one extraction pass.
   Write `spec/features.jsonl` with evidence paths. Then continue to step 3 without asking.

3. Confirm the scaffold, if `./scripts/check.sh` is red. It is already laid out; do not rewrite it:
   - `cmd/app` desktop entry, `cmd/server` headless entry (`-tags server`)
   - `internal/domain/flow`, `internal/domain/run` — types, `Validate`, `TopoSort`
   - `internal/nodes` (one file per type), `internal/exec` (orchestration), `internal/events`
   - `internal/sqlite` (migrations in `internal/sqlite/migrations/`), `internal/platform`
   - `internal/transport/httpapi` (HTTP + SSE), `internal/wails` (service `wails.Flows`),
     `internal/webui` (`go:embed dist`)
   - `frontend`: Vite + React + TS + `@xyflow/react`; aliases in `frontend/vite.config.ts`;
     generated bindings resolve through `@bindings/*` in `frontend/bindings`
   - regenerate bindings with `wails3 generate bindings -ts -d frontend/bindings`
   Get `./scripts/check.sh` green before touching any parity row. The desktop shell still needs the
   Wails v3 CLI and its `build/` assets (see `docs/platform-setup.md`); the headless path
   (`cmd/server` + Vite + Playwright) is what the loop uses.

4. Work the queue. Loop:
   - take the highest `MUST` row with status `todo`
   - write the plan in `.agents/STATE.md`
   - implement the smallest change that makes it pass
   - `./scripts/check.sh`, then mark the row `pass` and commit
   - repeat until `./scripts/dod.sh` exits 0

Rules: one row at a time, no refactors, never weaken a test, never edit frozen files.
If the same check fails three times, write the question to `BLOCKED.md` and stop.
