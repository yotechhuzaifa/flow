# Research pipeline (Layer 0 → Layer 5)

This repo already has a build loop (`parity.csv` → code → `dod.sh`). This file is the
missing research loop: how to turn a public product into evidence, features, parity rows,
and executable tests without a human watching the screen.

The field manual is `.opencode/skills/grabbing-flows/SKILL.md`. The missing layer is
`.opencode/skills/spec-compiler/SKILL.md` plus `scripts/spec.mjs`.

## Status

| Layer | State | Where |
|---|---|---|
| 0 — harvest free intel | **partial** | `grabbing-flows` Step 1 + this file |
| 1 — instrumented Playwright trace | complete | `opencode.jsonc` playwright MCP, `grabbing-flows` Step 3 |
| 2 — AI reads traces, not screens | complete | `grabbing-flows` Step 4, `spec/feature.schema.json` |
| 3 — compile features into executable specs | **missing until this change** | `scripts/spec.mjs`, `spec-compiler` skill |
| 4 — synthesise the domain model first | partial | `docs/CONTRACT.md`, `internal/domain/*`, hand-authored |
| 5 — clean-room hygiene | complete | `AGENTS.md` Research rules, `grabbing-flows` intro |

Layer 3 was the hole. Research produced `spec/features.jsonl`, then nothing carried it into
the frozen work queue or the test suite. The compiler closes that seam.

## The missing layer: compile research into work

```
traces/                         evidence captured by playwright MCP
  └─ spec/features.jsonl        one cited behaviour per line
       └─ scripts/spec.mjs check     reject uncited or malformed features
            └─ propose               spec/parity.proposed.csv + spec/parity.map.json
                 └─ human [contract] merge into frozen parity.csv
                      └─ plan        spec/TEST-PLAN.md
                           └─ build  one parity row per iteration
                                └─ audit   spec/AUDIT.md
                                     └─ spec-auditor tries to fake a pass
```

Artifacts:

| File | Owner | Purpose |
|---|---|---|
| `spec/features.jsonl` | `re-observer` | append-only feature list, one JSON object per line |
| `spec/parity.map.json` | `scripts/spec.mjs` | proposed row id → feature id |
| `spec/parity.proposed.csv` | `scripts/spec.mjs` | proposal only; `parity.csv` is frozen |
| `spec/TEST-PLAN.md` | `scripts/spec.mjs` | per-feature checklist the build agent executes |
| `spec/AUDIT.md` | `scripts/spec.mjs` | machine audit of tests, mappings, evidence |
| `spec/CONTRACT-PROPOSAL.md` | `re-observer` / human | OSS-derived contract changes; not approved |
| `spec/RESEARCH-DONE.md` | `re-observer` | every gap covered or marked unobservable |

`traces/` is gitignored on purpose. Evidence stays on the research machine; commit only
paraphrased claims, URLs, and feature lines. `spec.mjs check` verifies a cited trace exists
before the feature is allowed into the pipeline.

Commands:

```bash
node scripts/spec.mjs check      # validate JSONL, ids, priorities, evidence paths
node scripts/spec.mjs propose    # propose parity rows; never edits parity.csv
node scripts/spec.mjs plan       # write the test plan
node scripts/spec.mjs audit      # write the audit; exit 1 on a broken gate
node scripts/spec.mjs all        # all four
make spec                        # same as `all`
```

Rules the compiler enforces:

1. No evidence, no feature.
2. At least one evidence entry must be a real `traces/…` path that exists.
3. A feature id is stable and unique; `p01`-style row ids are assigned once and kept in
   `parity.map.json`.
4. `parity.csv`, `docs/CONTRACT.md`, `scripts/check.sh`, `scripts/dod.sh`, and `GOAL.md`
   are never written by the compiler. Only a human `[contract]` commit merges proposals.
5. The compiler does **not** generate passing tests. It writes a plan. A passing stub is a
   false green and `spec-auditor` exists to catch exactly that.
6. If a feature needs an entity, event, or field that is not in `docs/CONTRACT.md`, write
   `spec/CONTRACT-PROPOSAL.md` and stop. Do not widen the frozen contract from research.

