---
name: grabbing-flows
description: How to observe a third-party visual flow builder and turn it into spec data. Load before any research or observation session.
autoinvoke: false
---

# Observing a flow builder

You are the spec team. You never read their source. You produce evidence.

## Step 1 — public intel first (no automation)
Use `webfetch` on their docs, changelog, and help centre. One page at a time.
Write `spec/INTEL.md`: features grouped by area (palette, canvas, ports/edges, node config,
run, status, persistence). Cite the URL per claim. If a page is not reachable, skip it — do not guess.

Result: `spec/GAPS.md`, a numbered list of behaviours that the docs describe but do not pin down.
Those gaps are your observation scenarios. Nothing else is.

## Step 2 — one scenario per trace
One interaction path per session. Example scenario lists:

1. open the builder, find the palette, note the node types shown
2. drag one node onto the canvas
3. open the node, change its config, close, reopen
4. connect two nodes port to port
5. delete a node that has an edge
6. move a node, reload the page, check the position survived
7. run the flow, watch status per node
8. make a node fail, watch what happens downstream

## Step 3 — capture, per action
Use the playwright MCP. For each action record:

- `network`: request/response bodies (this is the data model, and the strongest evidence)
- `a11y`: accessibility snapshot before and after (structure without pixels)
- `dom`: elements added or removed, changed text, changed attributes
- `console`: log and error text

Store under `traces/<scenario>/` as `steps.jsonl` and `network.jsonl`.
Use `--caps network,console,testing`. Only enable `vision` if a control is not reachable through a11y.

## Step 4 — extract
Two passes, never one: a per-trace pass, then a merge pass over all of `spec/features.jsonl`.
Append to `spec/features.jsonl`:

```json
{"id":"edge.create","trigger":"dragging from an output port to an input port",
 "observable":"an edge appears only when the drop lands on an input port",
 "wire":["POST /api/flows/:id/edges"],"errors":["dropping on empty canvas removes the pending edge"],
 "priority":"must","evidence":["traces/connect/network.jsonl#req-4","traces/connect/steps.jsonl#7"]}
```

Rules:
- No `evidence` → do not write the feature.
- Behaviour only. Never name their functions, CSS classes, or internal ids.
- Where network and UI disagree, network wins; record the disagreement.
- One observable outcome per feature.

## Step 5 — stop
Write `spec/RESEARCH-DONE.md` when every gap in `GAPS.md` is either covered or marked
`unobservable, reason: …`. Then the research phase is over and implementation starts.
