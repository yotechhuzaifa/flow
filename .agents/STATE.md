# STATE

Append-only. Newest at the bottom.

## Current token
- Task: (none)
- Files:
- Verify with:

## Log
| # | Row | Check | Result | Note |
|---|-----|-------|--------|------|
| 0 | scaffold | ./scripts/check.sh | — | initialised |
| 1 | setup audit | ./scripts/check.sh | fail → pass | missing `internal/webui/dist` broke `go build`; vite alias pointed at a non-existent stub; `@playwright/test` was not a dependency. Dependency added: `@playwright/test` (the parity suite in `tests/parity` imports it; no runtime cost). |
| 2 | setup audit | `check.sh` · `npx playwright test` · `opencode debug config` | pass | check.sh PASS; f01 passes against Vite; config parses with all three MCP servers. Initial commit `2c7eeca` lets the `dod.sh` frozen-file guard compare against HEAD. `dod.sh` now fails only on the 12 open MUST rows. |
| 3 | research pipeline | `node scripts/spec.mjs check` · `bash scripts/check.sh` | pass | Added the missing Layer 3 compiler (`scripts/spec.mjs`), `spec/TARGET.md`, `scripts/research.sh`, `docs/RESEARCH.md`, the `spec-compiler` skill, and the `audit-spec` command. Expanded Layer 0 sources and the OSS reference table. No package dependency; Node built-ins only. `spec.mjs audit` is expected to FAIL until the open MUST rows have tests. Layer 4 stays hand-authored; OSS-derived changes are written to `spec/CONTRACT-PROPOSAL.md` and need a human `[contract]` commit. |
| 4 | target | `sed -n 's/^[[:space:]]*url:[[:space:]]*//p' spec/TARGET.md` · `node scripts/spec.mjs check` | pass | Target set to Postman Flows, docs entry `https://learning.postman.com/flows/overview`. Scope stays inside `GOAL.md`; Blocks/FQL are intel only. The live builder needs a human Postman account and must not be bypassed. Target source/sourcemaps remain off-limits. |
| 5 | Layer 0 intel | `node scripts/spec.mjs check` · `bash scripts/check.sh` | pass | Wrote `spec/INTEL.md` and `spec/GAPS.md` from the official docs, release notes, blog, community threads, and a public Flow Snippet. 15 gaps; the first three need the live builder. Fixed the Playwright MCP to select Chromium after it failed on a missing `chrome` binary. The public Delay example was observed read-only without login; trace at `traces/public-delay/`. |
| 6 | live builder observation | `node scripts/spec.mjs check` · `propose` · `plan` | pass | Authenticated session in SiamSami's Workspace. Captured palette order/search, click-to-place add, silent rejection of self/duplicate/cycle connections, the inactive→execution→success status timeline, run logs, and autosave on reload. Wrote 11 features; compiler produced 11 proposed rows and a test plan. Audit fails as expected until the open MUST rows have tests. |
| 7 | contract decision | `git diff HEAD -- GOAL.md parity.csv` · `bash scripts/check.sh` | pass | Human owner accepted CP-07 and CP-08. `GOAL.md` and f02 now use select-and-place; `p01`–`p11` merged into `parity.csv`. CP-08 keeps the internal event names and maps the target states in the UI. CP-01–CP-06 remain unapproved. |
| 8 | live builder round 2 | `node scripts/spec.mjs check` · `propose` · `plan` | pass | Captured Run→Stop control, Delay config clamping, an unconfigured HTTP block entering `error`, move undo, and a partial box-select. Added 4 features (`p12`–`p15`). `propose` now lists only rows not already in `parity.csv`. G-07 remains partial; G-12–G-14 not attempted. |
| 9 | live builder round 3 | `node scripts/spec.mjs check` · `propose` · `plan` | pass | Confirmed edge deletion by midpoint control: four edges to three, then five to four, with no confirmation dialog. Added `edge.delete` (`p16`). Reconnect and connected-block deletion remain open. |
| 10 | live builder round 4 | `node scripts/spec.mjs check` · `bash scripts/check.sh` | pass | Autosave cadence is not exposed as a separate HTTP save; reload persistence is already covered. Ctrl+click selected only one node; no Group action. Cross-flow copy/paste is unobservable: workspace offline and only one flow. Limits and variable precedence need dedicated setups. No new features. |

## Known gaps (setup audit, not fixed here)
- Desktop build needs the Wails v3 `build/` assets, which are not vendored. `wails3` cannot install on a
  host without `pkg-config`/GTK4. Headless (`cmd/server`, Vite, Playwright) is fully usable.
- `frontend/src/shared/api/client.ts` talks only to generated Wails bindings; there is no HTTP client for
  server mode yet. Later rows that need persistence/live runs in the browser must add it.
- The generated binding path `@bindings/flow/api` is a declared placeholder; verify it after
  `wails3 generate bindings` on a machine with the WebView toolchain.
