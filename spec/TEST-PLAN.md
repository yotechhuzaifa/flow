# Test plan

Generated from `spec/features.jsonl` by `scripts/spec.mjs`. Do not edit by hand.

## palette.categories — must

- Parity row: `p01` (proposed)
- Observable: the palette shows nine grouped categories in order: Primitives, Connectors, AI, Logic, Looping, Visualize, Outputs, Data, Other
- Trigger: open the block palette from the flow toolbar
- Evidence: traces/live-builder/palette.jsonl#1

Checklist:
- [ ] Playwright test under `tests/parity/`
- [ ] One assertion per observable outcome
- [ ] Preconditions set up explicitly
- [ ] Every expected error exercised
- [ ] Every interactive element has a `data-testid`
- [ ] `spec-auditor` cannot make it pass without the feature

## palette.search — must

- Parity row: `p02` (proposed)
- Observable: searching for delay returns Logic > Delay plus a Suggested String block, and a nonsense query still suggests a String block rather than showing an empty state
- Trigger: type a query into the block palette search field
- Evidence: traces/live-builder/palette.jsonl#2, traces/live-builder/palette.jsonl#3

Checklist:
- [ ] Playwright test under `tests/parity/`
- [ ] One assertion per observable outcome
- [ ] Preconditions set up explicitly
- [ ] Every expected error exercised
- [ ] Every interactive element has a `data-testid`
- [ ] `spec-auditor` cannot make it pass without the feature

## palette.escape — could

- Parity row: `p03` (proposed)
- Observable: pressing Escape closes the palette
- Trigger: press Escape while the block palette is open
- Evidence: traces/live-builder/palette.jsonl#4

Checklist:
- [ ] Playwright test under `tests/parity/`
- [ ] One assertion per observable outcome
- [ ] Preconditions set up explicitly
- [ ] Every expected error exercised
- [ ] Every interactive element has a `data-testid`
- [ ] `spec-auditor` cannot make it pass without the feature

## block.add.click-place — must

- Parity row: `p04` (proposed)
- Observable: the palette closes and the block is placed at the clicked canvas position
- Trigger: select a block in the palette and then click an empty spot on the canvas
- Evidence: traces/live-builder/add-block.jsonl#1, traces/live-builder/add-block.jsonl#2

Checklist:
- [ ] Playwright test under `tests/parity/`
- [ ] One assertion per observable outcome
- [ ] Preconditions set up explicitly
- [ ] Every expected error exercised
- [ ] Every interactive element has a `data-testid`
- [ ] `spec-auditor` cannot make it pass without the feature

## block.add.no-drag — should

- Parity row: `p05` (proposed)
- Observable: dragging from a palette item does not place a block; the select-then-click path is required
- Trigger: try to drag a block directly from the palette
- Evidence: traces/live-builder/add-block.jsonl#1, traces/live-builder/add-block.jsonl#3

Checklist:
- [ ] Playwright test under `tests/parity/`
- [ ] One assertion per observable outcome
- [ ] Preconditions set up explicitly
- [ ] Every expected error exercised
- [ ] Every interactive element has a `data-testid`
- [ ] `spec-auditor` cannot make it pass without the feature

## connection.self.reject — must

- Parity row: `p06` (proposed)
- Observable: a connection preview appears during the drag, but no self-edge is added and no error toast appears
- Trigger: drag a block output port onto its own input port
- Evidence: traces/live-builder/connections-clean.jsonl#1

Checklist:
- [ ] Playwright test under `tests/parity/`
- [ ] One assertion per observable outcome
- [ ] Preconditions set up explicitly
- [ ] Every expected error exercised
- [ ] Every interactive element has a `data-testid`
- [ ] `spec-auditor` cannot make it pass without the feature

## connection.duplicate.reject — must

- Parity row: `p07` (proposed)
- Observable: repeating an existing connection leaves the edge count unchanged and shows no error toast
- Trigger: repeat a connection that already exists
- Evidence: traces/live-builder/connections-clean.jsonl#2

Checklist:
- [ ] Playwright test under `tests/parity/`
- [ ] One assertion per observable outcome
- [ ] Preconditions set up explicitly
- [ ] Every expected error exercised
- [ ] Every interactive element has a `data-testid`
- [ ] `spec-auditor` cannot make it pass without the feature

## connection.cycle.reject — must

- Parity row: `p08` (proposed)
- Observable: connecting a downstream block back to an upstream input leaves the edge count unchanged and shows no error toast
- Trigger: connect a downstream block output back to an upstream input
- Evidence: traces/live-builder/connections-clean.jsonl#3

Checklist:
- [ ] Playwright test under `tests/parity/`
- [ ] One assertion per observable outcome
- [ ] Preconditions set up explicitly
- [ ] Every expected error exercised
- [ ] Every interactive element has a `data-testid`
- [ ] `spec-auditor` cannot make it pass without the feature

## run.status.timeline — must

- Parity row: `p09` (proposed)
- Observable: blocks move from inactive to execution to success; Start reaches success first, then the HTTP block stays in execution until it completes and the downstream blocks mark success
- Trigger: run a flow and watch the block headers
- Evidence: traces/live-builder/status-timeline.jsonl#1

Checklist:
- [ ] Playwright test under `tests/parity/`
- [ ] One assertion per observable outcome
- [ ] Preconditions set up explicitly
- [ ] Every expected error exercised
- [ ] Every interactive element has a `data-testid`
- [ ] `spec-auditor` cannot make it pass without the feature

## run.logs.summary — must

- Parity row: `p10` (proposed)
- Observable: the pane lists each block with its start time and duration and lets you select a block for details
- Trigger: open the run logs after a run
- Evidence: traces/live-builder/runlogs.jsonl#1, traces/live-builder/runlogs.jsonl#3

Checklist:
- [ ] Playwright test under `tests/parity/`
- [ ] One assertion per observable outcome
- [ ] Preconditions set up explicitly
- [ ] Every expected error exercised
- [ ] Every interactive element has a `data-testid`
- [ ] `spec-auditor` cannot make it pass without the feature

## flow.autosave.reload — must

- Parity row: `p11` (proposed)
- Observable: the added block and edges are still present after the reload
- Trigger: edit the flow and reload the page
- Evidence: traces/live-builder/reload.jsonl#1

Checklist:
- [ ] Playwright test under `tests/parity/`
- [ ] One assertion per observable outcome
- [ ] Preconditions set up explicitly
- [ ] Every expected error exercised
- [ ] Every interactive element has a `data-testid`
- [ ] `spec-auditor` cannot make it pass without the feature