## Layer 0, expanded

Layer 0 is free intel. Do it before opening a browser automation session. Every claim goes
into `spec/INTEL.md` with a URL or a repo path. `spec/GAPS.md` is the numbered list of
behaviours the docs describe but do not pin down; only those gaps become observation
scenarios.

| Source | How | Output | Guardrail |
|---|---|---|---|
| Public docs / Learning Center | `firecrawl_map`, then `firecrawl_scrape` one page per call | node catalogue, semantics, limits | URL per claim |
| Changelog / release blog / "What's new" | `firecrawl_scrape` | dated feature inventory | mark core vs recently added |
| Community feature requests | `firecrawl_search` with `site:reddit.com`, `site:news.ycombinator.com`, `site:stackoverflow.com` | demand list for differentiation | link the thread, no PII |
| i18n / locale JSON | playwright network log while using the product | feature-name inventory and counts | record names, never copy strings |
| Shipped JS / sourcemaps | **off-limits for the target** | none | clean-room rule; see Layer 5 |
| Public API docs / collection schema | `firecrawl_scrape` | entity and serialisation notes | compare with `docs/CONTRACT.md` |
| YouTube demos / walkthroughs | `firecrawl_search`, then scrape the transcript page | jobs-to-be-done and journeys | behaviour only, no screenshots |
| Open-source reference code | `firecrawl` developer search or a raw GitHub file fetch | architecture patterns | read-only; record repo+commit+path+license |

The target's shipped JavaScript and sourcemaps are not read. That is deliberate: the repo is
clean-room, and the bundled locale strings are microcopy. Use the locale JSON only as a
checklist of feature names; never copy the strings, assets, CSS, icons, or class names.

Open-source reference code is different. It may be read to understand architecture and
behaviour, but it is never copied into this repo. Record the repository, commit, file path,
and license in `spec/INTEL.md` for every claim that comes from code.

## Reference projects (behaviour, not code)

Licenses were checked before writing this table. A permissive license still does not make
copying acceptable here; the repo rule is behaviour-only.

