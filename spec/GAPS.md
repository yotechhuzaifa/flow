# GAPS — Postman Flows

Behaviours the docs describe but do not pin down. These are the observation scenarios; nothing
else. A gap is not a feature request and not a parity row until it has a trace and a
`spec/features.jsonl` line.

Status key:

- `open` — not observed yet.
- `needs live builder` — the editable canvas requires a Postman account; do not bypass auth.
- `public example` — partially observable on a read-only Flow Snippet without a login.
- `unobservable, reason: …` — cannot be observed safely or repeatably.

The first three are the core clone behaviours. All three need the live builder. Public Flow
Snippets can confirm layout and vocabulary, but not editing, connection rejection, or run/status.

## Queue

| # | Gap | Priority | Status |
|---|---|---|---|
| G-01 | Palette order, search, filter, category expansion | must | covered by `traces/live-builder/palette.jsonl` |
| G-02 | Add-block placement feedback and cancel | must | covered by `traces/live-builder/add-block.jsonl` |
| G-03 | Connection rejection rules and error feedback | must | partial: self/duplicate/cycle covered by `connections-clean.jsonl` |
| G-04 | Edge select, delete, reconnect | must | partial: midpoint delete observed; `traces/live-builder/edge-delete-clean.jsonl` |
| G-05 | Config validation, required fields, save timing | must | partial: Delay clamps invalid input; `traces/live-builder/config-validation.jsonl` |
| G-06 | Run concurrency, stop, re-run, edit during run | should | covered by `run-concurrency.jsonl` and `run-stop.jsonl` |
| G-07 | Failure propagation and skipped descendants | must | partial: unconfigured HTTP error observed; `failure-propagation.jsonl` |
| G-08 | Status vocabulary and transition timing | must | covered by `traces/live-builder/status-timeline.jsonl` |
| G-09 | Autosave cadence, crash recovery, reload | must | partial: reload covered by `traces/live-builder/reload.jsonl` |
| G-10 | Undo/redo scope | should | partial: move undo covered; `traces/live-builder/undo-redo.jsonl` |
| G-11 | Multi-select, box select, group/ungroup | could | partial: box select observed; `traces/live-builder/multiselect.jsonl` |
| G-12 | Limits: blocks, edges, run duration, >20 requests | should | needs live builder |
| G-13 | Variable precedence and interpolation | should | needs live builder |
| G-14 | Copy/paste across flows: ids, positions, config | could | needs live builder |
| G-15 | Public example read-only layout | info | covered by `traces/public-delay/` |

## G-01 — Palette order, search, filter, category expansion

