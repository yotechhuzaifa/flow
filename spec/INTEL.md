# INTEL — Postman Flows (Layer 0)

Captured 2026-11-23 from the official docs, release notes, blog, community forum, and public
example flows. Everything below is paraphrased behaviour. No target microcopy, assets, CSS,
icons, class names, or source were copied.

Docs version: Latest (v12). Scope: only the surfaces in `GOAL.md` — palette, canvas, nodes,
edges, run with live status, persistence. Cloud deploy, scheduling, AI, MCP, connectors, and
collaboration are recorded as context, not clone scope.

## Product shape

- Postman Flows is a visual low-code editor for API workflows made of blocks on a canvas.
  [docs overview](https://learning.postman.com/flows/overview)
- A flow is built from blocks; blocks perform a task or store data and connect through ports.
  [blocks overview](https://learning.postman.com/flows/reference/blocks/overview)
- Flows can run locally or in the Postman Cloud. Cloud runs add schedules, webhooks, deployed
  URLs, and MCP. The clone is local-first; cloud behaviour is out of scope.
  [build flows](https://learning.postman.com/flows/build-flows/overview)
- Two storage modes exist: Cloud View and Local View. Local View saves to a local repo and
  supports Native Git; Cloud View saves to Postman Cloud.
  [build flows](https://learning.postman.com/flows/build-flows/overview)
- Snapshots version a flow. Local View does not support snapshots.
  [snapshots](https://learning.postman.com/flows/build-flows/configure/snapshots)

## Palette

- A new flow starts with a Start block. The Start block triggers connected blocks, but
  connecting blocks to Start is optional.
  [start block](https://learning.postman.com/flows/reference/blocks/start)
- Blocks can be added three ways: right-click the canvas, use the toolbar Block action, or drag
  from an existing block's input/output port. All open a block list.
  [create blocks](https://learning.postman.com/flows/build-flows/create/create-blocks)
- The Start block also offers shortcuts: add blocks, create an API example, create from a
  collection, and generate with AI.
  [create blocks](https://learning.postman.com/flows/build-flows/create/create-blocks)
- Block categories: primitives/action, connectors, trigger, AI, logic, looping, visualize,
  outputs, and data.
  [blocks overview](https://learning.postman.com/flows/reference/blocks/overview)
- Copy/paste preserves block configuration and connections; multiple blocks can be selected
  with the platform modifier. Copy works across flows.
  [create blocks](https://learning.postman.com/flows/build-flows/create/create-blocks)
- A selected block has a context menu with duplicate, change colour, pause on block, copy,
  delete, and Agent Mode.
  [blocks overview](https://learning.postman.com/flows/reference/blocks/overview)
- Groups and annotations exist for organising a canvas.
  [build flows](https://learning.postman.com/flows/build-flows/overview)

## Canvas

- The canvas is a spatial editor; blocks are placed, moved, grouped, and connected.
  [docs overview](https://learning.postman.com/flows/overview)
- A block can be paused before it runs. Paused state is visible on the block, and the run
  resumes with a Continue action.
  [blocks overview](https://learning.postman.com/flows/reference/blocks/overview)
- Clicking a block previews its input and output data.
  [troubleshoot](https://learning.postman.com/flows/build-flows/troubleshoot/troubleshoot)
- Undo/redo is not described on the pages read. This is a gap, not a confirmed absence.
  [gap](./GAPS.md)

## Ports and edges

- Connections are unidirectional conduits that carry data between blocks.
  [create connections](https://learning.postman.com/flows/build-flows/create/create-connections)
- Most blocks have an input port, an output port, or both. Many output ports can fan out to
  multiple blocks.
  [create connections](https://learning.postman.com/flows/build-flows/create/create-connections)
- Dragging from an output port starts a connection; dropping on empty canvas can create a new
  block with the connection.
  [create connections](https://learning.postman.com/flows/build-flows/create/create-connections)
- Some blocks accept data blocks inline, so a value-producing block can live inside another
  block's configuration.
  [blocks overview](https://learning.postman.com/flows/reference/blocks/overview)
- The docs do not pin the rejection behaviour for self-edges, cycles, duplicate edges,
  incompatible ports, or a dropped connection. This is a gap.
  [gap](./GAPS.md)

## Node config

- Blocks are configured in place. The HTTP Request block selects a request from a collection,
  and the request's variables become input ports on the block.
  [HTTP Request](https://learning.postman.com/flows/reference/blocks/http-request)
- HTTP Request inputs: `Send` (optional trigger) and `Variables`. Outputs: `Success` for 2xx and
  `Fail` for any non-2xx response.
  [HTTP Request](https://learning.postman.com/flows/reference/blocks/http-request)
- HTTP Request can parse responses as Auto, JSON, XML, HTML, or Text.
  [HTTP Request](https://learning.postman.com/flows/reference/blocks/http-request)
- The If block has a variable input, a data input, an expression (FQL or TypeScript), and
  `Then` / `Else` outputs. Data is routed based on the expression; the routed payload is the
  data input, not the expression variables.
  [If block](https://learning.postman.com/flows/reference/blocks/if)
- The Delay block waits a configured number of milliseconds per received payload and preserves
  arrival order when multiple blocks feed it.
  [Delay block](https://learning.postman.com/flows/reference/blocks/delay)
- The Log block sends its input to the Postman Console; multiple inputs are logged in arrival
  order. HTTP requests are logged automatically as network events.
  [Log block](https://learning.postman.com/flows/reference/blocks/log)
- Variables can come from flow, environment, collection, or dynamic sources. The docs describe
  insertion and override per request, but not precedence when the same name exists in several
  scopes. This is a gap.
  [requests and variables](https://learning.postman.com/flows/build-flows/configure/requests-and-variables)

## Run and status

- The flow runs from the Run action. A connected `Send` input is optional for HTTP Request; the
  block runs when the flow runs even if `Send` is unconnected.
  [HTTP Request](https://learning.postman.com/flows/reference/blocks/http-request)
- Run logs show each run by start time. Selecting a run and a block reveals input, output, test
  results, eval results, and details. The pane filters by block type and run status and can
  export all test results.
  [troubleshoot](https://learning.postman.com/flows/build-flows/troubleshoot/troubleshoot)
- Cloud analytics add time-range and run filters, a Run Details pane, and per-block
  success/fail status. A run ID is exposed in the `x-session-id` response header.
  [analytics](https://learning.postman.com/flows/build-flows/analyze/analytics)
- The docs show final success/fail and a paused state, but not the full status vocabulary
  (idle/running/skipped), transition timing, or whether descendants of a failed block are
  skipped. These are gaps.
  [gap](./GAPS.md)
- Stop/cancel, concurrent runs, and re-run semantics are not described on the pages read.
  [gap](./GAPS.md)

## Persistence

- Cloud View stores flows in Postman Cloud. Local View automatically saves changes to a local
  repo and can sync with a remote repo through Native Git.
  [build flows](https://learning.postman.com/flows/build-flows/overview)
- Snapshots save, view, and restore flow versions. Viewing a snapshot loads it read-only but
  still runnable; editing requires restoring it. The newest snapshot is marked as latest, and
  the first is numbered.
  [snapshots](https://learning.postman.com/flows/build-flows/configure/snapshots)
- The docs do not pin autosave cadence, crash recovery, or reload behaviour for Local View.
  [gap](./GAPS.md)
- The clone uses SQLite, not Cloud View or Git. That is a deliberate difference, not a target
  behaviour to copy. [GOAL](../GOAL.md)

## Limits and costs

- HTTP Request supports form-data file attachments up to 50 MB per file.
  [HTTP Request](https://learning.postman.com/flows/reference/blocks/http-request)
- Cloud blocks consume credits; manual runs only consume credits for blocks that run in the
  cloud.
  [billing](https://learning.postman.com/docs/billing/flows-usage)
- A community thread reports flows with more than 20 API requests failing to run through the
  Postman Agent. This is a report, not a documented limit.
  [community](https://community.postman.com/t/not-able-to-run-flows-with-more-than-20-api-requests/57765)
- No documented maximum block count, edge count, or run duration was found on the pages read.
  This is a gap.
  [gap](./GAPS.md)

## Community signal (feature requests and pain)

These are demand signals, not specifications. They are useful for differentiation in
`docs/SCALE.md`, not for the clone's MUST rows.

- The flow builder's usefulness beyond testing is questioned.
  [community](https://community.postman.com/t/postman-flow-builder-usability-beyond-testing/67874)
- Environment variables set during a flow do not always persist as expected.
  [community](https://community.postman.com/t/postman-flows-not-updating-environment-variables/48277)
- The Run action sometimes appears to do nothing.
  [community](https://community.postman.com/t/i-click-run-in-flows-postman-and-nothing-happens/57785)
- Large flows with more than 20 API requests fail under the agent.
  [community](https://community.postman.com/t/not-able-to-run-flows-with-more-than-20-api-requests/57765)
- A flow can run infinitely when a connection bypasses a gate.
  [community](https://community.postman.com/t/postman-flow-runs-infinitely/75303)
- Request body display and sub-request selection have been reported as confusing.
  [community](https://community.postman.com/t/flow-request-block-response-body-displays-sub-request/62840)

## Changelog and release notes

- The official release notes are versioned and include Flows entries, for example Flows in the
  EU region and Native Git for Flows.
  [release notes](https://www.postman.com/release-notes/postman-app/)
- The Postman blog carries Flows release posts and product updates.
  [blog](https://blog.postman.com/whats-new-in-postman-flows-fabric-gateway-and-astro-ai/)
  [updates](https://blog.postman.com/december-2025-product-updates/)
- Docs are versioned (v12 latest, v11 archive), so a behaviour claim must state its version.
  [docs overview](https://learning.postman.com/flows/overview)

## Public example flows

Public Flow Snippets render read-only without a login. They show the canvas, block labels,
configured values, and ports, but Run requires the Postman Agent and a sign-in. They are a
legitimate observation source for layout and vocabulary, not for run/status.
[Delay example](https://www.postman.com/postman/flows-snippets/flow/63bcbaf6ead4e97c45493eb9)

Observed with Playwright on 2026-11-23 (`traces/public-delay/`): a read-only canvas rendered
with a Delay block and a value field; the toolbar exposed Run, a run-log action, fit-view, and
auto-layout; add-block was disabled; the run path reported the agent unavailable and offered
sign-in; some workspace API calls returned 403 without auth; the page fetched locale JSON
bundles. No run was attempted.

## Live builder observation (2026-11-23)

Authenticated session in SiamSami's Workspace, flow
`/workspace/093bebad-aa29-49af-80e2-33936570c6ca/flow/6b01c71bf452684a9606f801`.
Evidence: `traces/live-builder/`.

- The block palette opens from the toolbar and lists nine categories in order: Primitives,
  Connectors, AI, Logic, Looping, Visualize, Outputs, Data, Other.
- The palette has a search field. `delay` returns Logic → Delay plus Suggested → String “delay”.
  A nonsense query still suggests a String block; there is no empty state. Escape closes it.
- Blocks are added by selecting a palette item and then clicking the canvas. There is no
  drag-from-palette path; palette items are not draggable. This differs from `GOAL.md` and
  parity row f02.
- Invalid connections are rejected silently. Self-edge, duplicate, and cycle attempts each left
  the edge count unchanged and produced no error toast; a connection preview rendered during
  the drag.
- Run status values are `inactive`, `execution`, and `success`. Start completes first; an HTTP
  block stays in `execution` until it finishes, then downstream blocks mark success.
- The run logs pane lists each block with a start time and duration and opens block details,
  including input, output, tests, evals, and further details.
- A reload preserved the added Delay block and four edges, so the flow autosaved before the
  reload completed. The exact debounce interval was not measured.
- A happy-path HTTP run completed successfully (about 2–3 s). Failure propagation and skipped
  descendants were not forced, so `GAPS.md` G-07 remains open.

## Source log

- https://learning.postman.com/flows/overview
- https://learning.postman.com/flows/build-flows/overview
- https://learning.postman.com/flows/build-flows/create/create-blocks
- https://learning.postman.com/flows/build-flows/create/create-connections
- https://learning.postman.com/flows/reference/blocks/overview
- https://learning.postman.com/flows/reference/blocks/start
- https://learning.postman.com/flows/reference/blocks/http-request
- https://learning.postman.com/flows/reference/blocks/if
- https://learning.postman.com/flows/reference/blocks/delay
- https://learning.postman.com/flows/reference/blocks/log
- https://learning.postman.com/flows/build-flows/configure/requests-and-variables
- https://learning.postman.com/flows/build-flows/configure/snapshots
- https://learning.postman.com/flows/build-flows/troubleshoot/troubleshoot
- https://learning.postman.com/flows/build-flows/analyze/analytics
- https://learning.postman.com/docs/billing/flows-usage
- https://www.postman.com/release-notes/postman-app/
- https://blog.postman.com/whats-new-in-postman-flows-fabric-gateway-and-astro-ai/
- https://blog.postman.com/december-2025-product-updates/
- Community threads linked inline above.