| Project | License | What it proves | Do not copy |
|---|---|---|---|
| [n8n](https://github.com/n8n-io/n8n) | Sustainable Use (source-available, not OSI) | a node is a description + execute function + metadata; execution history is a first-class view | any code; the license is not permissive |
| [Node-RED](https://github.com/node-red/node-red) | Apache-2.0 | a node is a runtime JS file + editor HTML file; status and errors are part of the node contract | editor markup, CSS, icons |
| [Langflow](https://github.com/langflow-ai/langflow) | MIT | components are typed classes with declared inputs/outputs; a flow exports as JSON | component code, prompts |
| [Flowise](https://github.com/FlowiseAI/Flowise) | Apache-2.0 (enterprise dirs commercial) | several builders can share one graph/runtime core | enterprise code, UI |
| [Activepieces](https://github.com/activepieces/activepieces) | MIT (EE dirs separate) | a piece/action is declarative metadata over a typed executor | piece implementations |
| [Windmill](https://github.com/windmill-labs/windmill) | AGPL-3.0 | flows are DAGs in a versioned JSON format; scheduling, errors, early return are engine concerns | backend/frontend code; AGPL is viral |
| [ComfyUI](https://github.com/comfyanonymous/ComfyUI) | GPL-3.0 | nodes have typed inputs/outputs and explicit run states | node code; GPL is viral |
| [React Flow](https://github.com/xyflow/xyflow) | MIT | `isValidConnection` validates before an edge exists; undo/redo is a snapshot stack | example code, styling |

Patterns worth adopting, in order of value:

1. **One node manifest drives everything.** n8n has `*.node.ts` + codex JSON, Node-RED has a
   JS file + HTML file, Langflow has a typed component class. This repo currently repeats a
   node across `internal/nodes/<type>.go`, `frontend/src/features/nodes/`, `OUT_PORTS`, a
   parity row, and a test. That is the largest drift risk. A future `spec/nodes.json` should
   generate the palette, port list, config schema, and registry. It is a refactor, not a
   research change; propose it before doing it.
2. **Type the ports.** Langflow and ComfyUI declare input/output types; React Flow can reject
   a connection before it is added (`isValidConnection`). The current `CONTRACT.md` has
   untyped ports and rejects self-edges/cycles at save time. Add port types to the node
   manifest before adding third-party nodes.
3. **Version the flow document.** Windmill uses OpenFlow JSON; n8n and Node-RED version their
   flow formats. `Flow` in `CONTRACT.md` has no `schemaVersion`. SQLite has migrations; flow
   JSON does not. Propose `schemaVersion` + a migration function in `internal/domain/flow`.
4. **Persist per-node run data.** n8n keeps execution history with per-node data; Node-RED
   exposes debug/status. This repo stores only a `runs` row and keeps `nodeStatus` in memory.
   A `run_nodes` table or a JSON blob per run makes failures debuggable.
5. **Cancellation, retries, timeouts belong to the executor.** Windmill treats them as
   orchestration; the node stays simple. `internal/exec` currently has no context. Adding
   `context.Context`, `run.cancelled`, and per-node retry/timeout belongs there, not in nodes.
6. **Golden flow fixtures.** The references test both a node in isolation and a whole flow.
   Add `testdata/flows/*.json` and table-driven `exec` tests. Playwright is for the UI, not
   for execution semantics.
7. **Spec-to-test is not solved by the references.** None of them generate executable UI
   tests from a feature list. Keep the compiler small: validate, propose, plan, audit. Do not
   auto-generate assertions; a generated test that passes proves nothing.

Items 1–5 change the frozen contract or the executor; item 6 is additive. Write contract
changes to `spec/CONTRACT-PROPOSAL.md` first. `docs/CONTRACT.md` only changes in a human
`[contract]` commit.

## Minimal-human loop

1. **Human, once:** fill `url:` in `spec/TARGET.md`.
2. **Agent:** `bash scripts/research.sh` — runs `/research`, then the compiler.
3. **Human, once per batch:** review `spec/parity.proposed.csv` and merge the accepted rows
   into `parity.csv` in a `[contract]` commit. This is the only mandatory gate.
4. **Agent:** `./scripts/loop.sh` — builds one parity row per iteration until `dod.sh` is green.
5. **Agent:** `/next` runs `node scripts/spec.mjs audit`; `spec-auditor` attacks each new row
   and adds a failing test if the feature can be faked.
6. **Human:** resolve `BLOCKED.md`, `MAX_ITERS`, or a contract proposal. Everything else is
   unattended.

Remaining human touchpoints and why:

| Touchpoint | Why it cannot be automated safely |
|---|---|
| Target URL + scope | a product choice, not a research result |
| Login/consent | `AGENTS.md` forbids bypassing authentication or touching another user's data |
| `parity.csv` / `CONTRACT.md` merge | frozen by design; prevents an agent editing its own goal |
| Desktop packaging | Wails v3 `build/` assets and a WebView toolchain are not vendored |
| `BLOCKED.md` | the loop stops after repeated failure instead of guessing |

## MCP setup

`opencode.jsonc` wires the research tools:

| Server | Use |
|---|---|
| `firecrawl` | docs, changelog, public API, OSS code search, community search, page scraping |
| `playwright` | instrumented traces: network, console, a11y, DOM; headless |
| `context7` | current library docs (`@xyflow/react`, Playwright, Wails) |

The Playwright server runs with `--browser chromium`, `--caps network,console,testing`, and no
`vision` by default. Vision is a fallback only when a control is unreachable through the
accessibility tree. If it reports a missing `chrome` binary, the config is pointing at the
wrong browser; use Chromium and run `npx playwright install chromium`.

`RivalSearchMCP` is optional. It has been unavailable in this environment, so the documented
Layer 0 path uses `firecrawl` with `site:` filters and does not depend on it.

The workflow is intentionally split: Layer 0 uses `firecrawl` and costs no browser session;
Layer 1 uses `playwright` only for the gaps that the docs left open.