- Doc evidence: categories are listed at
  [blocks overview](https://learning.postman.com/flows/reference/blocks/overview); the add-block
  entry points are at
  [create blocks](https://learning.postman.com/flows/build-flows/create/create-blocks).
- Not pinned: visual order, search, filter, favourites, collapsed categories, empty state, and
  whether the list is grouped or flat.
- Scenario: open a new flow, open the block list, record the order; type a query; clear it;
  expand/collapse each category; record the empty state.
- Capture: a11y snapshot of the list, DOM text order, search result counts.
- Status: covered by `traces/live-builder/palette.jsonl`. Observed categories in order:
  Primitives, Connectors, AI, Logic, Looping, Visualize, Outputs, Data, Other. Search `delay`
  returned Logic → Delay plus Suggested → String “delay”; `zzzz` still suggested a String block;
  Escape closed the palette.

## G-02 — Drag/drop feedback and cancel

- Doc evidence: blocks can be added by drag from a port or by the block list.
  [create blocks](https://learning.postman.com/flows/build-flows/create/create-blocks)
- Not pinned: ghost preview, drop-target highlight, snap/alignment guides, invalid-drop
  feedback, Esc/cancel, and whether the block list closes on drop.
- Scenario: drag a block over the canvas, over an existing block, over a port, then drop; repeat
  and press Esc; repeat and drop outside the canvas.
- Capture: screenshots only if a11y cannot show state; DOM classes are not evidence.
- Status: covered by `traces/live-builder/add-block.jsonl`. The target does **not** drag from
  the palette: palette items have no `draggable` attribute. Selecting an item closes the palette,
  and the next canvas click places the block at that point. This contradicts `GOAL.md` and
  parity row f02, which say “drag from the palette”. See `CONTRACT-PROPOSAL.md` CP-07.

## G-03 — Connection rejection rules and error feedback

- Doc evidence: connections are unidirectional; output ports can fan out.
  [create connections](https://learning.postman.com/flows/build-flows/create/create-connections)
- Not pinned: self-edge, cycle, duplicate edge, incompatible port, dropping on empty canvas,
  reconnecting an edge, and whether an error is shown or the edge silently disappears.
- Scenario: attempt each invalid connection; record whether an edge is created and what the UI
  shows. Then create a valid two-node chain and try to close the cycle.
- Capture: `network` if the app posts edges; `dom`/`a11y` for the rejection message.
- Status: partial. `traces/live-builder/connections-clean.jsonl` confirms self-edge, duplicate,
  and cycle are all rejected silently: the edge count is unchanged and no error toast appears.
  Source-to-source, target-to-target, and empty-canvas drops still need a clean rerun without
  the palette overlay.

## G-04 — Edge select, delete, reconnect

- Doc evidence: connections are created by dragging ports.
  [create connections](https://learning.postman.com/flows/build-flows/create/create-connections)
- Not pinned: edge hit area, selection style, Delete/Backspace behaviour, reconnect by dragging
  an endpoint, and whether deleting a block also removes its edges with confirmation.
- Scenario: select an edge, delete it; drag one endpoint to another port; delete a connected
  block.
- Capture: `dom`/`a11y` before and after; `network` if edges are persisted.
- Status: partial. `traces/live-builder/edge-delete-clean.jsonl` shows that clicking the control
  at a connection’s midpoint removes the edge immediately, with no confirmation dialog. Edge
  selection styling, endpoint reconnect, and deleting a connected block are not covered.

## G-05 — Config validation, required fields, save timing

- Doc evidence: HTTP Request selects a collection request and exposes variables; Parse Response
  has enum options.
  [HTTP Request](https://learning.postman.com/flows/reference/blocks/http-request)
- Not pinned: required-field markers, inline error text, whether an invalid block can run,
  debounce/save timing, and whether closing the panel discards an edit.
- Scenario: clear a required field, enter an invalid URL, close and reopen the block, reload.
- Capture: `dom`/`a11y` for error text; `network` for save calls.
- Status: partial. `traces/live-builder/config-validation.jsonl` shows the Delay field rejecting
  negative and empty input and falling back to 0, with no error state. The HTTP request editor
  was hidden and the workspace reported offline, so required-field and invalid-URL validation
  on HTTP Request were not observed.

## G-06 — Run concurrency, stop, re-run, edit during run

- Doc evidence: Run action and pause-on-block exist.
  [blocks overview](https://learning.postman.com/flows/reference/blocks/overview)
- Not pinned: whether Run is disabled while running, how to stop/cancel, whether editing during
  a run is allowed, and whether a second run queues or replaces the first.
- Scenario: run a slow flow; press Run again; try to stop; edit a block mid-run.
- Capture: `steps` timestamps, button disabled states, status changes.
- Status: covered by `traces/live-builder/run-concurrency.jsonl` and `run-stop.jsonl`. During a
  run the Run action becomes Stop, so there is no second Run to click. Pressing Stop returns the
  action to Run; the in-flight block switches to `success` and downstream blocks stay
  `inactive`. Editing during a run was not tested.

## G-07 — Failure propagation and skipped descendants

- Doc evidence: HTTP Request has Success and Fail outputs; non-2xx goes to Fail.
  [HTTP Request](https://learning.postman.com/flows/reference/blocks/http-request)
- Not pinned: whether a failed block fails the run, whether unconnected descendants are
  skipped, what status they show, and whether the Fail port is optional.
- Scenario: force a non-2xx response; watch the block and every downstream block; repeat with
  Fail unconnected.
- Capture: per-block status over time; run summary.
- Status: partial. `traces/live-builder/failure-propagation.jsonl` shows an HTTP Request block
  with no request selected entering `error` while the rest of the flow still ran. The downstream
  block connected to its Success port remained `inactive` before the run, but it disappeared
  when the run started, so the skip path is not clean enough for a feature line. A configured
  request returning a non-2xx is still needed.

## G-08 — Status vocabulary and transition timing

- Doc evidence: run logs show success/fail and a paused state.
  [troubleshoot](https://learning.postman.com/flows/build-flows/troubleshoot/troubleshoot)
- Not pinned: idle, queued, running, success, fail, skipped, cancelled; whether transitions are
  animated; how long a completed state persists.
- Scenario: run a three-block chain and sample each block's status every 100 ms until the run
  ends.
- Capture: status timeline, not a single screenshot.
- Status: covered by `traces/live-builder/status-timeline.jsonl`. Observed values are
  `inactive`, `execution`, and `success`. Start reached success at about 150 ms; the HTTP block
  entered `execution` and later switched to `success`; Evaluate and Display switched to success
  when the HTTP block completed. This is a different vocabulary from the clone’s
  `idle/running/ok/failed/skipped`; see `CONTRACT-PROPOSAL.md` CP-08.

## G-09 — Autosave cadence, crash recovery, reload

- Doc evidence: Local View automatically saves changes to a local repo.
  [build flows](https://learning.postman.com/flows/build-flows/overview)
- Not pinned: debounce interval, save indicator, conflict handling, crash recovery, and whether
  a reload restores the exact viewport.
- Scenario: add/move/configure a block; wait; reload; kill the process and reopen.
- Capture: `network` save calls with timestamps; reload result.
- Status: `needs live builder`.

## G-10 — Undo/redo scope

- Doc evidence: not described on the pages read.
- Not pinned: whether undo covers add, move, delete, connect, config, paste, group, and whether
  the stack survives reload.
- Scenario: perform each edit type, undo/redo each, then reload and undo.
- Capture: `steps` diff; available shortcuts.
- Status: partial. `traces/live-builder/undo-redo.jsonl` shows Ctrl+Z restoring a moved block to
  its previous position. Ctrl+Shift+Z re-applied an intermediate move rather than the exact
  final position, so redo needs a cleaner multi-step test. Add, delete, connect, config, and
  reload survival are not covered yet.

## G-11 — Multi-select, box select, group/ungroup

- Doc evidence: multi-select with the platform modifier is documented in copy/paste; groups are
  mentioned.
  [create blocks](https://learning.postman.com/flows/build-flows/create/create-blocks)
- Not pinned: box-select gesture, additive selection, group bounds, collapse, and whether group
  membership survives copy/paste.
- Scenario: box-select two blocks, group, move, ungroup, copy the group.
- Capture: `dom`/`a11y`; selection counts.
- Status: partial. `traces/live-builder/multiselect.jsonl` shows a drag-selection gesture
  selecting one block; no Group action was exposed in the toolbar. The gesture needs a cleaner
  start point and a second pass for group/ungroup and cross-feature copy/paste.

## G-12 — Limits: blocks, edges, run duration, request count

- Doc evidence: 50 MB per form-data file is documented.
  [HTTP Request](https://learning.postman.com/flows/reference/blocks/http-request)
- Not pinned: maximum blocks, edges, run duration, output size, or request count. A community
  thread reports more than 20 API requests failing under the agent.
  [community](https://community.postman.com/t/not-able-to-run-flows-with-more-than-20-api-requests/57765)
- Scenario: build a flow with N blocks and N requests; run; record the first failure point.
  Treat the community report as a hypothesis, not a limit.
- Capture: `network` errors, console errors, run summary.
- Status: `needs live builder`.

## G-13 — Variable precedence and interpolation

- Doc evidence: flow, environment, collection, and dynamic variables are described.
  [requests and variables](https://learning.postman.com/flows/build-flows/configure/requests-and-variables)
- Not pinned: precedence when the same name exists in several scopes, interpolation syntax,
  resolution time, and whether a missing variable fails the run or resolves empty.
- Scenario: define the same name in two scopes; run; remove the variable; run again.
- Capture: resolved request values and error state.
- Status: `needs live builder`.

## G-14 — Copy/paste across flows: ids, positions, config

- Doc evidence: copy/paste preserves configuration and connections and works across flows.
  [create blocks](https://learning.postman.com/flows/build-flows/create/create-blocks)
- Not pinned: new ids, position offset, whether edges to unselected blocks are dropped, and
  whether cross-flow paste keeps collection references.
- Scenario: copy two connected blocks from flow A to flow B; inspect ids, positions, edges, and
  config.
- Capture: `network` payloads if available; otherwise `dom`/`a11y`.
- Status: `needs live builder`.

## G-15 — Public example read-only layout

- Doc evidence: public Flow Snippets render without a login.
  [Delay example](https://www.postman.com/postman/flows-snippets/flow/63bcbaf6ead4e97c45493eb9)
- Observable now: canvas layout, block labels, configured values, ports, toolbar vocabulary.
- Not observable: add, drag, connect, configure, run, status.
- Scenario: load the Delay snippet; capture a11y/DOM; record block order and port count.
- Capture: `traces/public-delay/` — `steps.jsonl`, `network.jsonl`, `console.jsonl`.
- Status: `covered by traces/public-delay/`. Observed without login: read-only canvas, Delay
  block, Run and run-log actions, add-block disabled, fit-view/auto-layout, agent/sign-in gate,
  and locale JSON bundles. Editing, connection rejection, and run/status remain unobservable
  without an authorized account.

## Next observation session

G-01, G-02, G-06, G-08, G-09 (reload), and G-15 are covered. G-03, G-05, G-07, G-09 (cadence),
G-10, and G-11 are partial.

Next: G-04 (edge select/delete/reconnect), then G-12–G-14 (limits, variable precedence,
cross-flow copy/paste). G-07 needs a configured request that returns a non-2xx; the current
flow’s request editor is hidden and the workspace reports offline. A fresh flow or a reconnected
session should clear that. The agent must not bypass authentication.
