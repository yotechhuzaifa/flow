# Contract proposals

Status: **partly approved**. `docs/CONTRACT.md` is frozen. Each item needs a human `[contract]`
commit before code changes. Evidence and the reference projects are in `docs/RESEARCH.md`.

## Decisions (2026-11-23)

- **CP-07 accepted.** Add-block is select-then-place. `GOAL.md` and `parity.csv` f02 now say
  “select a node from the palette and place it on the canvas”. The 11 researched rows `p01`–`p11`
  are merged into `parity.csv`.
- **CP-08 accepted as the mapping option.** Keep the internal `Run` / `NodeStatus` names. The
  frontend maps the target’s `inactive` / `execution` / `success` states to the clone’s
  `idle` / `running` / `ok`. No `docs/CONTRACT.md` change.
- **CP-01 through CP-06 remain unapproved.** They still need their own `[contract]` decision.

Do not implement these from a research session. Write the proposal here, link the feature id
and trace, and stop.

## CP-01 — version the flow document

- Problem: `Flow` has no `schemaVersion`. SQLite has migrations; flow JSON does not.
- Evidence: Windmill serialises flows in the versioned OpenFlow format; n8n and Node-RED
  version their flow documents.
- Proposal: `Flow.schemaVersion: number`, `internal/domain/flow.Migrate(raw) []byte`, and a
  migration test that opens an old fixture.
- Needed by: import/export, third-party nodes, any persisted format change.

## CP-02 — typed ports

- Problem: `docs/CONTRACT.md` declares ports but not their types. Any output can reach any
  input, so invalid edges are only caught by the executor.
- Evidence: Langflow and ComfyUI declare typed inputs/outputs; React Flow can reject a
  connection before it is added with `isValidConnection`.
- Proposal: `Port { id, type }` in the node manifest; `Validate` checks source type against
  target type; the canvas uses the same rule through `isValidConnection`.
- Needed by: `if` branches, HTTP/JSON nodes, any node that produces structured data.

## CP-03 — one node manifest, generated surfaces

- Problem: a node is repeated in `internal/nodes/<type>.go`, `frontend/src/features/nodes/`,
  `OUT_PORTS`, a parity row, and a test. Drift is the top maintenance risk.
- Evidence: n8n (`*.node.ts` + codex JSON), Node-RED (JS + HTML), Langflow (typed component),
  ComfyUI (`INPUT_TYPES`) all keep the node definition in one place.
- Proposal: `spec/nodes.json` with type, label, ports, config schema, and executor key.
  Generate the Go registry stub and the TS palette/port map from it. Keep the hand-written
  `Run` function; generate only the metadata surfaces.
- Needed by: third-party nodes, config editors, schema validation.

## CP-04 — persist per-node run data

- Problem: `runs` stores status and timestamps; `nodeStatus` lives in memory. A failed run
  cannot be inspected after restart.
- Evidence: n8n keeps execution history with per-node data; Node-RED exposes debug and status.
- Proposal: `run_nodes(run_id, node_id, status, started_at, finished_at, error, output)` or a
  JSON blob on `runs`. Keep the event names unchanged.
- Needed by: run history drawer, retry, support/debugging.

## CP-05 — cancellation, retry, timeout in the executor

- Problem: `internal/exec` has no context, cancellation, retry, or timeout. A slow `send` node
  cannot be stopped, and a transient failure cannot be retried.
- Evidence: Windmill treats scheduling, errors, and early return as engine concerns; the node
  stays simple. ComfyUI has explicit run states.
- Proposal: `Run(ctx context.Context, ...)`, `run.cancelled` event, per-node `timeout` and
  `retry` in config, all enforced by `exec`. Nodes stay unaware.
- Needed by: f11/f12 hardening, SCALE Phase 1.

## CP-07 — add-block interaction is click-to-place

- Problem: `GOAL.md` and parity row f02 say “drag a node from the palette onto the canvas”.
  The live target does not support that: palette items are not draggable. The user selects an
  item, the palette closes, and the next canvas click places the block.
- Evidence: `traces/live-builder/add-block.jsonl`.
- Proposal: change f02 to “select a block from the palette and place it on the canvas”, or
  keep drag as a deliberate divergence and add a separate parity row for it.
- Needs: a human `[contract]` decision, because `GOAL.md` and `parity.csv` are frozen.

## CP-08 — execution-state names

- Problem: the live target exposes `inactive`, `execution`, and `success` on block headers.
  The clone’s contract uses `idle`, `running`, `ok`, `failed`, `skipped`.
- Evidence: `traces/live-builder/status-timeline.jsonl`.
- Proposal: keep the internal event names and map them in the UI, or rename the contract values
  in a migration. Add a parity test for the visible state, not the internal attribute.
- Needs: no contract change if the UI maps; a `[contract]` change if event names are renamed.

## CP-06 — node type version

- Problem: changing a node's config shape breaks saved flows with no upgrade path.
- Evidence: n8n recommends node versioning and keeps version directories.
- Proposal: `Node.typeVersion: number`; the registry resolves `(type, version)`; a migration
  maps old config to new. Default `1`.
- Needed by: any config-shape change after the clone.
