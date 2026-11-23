# Research phase closed

Date: 2026-11-23. Target: Postman Flows (`spec/TARGET.md`). Docs version: v12.

The research loop is closed. Every gap in `spec/GAPS.md` is either covered by a trace or marked
`unobservable` with a reason. Implementation starts from `parity.csv`.

## Outcome

- 16 trace-backed features in `spec/features.jsonl`.
- 16 parity rows merged: `p01`–`p16`.
- Feature-to-row map: `spec/parity.map.json`.
- Test plan: `spec/TEST-PLAN.md`.
- No unmerged rows remain in `spec/parity.proposed.csv`.

## Covered

| Gap | Evidence |
|---|---|
| G-01 palette order/search | `traces/live-builder/palette.jsonl` |
| G-02 add-block placement | `traces/live-builder/add-block.jsonl` |
| G-03 self/duplicate/cycle rejection | `traces/live-builder/connections-clean.jsonl` |
| G-04 edge midpoint delete | `traces/live-builder/edge-delete-clean.jsonl` |
| G-05 Delay config clamp | `traces/live-builder/config-validation.jsonl` |
| G-06 run control / stop | `traces/live-builder/run-concurrency.jsonl`, `run-stop.jsonl` |
| G-07 unconfigured HTTP error | `traces/live-builder/failure-propagation.jsonl` |
| G-08 status timeline | `traces/live-builder/status-timeline.jsonl` |
| G-09 reload persistence | `traces/live-builder/reload.jsonl` |
| G-10 move/edge undo | `traces/live-builder/undo-redo.jsonl` |
| G-15 public read-only canvas | `traces/public-delay/` |

## Unobservable in this session

| Gap | Reason |
|---|---|
| G-03 source↔source, target↔target, empty-canvas drop | dragging a port opens the block palette, so the drop target cannot be isolated without a controlled harness |
| G-04 endpoint reconnect, delete a connected block | endpoint drag opens the palette; no stable handle to reconnect |
| G-05 HTTP required-field and invalid-URL validation | request editor hidden and the workspace reported offline |
| G-07 configured non-2xx and skipped descendants | the request editor was hidden; the downstream block on the failed branch disappeared when the run started |
| G-09 autosave cadence, crash recovery | no distinct save HTTP request appeared within five seconds; the app likely saves over its existing connection or debounces longer |
| G-10 add/connect/config undo, reload survival | the redo stack returned an intermediate move, not the final state; more clean runs are needed |
| G-11 modifier multi-select, group/ungroup | Ctrl+click selected only one node and no Group action appeared in the toolbar or context menu |
| G-12 limits: blocks, edges, run duration, >20 requests | needs a generated stress flow and multiple runs, not a single observation |
| G-13 variable precedence and interpolation | needs flow/environment/collection variables set up before the run |
| G-14 cross-flow copy/paste | workspace reported offline and only one flow was available |

These are not blockers for the clone. They are follow-up research items for a fresh authorized
session, recorded with their scenarios in `spec/GAPS.md`.

## Notes for the build loop

- `GOAL.md` and f02 now say select-and-place, matching the target.
- Keep the internal `Run` / `NodeStatus` names; the frontend maps the target’s
  `inactive` / `execution` / `success` states (CP-08).
- `dod.sh` counts every MUST row; do not reintroduce commas into the `parity.csv` feature column.
- `traces/` is gitignored. Evidence stays on the research machine; the committed rows cite it.
