---
name: grabbing-flows
description: How to observe a third-party visual flow builder and turn it into spec data. Load before any research or observation session.
autoinvoke: false
---

# Observing a flow builder

You are the spec team. You never read their source. You produce evidence.

## Step 1 — Layer 0: free intel first (no automation)
Read `docs/RESEARCH.md` for the full source list and the open-source reference projects.
Harvest in this order, one page per call:

1. docs / Learning Center — `firecrawl_map`, then `firecrawl_scrape`.
2. changelog / release blog / "What's new" — `firecrawl_scrape`; record dates.
3. community feature requests — `firecrawl_search` with `site:` filters.
4. public API docs / collection schema — `firecrawl_scrape`.
5. demos / walkthroughs — `firecrawl_search`; record journeys, not frames.
6. open-source reference code — read-only; record repo + commit + path + license.

Write `spec/INTEL.md`: features grouped by area (palette, canvas, ports/edges, node config,
run, status, persistence). Cite the URL or repo path per claim. If a page is not reachable,
skip it — do not guess.

Never read the target's shipped JS or sourcemaps. Never copy target microcopy, assets, CSS,
icons, class names, or code. Open-source reference code may be read for architecture, but it
is never copied into this repo.

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
Append to `spec/features.jsonl`, then run `node scripts/spec.mjs check`:

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

## Step 5 — compile, then stop
Run `node scripts/spec.mjs propose` and `node scripts/spec.mjs plan`. The compiler writes
`spec/parity.proposed.csv`, `spec/parity.map.json`, and `spec/TEST-PLAN.md`. It never edits
`parity.csv`; a human merges accepted rows with a `[contract]` commit.

Write `spec/RESEARCH-DONE.md` when every gap in `GAPS.md` is either covered or marked
`unobservable, reason: …`. Then the research phase is over and implementation starts.
Load the `spec-compiler` skill for the handoff.
